package rp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewClient("rp.epam.com", "test_project", "1234")

	assert.NotNil(t, c)
	assert.Equal(t, "https://rp.epam.com/api/v1", c.Endpoint)
}
