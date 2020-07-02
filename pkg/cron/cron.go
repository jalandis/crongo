package cron

import (
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
	Run      func()
	Schedule Schedule
}

type cron struct {
	mx      sync.Mutex
	wg      sync.WaitGroup
	done    chan (struct{})
	running bool
}

func log(s string) {
	fmt.Printf("%s : %s\n", time.Now().Format("Mon Jan _2 15:04:05 2006"), s)
}

func getNextRun(s Schedule, t time.Time) time.Time {
	if s.Start.After(t) {
		return s.Start.Add(s.Interval)
	}

	return t.Add(s.Interval)
}

type runningJob struct {
	LastRun time.Time
	Job     Job
}

func nextJob(queue []runningJob) (int, time.Time) {
	result := 0
	nextRun := getNextRun(queue[0].Job.Schedule, queue[0].LastRun)
	for index, item := range queue {
		check := getNextRun(item.Job.Schedule, item.LastRun)
		if check.Before(nextRun) {
			result = index
			nextRun = check
		}
	}

	return result, nextRun
}

func Start(jobs []Job) (*cron, error) {
	if len(jobs) == 0 {
		return nil, errors.New("at least one job is required")
	}

	startTime := time.Now()
	var q []runningJob
	for _, j := range jobs {
		q = append(q, runningJob{
			Job:     j,
			LastRun: startTime,
		})
	}

	log("starting jobs")
	c := &cron{done: make(chan struct{})}
	go func() {
		for {
			index, nextRun := nextJob(q)
			select {
			case <-time.After(time.Until(nextRun)):
				c.wg.Add(1)
				go func(j Job) {
					defer func() {
						c.wg.Done()
						if r := recover(); r != nil {
							log(fmt.Sprintf("panic with job (%s) : %v", j.Name, r))
						}
					}()
					log(fmt.Sprintf("starting %s", j.Name))
					j.Run()
					log(fmt.Sprintf("finished %s", j.Name))
				}(q[index].Job)
				q[index].LastRun = time.Now()
			case <-c.done:
				log("done signaled")
				return
			}
		}
	}()

	return c, nil
}

func (c *cron) Stop() {
	log("halting jobs")
	c.done <- struct{}{}
	c.wg.Wait()
}
