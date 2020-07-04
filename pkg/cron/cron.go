package cron

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Schedule interface {
	NextRunTime(now time.Time) time.Time
}

type Job struct {
	Name     string
	Run      func(context.Context)
	Schedule Schedule
}

type Cron struct {
	mx      sync.Mutex
	wg      sync.WaitGroup
	done    chan (struct{})
	running bool
}

type timedJob struct {
	NextRun time.Time
	Job     Job
}

type ConstantInterval struct {
	Start    time.Time
	Interval time.Duration
}

func (s ConstantInterval) NextRunTime(now time.Time) time.Time {
	if s.Start.After(now) {
		return s.Start.Add(s.Interval)
	}

	return now.Add(s.Interval)
}

func nextJob(queue []timedJob) int {
	result := 0
	for index, item := range queue {
		if item.NextRun.Before(queue[result].NextRun) {
			result = index
		}
	}

	return result
}

func Start(ctx context.Context, jobs []Job) (*Cron, error) {
	if len(jobs) == 0 {
		return nil, errors.New("at least one job is required")
	}

	startTime := time.Now()
	var q []timedJob
	for _, j := range jobs {
		q = append(q, timedJob{
			Job:     j,
			NextRun: j.Schedule.NextRunTime(startTime),
		})
	}

	fmt.Printf("starting cron with %d jobs\n", len(jobs))
	cancelCtx, cancel := context.WithCancel(ctx)
	c := &Cron{done: make(chan struct{})}
	go func() {
		for {
			index := nextJob(q)
			select {
			case <-time.After(time.Until(q[index].NextRun)):
				c.wg.Add(1)
				go func(j Job) {
					defer func() {
						c.wg.Done()
						if r := recover(); r != nil {
							fmt.Printf("panic with job (%s) : %v\n", j.Name, r)
						}
					}()
					j.Run(cancelCtx)
				}(q[index].Job)
				q[index].NextRun = q[index].Job.Schedule.NextRunTime(time.Now())
			case <-c.done:
				cancel()
				return
			}
		}
	}()

	return c, nil
}

func (c *Cron) Stop() {
	fmt.Println("halting jobs")
	c.done <- struct{}{}
	c.wg.Wait()
}
