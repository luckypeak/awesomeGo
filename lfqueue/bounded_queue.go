package lfqueue

import "sync"

type BoundedQueue struct {
	queue []interface{}
	mu sync.Mutex
	capacity int
}

func NewBoundedQueue(capacity int) (q * BoundedQueue)  {
	q = new(BoundedQueue)
	q.capacity = capacity
	return
}

func (q *BoundedQueue) EnQueue(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.queue) > q.capacity  {

	}
	q.queue = append(q.queue, v)

}

func (q *BoundedQueue)DeQueue() (v interface{})  {
	q.mu.Lock()
	defer q.mu.Unlock()
	v = q.queue[0]
	q.queue = q.queue[1:len(q.queue)]
	return
}