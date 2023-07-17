package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errorCounter struct {
	mu        sync.Mutex
	val       int
	threshold int
}

func (counter *errorCounter) inc() {
	counter.mu.Lock()
	counter.val++
	counter.mu.Unlock()
}

func (counter *errorCounter) reachedThreshold() bool {
	counter.mu.Lock()
	thReached := counter.val >= counter.threshold
	counter.mu.Unlock()
	return thReached
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	wg.Add(n)

	ch := make(chan int)
	errorCnt := errorCounter{
		threshold: m,
		val:       0,
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
