package lfqueue

import (
	"fmt"
	"sync"
	"testing"
)

func TestBoundedQueue_DeQueue(t *testing.T) {

	var wg sync.WaitGroup
	count := 100
	wg.Add(2)
	q := NewBoundedQueue(10)
	go func() {
		defer wg.Done()
		for  i := 0; i < 10000; i++{
			q.EnQueue(i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < count; i ++{
			v := q.DeQueue()
			if v == nil{

			}
			fmt.Printf("i is: %v v is: %v\n",i, v)

		}
	}()
	wg.Wait()
	if len(q.queue) != 0{
		t.Fatalf("queue size should be is:%v, cur is: %v", 0, len(q.queue))
	}
}