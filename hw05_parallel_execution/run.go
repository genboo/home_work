package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrNoWorkers = errors.New("workers counts must be greater 0")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrNoWorkers
	}
	tch := make(chan Task, n)
	go func() {
		// кладем таски в канал, чтобы другие горутины могли их по одному забирать
		for _, t := range tasks {
			tch <- t
		}
		close(tch)
	}()
	var wg sync.WaitGroup
	var errCounter int32
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// читаем, пока канал не закроется
			for t := range tch {
				e := int(atomic.LoadInt32(&errCounter))
				if m != 0 && e > m {
					break
				}
				err := t()
				if err != nil {
					atomic.AddInt32(&errCounter, 1)
				}
			}
		}()
	}
	wg.Wait()
	if m != 0 && int(errCounter) > m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
