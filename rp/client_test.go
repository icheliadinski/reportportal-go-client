package rp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient("", "", "", 0)
	assert.NotNil(t, c)

	var endpoints = []struct {
		endpoint string
		version  int
		expected string
	}{
		{"rp.epam.com", 0, "https://rp.epam.com/api/v1"},
		{"rp.epam.com", -1, "https://rp.epam.com/api/v1"},
		{"rp.epam.com", 1, "https://rp.epam.com/api/v1"},
		{"rp.epam.com", 2, "https://rp.epam.com/api/v2"},
		{"https://rp.epam.com", 1, "https://rp.epam.com/api/v1"},
		{"https://rp.epam.com/", 1, "https://rp.epam.com/api/v1"},
		{"rp.epam.com/api/v1", 1, "https://rp.epam.com/api/v1"},
		{"http://rp.epam.com", 1, "http://rp.epam.com/api/v1"},
	}

	for _, tt := range endpoints {
		actual := NewClient(tt.endpoint, "", "", tt.version)
		if actual.Endpoint != tt.expected {
			t.Errorf(`NewClient(%s, "", "", %d): expected %s, actual %s`, tt.endpoint, tt.version, tt.expected, actual.Endpoint)
		}
	}
}

func TestCheckConnect(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("response"))
	}))
	defer testServer.Close()

	c := NewClient(testServer.URL, "test_project", "1234", 1)
	err := c.CheckConnect()
	assert.NoError(t, err)
}
