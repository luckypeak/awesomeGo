package lfqueue

import "testing"

func TestQueue(t *testing.T)  {
	count := 100
	q := NewLockFreeQueue()
	for i := 0; i < count; i++{
		q.EnQueue(i)
	}
	for i:= 0; i < count; i++{
		v, _ := q.DeQueue()
		if v == nil{
			t.Fatalf("nil")
		}
		if v.(int) != i{
			t.Fatalf("expect is %d, %d", i, v)
		}
	}
}