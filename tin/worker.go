package tin

import (
	"time"
)

// Worker executes the task on intervals.
type Worker struct {
	task   func()
	ticker *time.Ticker
	stop   chan struct{}
}

// Stop stops the ticker and closes the channel.
func (s *Worker) Stop() {
	s.ticker.Stop()
	close(s.stop)
}

// NewWorker returns a tin.Worker.
func NewWorker(interval time.Duration, task func()) *Worker {
	go task() // Executes the task immediately.

	w := &Worker{
		task:   task,
		ticker: time.NewTicker(interval),
		stop:   make(chan struct{}, 1),
	}

	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.task()
			case <-w.stop:
				return
			}
		}
	}()

	return w
}
