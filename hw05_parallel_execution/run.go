package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorCounter struct {
	mu        sync.RWMutex
	val       int
	threshold int
}

func (counter *errorCounter) inc() {
	defer counter.mu.Unlock()
	counter.mu.Lock()
	counter.val++
}

func (counter *errorCounter) reachedThreshold() bool {
	defer counter.mu.RUnlock()
	counter.mu.RLock()
	return counter.val >= counter.threshold
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)

	ch := make(chan int)
	errorCnt := errorCounter{
		threshold: m,
	}

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for taskInd := range ch {
				if err := tasks[taskInd](); err != nil {
					errorCnt.inc()
				}
			}
		}()
	}

	for i := range tasks {
		if errorCnt.reachedThreshold() {
			break
		}
		ch <- i
	}

	close(ch)
	wg.Wait()

	if errorCnt.reachedThreshold() {
		return ErrErrorsLimitExceeded
	}
	return nil
}
