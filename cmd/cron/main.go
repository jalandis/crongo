package main

import (
	"errors"
	"log"
	"time"

	"github.com/jalandis/crongo/pkg/cron"
)

func slowWork(i int) {
	time.Sleep(time.Duration(i) * time.Second)
}

func main() {
	c, err := cron.Start([]cron.Job{{
		Name: "A - 1 second job every 5 seconds",
		Run:  func() { slowWork(1) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 5 * time.Second,
		},
	}, {
		Name: "B - 2 second job every 2 seconds",
		Run:  func() { slowWork(2) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 2 * time.Second,
		},
	}, {
		Name: "C - 2 second job every 1 seconds",
		Run:  func() { slowWork(2) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: time.Second,
		},
	}, {
		Name: "D - 1 second job every 1 seconds",
		Run:  func() { slowWork(1) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: time.Second,
		},
	}, {
		Name: "E - panic!",
		Run:  func() { panic(errors.New("unknown error")) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 4 * time.Second,
		},
	}})

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)
	c.Stop()
}
