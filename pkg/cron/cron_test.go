package cron

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCron(t *testing.T) {
	called := make(chan bool)
	c, err := Start([]Job{{
		Name:     "testing",
		Run:      func(ctx context.Context) { called <- true },
		Schedule: Schedule{Interval: time.Millisecond},
	}}, context.Background())
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
		Run: func(ctx context.Context) {
			called <- true
			panic("testing")
		},
		Schedule: Schedule{Interval: time.Millisecond},
	}}, context.Background())
	assert.NoError(t, err)

	select {
	case b := <-called:
		assert.True(t, b, "cron job called")
	case <-time.After(time.Second):
		assert.Fail(t, "timeout waiting for cron job to be called")
	}

	c.Stop()
}

func TestCancelWork(t *testing.T) {
	called := make(chan bool)
	canceled := false
	c, err := Start([]Job{{
		Name: "testing",
		Run: func(ctx context.Context) {
			called <- true
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				canceled = true
			}
		},
		Schedule: Schedule{Interval: time.Millisecond},
	}}, context.Background())
	assert.NoError(t, err)

	select {
	case b := <-called:
		assert.True(t, b, "cron job called")
		c.Stop()
		assert.True(t, canceled, "cron job canceled")
	case <-time.After(time.Second):
		assert.Fail(t, "timeout waiting for cron job to be called")
	}
}
