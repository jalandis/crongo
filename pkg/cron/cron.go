package cron

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Schedule struct {
	Start    time.Time
	Interval time.Duration
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

func log(s string) {
	fmt.Printf("%s : %s\n", time.Now().Format("Mon Jan _2 15:04:05 2006"), s)
}

func getNextRunTime(s Schedule, t time.Time) time.Time {
	if s.Start.After(t) {
		return s.Start.Add(s.Interval)
	}

	return t.Add(s.Interval)
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

func Start(jobs []Job, ctx context.Context) (*Cron, error) {
	if len(jobs) == 0 {
		return nil, errors.New("at least one job is required")
	}

	startTime := time.Now()
	var q []timedJob
	for _, j := range jobs {
		q = append(q, timedJob{
			Job:     j,
			NextRun: getNextRunTime(j.Schedule, startTime),
		})
	}

	log(fmt.Sprintf("starting cron with %d jobs", len(jobs)))
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
							log(fmt.Sprintf("panic with job (%s) : %v", j.Name, r))
						}
					}()
					j.Run(cancelCtx)
				}(q[index].Job)
				q[index].NextRun = getNextRunTime(q[index].Job.Schedule, time.Now())
			case <-c.done:
				cancel()
				return
			}
		}
	}()

	return c, nil
}

func (c *Cron) Stop() {
	log("halting jobs")
	c.done <- struct{}{}
	c.wg.Wait()
}
