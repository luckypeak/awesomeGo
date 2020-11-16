package lfqueue

import "sync"

type SliceQueue struct {
	mu sync.Mutex
	queue []interface{}
}

func NewSliceQueue()  (q *SliceQueue){
  q = new(SliceQueue)
  return
}


func (q *SliceQueue) EnQueue(v interface{})  {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, v)
}

func (q *SliceQueue)DeQueue()(v interface{})  {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return
	}
	v = q.queue[0]
	q.queue = q.queue[1:]
	return
}