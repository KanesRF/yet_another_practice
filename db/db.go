package db

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

//How mach structs can we store before INSERT db call?
const StorageLimit int32 = 30000

type ClientDate struct {
	time.Time
}

func (c *ClientDate) UnmarshalJSON(b []byte) error {
	var err error
	layout := "2006-01-02 15:04:05"
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}
	c.Time, err = time.Parse(layout, s)
	return err
}

type EventLog struct {
	ClientTime ClientDate `json:"client_time"`
	DeviceId   string     `json:"device_id"`
	DeviceOS   string     `json:"device_os"`
	Session    string     `json:"session"`
	Sequence   uint64     `json:"sequence"`
	Event      string     `json:"event"`
	ParamInt   int64      `json:"param_int"`
	ParamStr   string     `json:"param_str"`
}

type EventLogExtended struct {
	EventLog
	Ip         net.IP    `json:"ip"`
	ServerTime time.Time `json:"server_time"`
}

//Can we add interface? Yes. But what the point, if there will be no expansion of codebase?
type DBInstance struct {
	conn      clickhouse.Conn
	mtx       sync.Mutex
	stop_chan chan int
	storage   []EventLogExtended
}

var db_instance *DBInstance

func (db *DBInstance) flush() {
	defer func() {
		db.storage = db.storage[:0]
		db.mtx.Unlock()
	}()
	db.mtx.Lock()
	if len(db.storage) == 0 {
		return
	}
	batch, err := db.conn.PrepareBatch(context.Background(), "INSERT INTO event_db.events_buf")
	if err != nil {
		infoLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
		infoLog.Println(err)
		return
	}
	for _, curEvent := range db.storage {
		// Why not using AppendStruct? Becouse of name policy in Clickhouse and Golang
		//err = batch.AppendStruct(&curEvent)
		err := batch.Append(
			curEvent.ClientTime.Time,
			curEvent.DeviceId,
			curEvent.DeviceOS,
			curEvent.Session,
			curEvent.Sequence,
			curEvent.Event,
			curEvent.ParamInt,
			curEvent.ParamStr,
			curEvent.Ip,
			curEvent.ServerTime)
		if err != nil {
			infoLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
			infoLog.Println(err)
		}
	}

	err = batch.Send()
	if err != nil {
		infoLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
		infoLog.Println(err)
	}
}

//it is just signal to stop. Anyway, it is not supposed to be stoped
func (db *DBInstance) StopInstance() {
	close(db.stop_chan)
}

func (db *DBInstance) flushLoop() {
	defer db.conn.Close()
	select {
	case <-db.stop_chan:
		db.flush()
		break
	case <-time.After(time.Second):
		db.flush()
	}
}

func GetInstance() *DBInstance {
	var init sync.Once
	init.Do(func() {
		clickConn, err := clickhouse.Open(&clickhouse.Options{
			Addr: []string{"127.0.0.1:9000"},
			Auth: clickhouse.Auth{
				Database: "event_db",
				Username: "event_writer",
				Password: "passwd",
			},
			DialTimeout:     time.Second,
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		})
		if err != nil {
			return
		}
		db_instance = new(DBInstance)
		db_instance.conn = clickConn
		db_instance.stop_chan = make(chan int)
		db_instance.storage = make([]EventLogExtended, 0, StorageLimit)
		go db_instance.flushLoop()
	})
	return db_instance
}

func (db *DBInstance) AddData(data []EventLogExtended) {
	defer db.mtx.Unlock()
	db.mtx.Lock()
	db_instance.storage = append(db_instance.storage, data...)
}
