package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Scheduler struct {
	jobs []*job
}

type job struct {
	interval time.Duration
	f        func() error
}

func newJob(interval time.Duration, f func() error) *job {
	return &job{interval: interval, f: f}
}

func (s *Scheduler) SetupJob(interval time.Duration, f func() error) {
	s.jobs = append(s.jobs, newJob(interval, f))
}

func (s *Scheduler) Run(ctx context.Context) error {
	for _, j := range s.jobs {
		err := j.f()
		if err != nil {
			return fmt.Errorf("can't run scheduler")
		}
		go j.run()
	}
	return nil
}

func (j *job) run() error {
	timer := time.NewTicker(j.interval)
	for {
		select {
		case <-timer.C:
			err := j.f()
			if err != nil {
				log.Println("can't run update func")
			}
		}
	}
}
