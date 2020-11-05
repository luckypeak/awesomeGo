package lfqueue

type Queue interface {
	EnQueue(v interface{})
	DeQueue() interface{}
}
