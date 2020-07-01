package main

import (
	"time"

	"github.com/jalandis/crongo/pkg/cron"
)

func slowWork(i int) {
	time.Sleep(time.Duration(i) * time.Second)
}

func main() {
	c := cron.Init()
	c.AddJob(
		"A - 1 second job every 5 seconds",
		func() { slowWork(1) },
		cron.Schedule(5*time.Second),
	)
	c.Start()

	time.Sleep(1 * time.Second)
	c.AddJob(
		"B - 2 second job every 2 seconds",
		func() { slowWork(2) },
		cron.Schedule(2*time.Second),
	)

	time.Sleep(1 * time.Second)
	c.AddJob(
		"C - 2 second job every 1 seconds",
		func() { slowWork(2) },
		cron.Schedule(1*time.Second),
	)
	time.Sleep(4 * time.Second)
	c.Stop()
}
