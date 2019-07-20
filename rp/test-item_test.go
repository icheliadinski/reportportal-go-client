package rp

import (
	"io/ioutil"
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
	t.Run("Successful start without parent id", func(t *testing.T) {
		okResponse := `{"id": "testid"}`
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/item", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			d, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(d), `"name":"item name"`)
			assert.Contains(t, string(d), `"description":"item description"`)
			assert.Contains(t, string(d), `"tags":["test","tag"]`)
			assert.Contains(t, string(d), `"type":"item type"`)
			assert.Contains(t, string(d), `"launch_id":"id123"`)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(okResponse))
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Name:        "item name",
			Description: "item description",
			Tags:        []string{"test", "tag"},
			Type:        "item type",
			launch: &Launch{
				Id: "id123",
			},
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}
		err := ti.Start()
		assert.NoError(t, err)
		assert.Equal(t, "testid", ti.Id)
	})

	t.Run("Valid url with parent id", func(t *testing.T) {
		okResponse := `{"id": "testid"}`
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/item/parent123", r.URL.Path)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(okResponse))
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Parent: &TestItem{
				Id: "parent123",
			},
			launch: &Launch{
				Id: "id123",
			},
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}

		err := ti.Start()
		assert.NoError(t, err)
	})

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

func TestFinishTestItem(t *testing.T) {
	t.Run("Successful finish", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/item/id123", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			d, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(d), `"status":"finish status"`)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Id: "id123",
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}

		err := ti.Finish("finish status")
		assert.NoError(t, err)
	})
	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}
		err := ti.Finish("")
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}
