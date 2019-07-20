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

func TestDashboard(t *testing.T) {
	t.Run("Successful result", func(t *testing.T) {
		okResponse := `[{"owner":"user","share": true,"id":"id123","name":"main","widgets":[{"widgetId":"wid123", "widgetSize":[12,7],"widgetPosition":[0,0]}]}]`
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/dashboard", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.Write([]byte(okResponse))
		})
		s := httptest.NewServer(h)

		c := &Client{
			Endpoint: s.URL,
			Project:  "test_project",
		}

		expected := &Dashboard{
			{
				Owner: "user",
				Share: true,
				Id:    "id123",
				Name:  "main",
				Widgets: []*Widget{
					{
						Id:       "wid123",
						Size:     []int{12, 7},
						Position: []int{0, 0},
					},
				},
			},
		}

		d, err := c.GetDashboard()
		assert.NoError(t, err)
		assert.NotNil(t, d)

		assert.Equal(t, expected, d)
	})

	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)

		c := &Client{
			Endpoint: s.URL,
			Project:  "test_project",
		}

		d, err := c.GetDashboard()
		assert.Nil(t, d)
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}
