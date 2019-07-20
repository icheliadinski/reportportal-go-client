package rp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestItem(t *testing.T) {
	l := &Launch{}
	ti := NewTestItem(l, "", "", "", nil, nil)
	assert.NotNil(t, ti)
}

func TestStartTestItem(t *testing.T) {
	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			launch: &Launch{
				Id: "id123",
			},
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}
		err := ti.Start()
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}
