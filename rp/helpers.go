package rp

import (
	"fmt"
	"net/http"
	"time"
)

// toTimestamp returns unix timestamp for time object
func toTimestamp(t time.Time) int64 {
	return t.Unix() * int64(time.Microsecond)
}

// doRequest do request with authorization token
func doRequest(req *http.Request, token string) (*http.Response, error) {
	auth := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", auth)

	client := http.Client{}
	return client.Do(req)
}
