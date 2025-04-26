package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v .
func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	// Приоритет поменялся, соот-но просто сравнивать с task1 нельзя
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
