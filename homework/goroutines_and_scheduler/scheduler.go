package main

import (
	"container/heap"
)

type Scheduler struct {
	tasks *TaskQueue
}

func NewScheduler() Scheduler {
	return Scheduler{
		tasks: NewTaskQueue(),
	}
}

func (s *Scheduler) AddTask(task Task) {
	heap.Push(s.tasks, task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	s.tasks.UpdatePriority(taskID, newPriority)
}

func (s *Scheduler) GetTask() Task {
	if s.tasks.Len() > 0 {
		return heap.Pop(s.tasks).(Task)
	}
	return Task{}
}
