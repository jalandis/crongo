package cron

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCron(t *testing.T) {
	called := make(chan bool)
	c, err := Start([]Job{{
		Name:     "testing",
		Run:      func() { called <- true },
		Schedule: Schedule{Interval: time.Millisecond},
	}})
	assert.NoError(t, err)

	select {
	case b := <-called:
		assert.True(t, b, "cron job called")
	case <-time.After(time.Second):
		assert.Fail(t, "timeout waiting for cron job to be called")
	}

	c.Stop()
}

func TestPanicCaught(t *testing.T) {
	called := make(chan bool)
	c, err := Start([]Job{{
		Name: "testing",
		Run: func() {
			called <- true
			panic("testing")
		},
		Schedule: Schedule{Interval: time.Millisecond},
	}})
	assert.NoError(t, err)

	select {
	case b := <-called:
		assert.True(t, b, "cron job called")
	case <-time.After(time.Second):
		assert.Fail(t, "timeout waiting for cron job to be called")
	}

	c.Stop()
}
