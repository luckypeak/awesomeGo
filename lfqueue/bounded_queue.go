package lfqueue

import "sync"

type BoundedQueue struct {
	queue []interface{}
	cond *sync.Cond
	capacity int
	len int
}

func NewBoundedQueue(capacity int) (q * BoundedQueue)  {
	q = new(BoundedQueue)
	q.capacity = capacity
	q.cond = sync.NewCond(&sync.Mutex{})
	return
}

func (q *BoundedQueue) EnQueue(v interface{}) {
	q.cond.L.Lock()
	for q.len == q.capacity {
		q.cond.Wait()
	}
	defer q.cond.L.Unlock()
	q.len ++
	q.queue = append(q.queue, v)
	q.cond.Broadcast()

}

func (q *BoundedQueue)DeQueue() (v interface{})  {
	q.cond.L.Lock()
	for q.len == 0{
		q.cond.Wait()
	}
	defer q.cond.L.Unlock()
	v = q.queue[0]
	q.len--
	q.queue = q.queue[1:len(q.queue)]
	q.cond.Broadcast()
	return
}