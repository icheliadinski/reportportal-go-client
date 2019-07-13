package rp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type params struct {
	Endpoint string
	Version  int
}

func TestNewClient(t *testing.T) {
	c := NewClient("", "", "", 0)
	assert.NotNil(t, c)

	var endpoints = []struct {
		client   *params
		expected string
	}{
		{&params{"rp.epam.com", 0}, "https://rp.epam.com/api/v1"},
		{&params{"rp.epam.com", -1}, "https://rp.epam.com/api/v1"},
		{&params{"rp.epam.com", 1}, "https://rp.epam.com/api/v1"},
		{&params{"rp.epam.com", 2}, "https://rp.epam.com/api/v2"},
		{&params{"https://rp.epam.com", 1}, "https://rp.epam.com/api/v1"},
		{&params{"https://rp.epam.com/", 1}, "https://rp.epam.com/api/v1"},
		{&params{"rp.epam.com/api/v1", 1}, "https://rp.epam.com/api/v1"},
	}

	for _, tt := range endpoints {
		actual := NewClient(tt.client.Endpoint, "", "", tt.client.Version)
		if actual.Endpoint != tt.expected {
			t.Errorf(`NewClient(%s, "", "", %d): expected %s, actual %s`, tt.client.Endpoint, tt.client.Version, tt.expected, actual.Endpoint)
		}
	}
}
