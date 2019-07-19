package rp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestItem(t *testing.T) {
	l := &Launch{}
	ti := NewTestItem(l, "", "", "", nil, nil)
	assert.NotNil(t, ti)
}
