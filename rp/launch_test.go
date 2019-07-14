package rp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLaunch(t *testing.T) {
	l := NewLaunch(nil, "test", "test", "test", nil)
	assert.NotNil(t, l)
}

func TestStartLaunch(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test_project/launch", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusCreated)
		resp := `{"id": "testid123"}`
		w.Write([]byte(resp))
	}))
	c := &Client{
		Endpoint: testServer.URL,
		Project:  "test_project",
		Token:    "1234",
	}
	l := NewLaunch(c, "test launch", "test description", "test mode", nil)
	err := l.Start()

	assert.NoError(t, err)
}
