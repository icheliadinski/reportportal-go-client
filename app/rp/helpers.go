package rp

import (
	"fmt"
	"net/http"
	"time"
)

func addContentTypeJSON(r *http.Request) {
	r.Header.Set("Content-Type", "application.json")
}

func toTimestamp(t time.Time) int64 {
	return t.Unix() * int64(time.Microsecond)
}

func doRequest(req *http.Request, token string) (*http.Response, error) {
	client := http.Client{}
	if token != "" {
		auth := fmt.Sprintf("Bearer %s", token)
		req.Header.Set("Authorization", auth)
	}
	return client.Do(req)
}
