package internal

import (
	"encoding/json"
	"net"
	"net/http"
	"practice_1/db"
	"strings"
	"time"
)

func getIP(req *http.Request) (string, error) {
	ip := req.Header.Get("X-REAL-IP")
	if ip != "" {
		return ip, nil
	}
	ips := strings.Split(req.Header.Get("X-FORWARDED-FOR"), ",")
	if len(ips) > 0 && ips[0] != "" {
		return ips[0], nil
	}
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, err
}

func ParseAndExpand(body []string, req *http.Request) []db.EventLogExtended {
	result := make([]db.EventLogExtended, 0, len(body))
	for _, curStr := range body {
		var curStruct db.EventLogExtended
		err := json.Unmarshal([]byte(curStr), &curStruct)
		if err != nil {
			continue
		}
		curStruct.ServerTime = time.Now()
		ipStr, err := getIP(req)
		if err != nil {
			continue
		}
		curStruct.Ip = net.ParseIP(ipStr)
		result = append(result, curStruct)
	}
	return result
}
