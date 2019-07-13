package rp

import (
	"net/http"
	"time"
)

// toTimestamp returns unix timestamp for time object
func toTimestamp(t time.Time) int64 {
	return t.Unix() * int64(time.Microsecond)
}

// doRequest
func doRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(req)
}
