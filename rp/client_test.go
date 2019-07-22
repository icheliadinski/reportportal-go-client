package rp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	t.Run("Successful check", func(t *testing.T) {
		h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			res.Write([]byte("response"))
		})
		s := httptest.NewServer(h)
		defer s.Close()

		c := NewClient(s.URL, "test_project", "1234", 1)
		err := c.CheckConnect()
		assert.NoError(t, err)
	})

	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)
		defer s.Close()

		c := &Client{
			Endpoint: s.URL,
		}
		err := c.CheckConnect()
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
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
		defer s.Close()

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
		defer s.Close()

		c := &Client{
			Endpoint: s.URL,
		}

		d, err := c.GetDashboard()
		assert.Nil(t, d)
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}

func TestActivity(t *testing.T) {
	t.Run("Successful result", func(t *testing.T) {
		okResponse := `{
			"content": [
				{
					"actionType": "test activity",
					"activityId": "test act id",
					"history": [
						{
							"field": "test field",
							"newValue": "new value",
							"oldValue": "old value"
						}
					],
					"lastModifiedDate": "2019-07-22T10:10:10.000Z",
					"loggedObjectRef": "object ref",
					"objectName": "object name",
					"objectType": "object type",
					"projectRef": "project ref",
					"userRef": "user ref"
				}
			],
			"page": {
				"number": 1,
				"size": 2,
				"totalElements": 3,
				"totalPages": 4
			}
		}`
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/test_project/activity", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			w.Write([]byte(okResponse))
		})
		s := httptest.NewServer(h)
		defer s.Close()

		c := &Client{
			Endpoint: s.URL,
			Project:  "test_project",
		}

		expected := &Activity{
			Content: []*ActivityContent{
				{
					ActionType: "test activity",
					ActivityId: "test act id",
					History: []*ActivityHistory{
						{
							Field:    "test field",
							NewValue: "new value",
							OldValue: "old value",
						},
					},
					LastModifiedDate: time.Date(2019, time.July, 22, 10, 10, 10, 0, time.UTC),
					LoggedObjectRef:  "object ref",
					ObjectName:       "object name",
					ObjectType:       "object type",
					ProjectRef:       "project ref",
					UserRef:          "user ref",
				},
			},
			Page: &ActivityPage{
				Number:        1,
				Size:          2,
				TotalElements: 3,
				TotalPages:    4,
			},
		}

		a, err := c.GetActivity()
		assert.NoError(t, err)
		assert.NotNil(t, a)

		assert.Equal(t, expected, a)
	})

	t.Run("Wrong status code", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		s := httptest.NewServer(h)
		defer s.Close()

		c := &Client{
			Endpoint: s.URL,
		}

		d, err := c.GetActivity()
		assert.Nil(t, d)
		assert.EqualError(t, err, "failed with status 500 Internal Server Error")
	})
}
