package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jalandis/crongo/pkg/cron"
)

func slowWork(i int) {
	fmt.Printf("Running job for %d seconds\n", i)
	time.Sleep(time.Duration(i) * time.Second)
}

func main() {
	c, err := cron.Start([]cron.Job{{
		Name: "A - 10 second job every 5 mintes",
		Run:  func() { slowWork(10) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 5 * time.Minute,
		},
	}, {
		Name: "B - 20 second job every 2 minutes",
		Run:  func() { slowWork(20) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 2 * time.Minute,
		},
	}, {
		Name: "C - 20 second job every 15 minutes",
		Run:  func() { slowWork(20) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 15 * time.Second,
		},
	}, {
		Name: "D - 10 second job every 10 minutes",
		Run:  func() { slowWork(10) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 10 * time.Minute,
		},
	}, {
		Name: "E - panic!",
		Run:  func() { panic(errors.New("unknown error")) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 4 * time.Minute,
		},
	}})
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	fmt.Println("Signal received. Shutting down.")
	c.Stop()
}
