package rp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToTimestamp(t *testing.T) {
	u, _ := time.Parse("2006-01-02", "2019-01-01")
	ts := int64(1546300800000) // Jan 1 2019 timestamp
	assert.Equal(t, ts, toTimestamp(u))
}

func TestDoRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer 1234", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	mockReq := httptest.NewRequest(http.MethodGet, mockServer.URL, nil)
	doRequest(mockReq, "1234")
}
