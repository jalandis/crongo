package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCron(t *testing.T) {
	called := make(chan bool)
	c := Init()
	c.AddJob("testing", func() { called <- true }, Schedule(0))
	c.Start()

	select {
	case b := <-called:
		assert.True(t, b, "cron job called")
	case <-time.After(time.Second):
		assert.Fail(t, "timeout waiting for cron job to be called")
	}
}
