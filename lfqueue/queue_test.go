package lfqueue

import "testing"

func TestQueue(t *testing.T)  {
	count := 100
	q := NewSliceQueue()
	for i := 0; i < count; i++{
		q.EnQueue(i)
	}
	for i:= 0; i < count; i++{
		v := q.DeQueue()
		if v == nil{
			t.Fatalf("nil except is: %v", i)
		}
		if v.(int) != i{
			t.Fatalf("expect is %d, %d", i, v)
		}
	}
}