package lfqueue

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	val interface{}
	next unsafe.Pointer
}
type LockFreeQueue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewLockFreeQueue() (q *LockFreeQueue)  {
	q = new(LockFreeQueue)
	n := unsafe.Pointer(&node{})
	q.head = n
	q.tail = n
	return
}

func (q *LockFreeQueue) EnQueue(val interface{}) {

	n := unsafe.Pointer(&node{val:val, next:nil})
	var tail unsafe.Pointer
	for {
		tail = q.tail
		next := ((*node)(tail)).next
		if tail == q.tail {
			if next == nil{
				if atomic.CompareAndSwapPointer(& ((*node)(tail)).next, next, n){
					atomic.CompareAndSwapPointer(&q.tail, tail,n)
					return
				}
			}else{
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}



	}


}

func (q *LockFreeQueue)DeQueue() (val interface{}, success bool)  {
    for {
    	head := q.head
    	tail := q.tail
    	next := ((*node)(head)).next
    	if head == q.head {
    		if head == tail{
    			if next == nil{
    				return nil, false
				}
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}else{
				val = ((*node)(next)).val
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					break
				}
			}
		}
	}
	return
}