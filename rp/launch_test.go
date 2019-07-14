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
	t.Run("Correctly created", func(t *testing.T) {
		okResponse := `{"id": "testid"}`
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "/test_project/launch", r.URL.Path)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(okResponse))
		})
		s := httptest.NewServer(h)

		c := &Client{
			Endpoint: s.URL,
			Project:  "test_project",
		}
		l := &Launch{
			client:      c,
			Name:        "test launch",
			Description: "test description",
			Mode:        "test mode",
			Tags:        nil,
		}
		err := l.Start()

		assert.Equal(t, "testid", l.Id)
		assert.NoError(t, err)
	})

	t.Run("Differ status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		s := httptest.NewServer(h)

		c := &Client{
			Endpoint: s.URL,
			Project:  "test_project",
		}
		l := &Launch{
			client:      c,
			Name:        "test launch",
			Description: "test description",
			Mode:        "test mode",
			Tags:        nil,
		}
		err := l.Start()

		assert.Error(t, err)
		assert.Equal(t, err.Error(), "failed with status 200 OK")
	})
}
