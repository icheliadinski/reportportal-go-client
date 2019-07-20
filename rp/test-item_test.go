package rp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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

			rx, _ := regexp.Compile(`\{\"name\"\:\"item name\"\,\"description\"\:\"item description\"\,\"tags\"\:\[\"test\"\,\"tag\"\]\,\"start\_time\"\:\d+\,\"launch\_id\"\:\"id123\"\,\"type\"\:\"item type\"\,\"parameters\"\:null\}`)
			assert.Regexp(t, rx, string(d))

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

func TestLogTestItem(t *testing.T) {
	t.Run("Successful write without attachment", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/log", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			d, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)

			rx, _ := regexp.Compile(`\{\"item\_id\"\:\"item id\"\,\"message\"\:\"log message\"\,\"level\"\:\"log level\"\,\"time\"\:\d+\}`)
			assert.Regexp(t, rx, string(d))

			w.WriteHeader(http.StatusCreated)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Id: "item id",
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}

		err := ti.Log("log message", "log level", nil)
		assert.NoError(t, err)
	})

	t.Run("Successful write with attachment", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/log", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data; boundary=")

			d, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(d), `Content-Disposition: form-data; name="json_request_part"`)
			assert.Contains(t, string(d), `Content-Type: application/json`)

			rx, _ := regexp.Compile(`\[\{\"file"\:\{\"name":\"test\-text\.txt"},\"item_id\"\:\"item id\"\,\"level\"\:\"log level\"\,\"message\"\:\"log message\"\,\"time\"\:\d+\}\]`)
			assert.Regexp(t, rx, string(d))

			assert.Contains(t, string(d), `Content-Disposition: form-data; name="file"; filename="test-text.txt"`)
			assert.Contains(t, string(d), `Content-Type: text/plain`)
			assert.Contains(t, string(d), `test text in file`)
			w.WriteHeader(http.StatusCreated)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Id: "item id",
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}

		err := ti.Log("log message", "log level", &Attachment{
			Name:     "test-text.txt",
			MimeType: "text/plain",
			Data:     strings.NewReader("test text in file"),
		})
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

		err := ti.Log("", "", nil)
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}

func TestUpdateTestItem(t *testing.T) {
	t.Run("Successful update", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/item/id123/update", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Id: "id123",
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}

		err := ti.Update("new description", []string{"new", "tags"})
		assert.NoError(t, err)
		assert.Equal(t, "new description", ti.Description)
		assert.Equal(t, []string{"new", "tags"}, ti.Tags)
	})

	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)

		ti := &TestItem{
			Id: "id123",
			client: &Client{
				Endpoint: s.URL,
				Project:  "test_project",
			},
		}
		err := ti.Update("", nil)
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}
