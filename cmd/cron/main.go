package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jalandis/crongo/pkg/cron"
)

func slowWork(i int, ctx context.Context) {
	fmt.Printf("Running job for %d seconds\n", i)
	select {
	case <-time.After(time.Duration(i) * time.Second):
	case <-ctx.Done():
		fmt.Println("Work canceled")
	}
}

func main() {
	c, err := cron.Start([]cron.Job{{
		Name: "A - 10 second job every 5",
		Run:  func(ctx context.Context) { slowWork(10, ctx) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 5 * time.Minute,
		},
	}, {
		Name: "B - 20 second job every 2",
		Run:  func(ctx context.Context) { slowWork(20, ctx) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 2 * time.Minute,
		},
	}, {
		Name: "C - 20 second job every 15",
		Run:  func(ctx context.Context) { slowWork(20, ctx) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 15 * time.Minute,
		},
	}, {
		Name: "D - 10 second job every 10",
		Run:  func(ctx context.Context) { slowWork(10, ctx) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 10 * time.Minute,
		},
	}, {
		Name: "E - panic!",
		Run:  func(ctx context.Context) { panic(errors.New("unknown error")) },
		Schedule: cron.Schedule{
			Start:    time.Now(),
			Interval: 4 * time.Minute,
		},
	}}, context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	<-sigCh
	fmt.Println("Signal received. Shutting down.")
	c.Stop()
}
