package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type WorkerPool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		tasks: make(chan func(), workersNumber),
	}

	wp.wg.Add(workersNumber)

	// run workers
	for i := 0; i < workersNumber; i++ {
		go func() {
			defer wp.wg.Done()
			for t := range wp.tasks {
				t()
			}
		}()
	}

	return wp
}

// Return an error if the pool is full
func (wp *WorkerPool) AddTask(task func()) (err error) {
	if len(wp.tasks) >= cap(wp.tasks) {
		err = errors.New("worker pull is full")
		go func() {
			wp.tasks <- task
		}()
	} else {
		wp.tasks <- task
	}

	return
}

// Shutdown all workers and wait for all
// tasks in the pool to complete
func (wp *WorkerPool) Shutdown() {
	close(wp.tasks)
	wp.wg.Wait()
}

func TestWorkerPool(t *testing.T) {
	var counter atomic.Int32
	task := func() {
		time.Sleep(time.Millisecond * 500)
		counter.Add(1)
	}

	pool := NewWorkerPool(2)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(2), counter.Load())

	time.Sleep(time.Millisecond * 600)
	assert.Equal(t, int32(3), counter.Load())

	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	_ = pool.AddTask(task)
	pool.Shutdown() // wait tasks

	assert.Equal(t, int32(6), counter.Load())
}
