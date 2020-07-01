package cron

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type LoggingLevel int

const (
	DebugLogging LoggingLevel = iota
)

type Schedule time.Duration

type Job struct {
	Name     string
	Run      func()
	Schedule Schedule
	NextRun  time.Time
}

type cron struct {
	mx       sync.Mutex
	wg       sync.WaitGroup
	done     chan (struct{})
	jobs     Jobs
	logLevel LoggingLevel
}

type Jobs []Job

func (c Jobs) Len() int           { return len(c) }
func (c Jobs) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Jobs) Less(i, j int) bool { return c[i].NextRun.Before(c[j].NextRun) }

func log(s string) {
	fmt.Printf("%s : %s\n", time.Now().Format("Mon Jan _2 15:04:05 2006"), s)
}

func getNextRun(s Schedule) time.Time {
	return time.Now().Add(time.Duration(s))
}

func runJob(c *cron) {
	c.mx.Lock()
	defer c.mx.Unlock()

	if len(c.jobs) == 0 {
		return
	}

	c.wg.Add(1)
	go func(j Job) {
		defer c.wg.Done()
		log(fmt.Sprintf("starting %s", j.Name))
		j.Run()
		log(fmt.Sprintf("finished %s", j.Name))
	}(c.jobs[0])

	c.jobs[0].NextRun = getNextRun(c.jobs[0].Schedule)
	sort.Sort(c.jobs)
}

func nextJob(c *cron) time.Duration {
	c.mx.Lock()
	defer c.mx.Unlock()

	if len(c.jobs) == 0 {
		return time.Duration(math.MaxInt64)
	}

	return c.jobs[0].NextRun.Sub(time.Now())
}

func Init() *cron {
	return &cron{done: make(chan struct{})}
}

func (c *cron) AddJob(n string, r func(), s Schedule) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.Stop()
	c.jobs = append(c.jobs, Job{
		Name:     n,
		Run:      r,
		Schedule: s,
		NextRun:  getNextRun(s),
	})
	sort.Sort(c.jobs)
	c.Start()
}

func (c *cron) Start() {
	go func() {
		for {
			select {
			case <-time.After(nextJob(c)):
				runJob(c)
			case <-c.done:
				log("done signaled")
				return
			}
		}
	}()
}

func (c *cron) Stop() {
	log("halting jobs")
	go func() { c.done <- struct{}{} }()
	c.wg.Wait()
}
