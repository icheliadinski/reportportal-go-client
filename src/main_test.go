package main

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	go func() {
		time.Sleep(5000 * time.Millisecond)
		p, e := os.FindProcess(os.Getpid())
		e = p.Kill()
		require.Nil(t, e)
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		st := time.Now()
		main()
		assert.True(t, time.Since(st).Seconds() >= 5, "should take about 5s")
		wg.Done()
	}()
}
