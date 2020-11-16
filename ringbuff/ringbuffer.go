package ringbuff

type RingBuffer struct {
	readPos  int
	writePos int
	buffer   []byte
	size     int
	isEmpty  bool
}

func New(size int) *RingBuffer {
	r := &RingBuffer{
		buffer: make([]byte, size),
		isEmpty: true
	}
	return r
}
