package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"practice_1/db"
	"practice_1/internal"
	"strings"
	"time"
)

const (
	ResultOK = iota + 1
	ResultError
)

func MainPageHandle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	//parse, then send to channel struct
	separatedJson := strings.Split(string(body), "\n")
	if len(separatedJson) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resultChan := make(chan int)
	go func() {
		defer close(resultChan)
		structs := internal.ParseAndExpand(separatedJson, r)
		if len(structs) == 0 {
			resultChan <- ResultError
		}
		fmt.Println(structs)
		db.GetInstance().AddData(structs)
		resultChan <- ResultOK
	}()
	select {
	case res := <-resultChan:
		if res != ResultOK {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
