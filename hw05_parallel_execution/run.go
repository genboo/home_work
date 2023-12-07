package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tch := make(chan Task, n)
	go func() {
		// кладем таски в канал, чтобы другие горутины могли их по одному забирать
		for _, t := range tasks {
			tch <- t
		}
		close(tch)
	}()
	var wg sync.WaitGroup
	var errCounter int
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// читаем, пока канал не закроется
			for t := range tch {
				if m != 0 && errCounter > m {
					break
				}
				err := t()
				if err != nil {
					errCounter++
				}
			}
		}(i)
	}
	wg.Wait()
	if m != 0 && errCounter > m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
