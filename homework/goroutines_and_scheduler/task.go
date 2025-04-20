package main

import (
	"container/heap"
)

type Task struct {
	Identifier int
	Priority   int
}

type TaskQueue struct {
	items  []Task
	id2Idx map[int]int
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		items:  []Task{},
		id2Idx: make(map[int]int),
	}
}

func (q *TaskQueue) Len() int {
	return len(q.items)
}

func (q *TaskQueue) Less(i, j int) bool {
	return q.items[i].Priority > q.items[j].Priority
}

func (q *TaskQueue) Swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
	q.id2Idx[q.items[i].Identifier] = i
	q.id2Idx[q.items[j].Identifier] = j
}

func (q *TaskQueue) Push(t any) {
	task := t.(Task)
	q.items = append(q.items, task)
	q.id2Idx[task.Identifier] = len(q.items)
}

func (q *TaskQueue) Pop() any {
	n := len(q.items)
	x := q.items[n-1]
	q.items = q.items[0 : n-1]
	delete(q.id2Idx, x.Identifier)
	return x
}

func (q *TaskQueue) UpdatePriority(id int, priority int) {
	idx, ok := q.id2Idx[id]
	if !ok {
		return
	}

	q.items[idx].Priority = priority
	heap.Fix(q, idx)
}
