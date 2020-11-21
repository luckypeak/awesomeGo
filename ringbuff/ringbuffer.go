package ringbuff

import (
	"errors"
	"sync"
)

var (
	Empty         = errors.New("ringbuffer is empty")
	Full          = errors.New("ringbuffer is full")
	NoEnoughSpace = errors.New("ringbuffer no enough space")
)

type RingBuffer struct {
	r       int
	w       int
	buf     []byte
	size    int
	isEmpty bool
	mu      sync.Mutex
}

func New(size int) *RingBuffer {
	r := &RingBuffer{
		buf:     make([]byte, size),
		isEmpty: true,
		size:    size,
	}
	return r
}

func (r *RingBuffer) Read(p []byte) (n int, err error) {
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	r.mu.Lock()
	if r.r == r.w && r.isEmpty {
		r.mu.Unlock()
		return 0, Empty
	}
	length := len(p)
	if r.r < r.w {
		n = r.w - r.r
		if n > length {
			n = length
		} else {
			r.isEmpty = true
		}
		copy(p, r.buf[r.r:r.r+n])
		r.r += n
		r.mu.Unlock()
		return
	}
	n = r.size - r.r + r.w
	if n > length {
		n = length
	} else {
		r.isEmpty = true
	}

	if r.r+n < r.size {
		copy(p, r.buf[r.r:r.r+n])
		r.r += n
		r.mu.Unlock()
		return
	} else {
		r1 := r.size - r.r
		copy(p, r.buf[r.r:r.size])
		r2 := n - r1
		copy(p[r1:], r.buf[0:r2])
		r.r = r2
	}
	r.mu.Unlock()
	return
}

func (r *RingBuffer) Write(p []byte) (err error) {
	if p == nil || len(p) == 0 {
		return nil
	}
	r.mu.Lock()
	if r.r == r.w && !r.isEmpty {
		r.mu.Unlock()
		return Full
	}
	length := len(p)
	var avail int
	if r.w >= r.r {
		avail = r.size - r.w + r.r
	} else {
		avail = r.r - r.w
	}
	if avail < length {
		r.mu.Unlock()
		return NoEnoughSpace
	}
	r.isEmpty = false
	if r.r <= r.w {
		w1 := r.size - r.w
		if w1 >= length {
			copy(r.buf[r.w:], p)
			r.w += length
		} else {
			copy(r.buf[r.w:], p[0:w1])
			copy(r.buf[0:], p[w1:])
			r.w = length - w1
		}
	} else {
		copy(r.buf[r.w:], p)
		r.w += length
	}

	if r.w == r.size {
		r.w = 0
	}
	r.mu.Unlock()
	return
}

func (r *RingBuffer) IsEmpty() (isEmpty bool) {
	return r.isEmpty
}

func (r *RingBuffer) IsFull() (isFull bool) {
	r.mu.Lock()
	isFull = r.r == r.w && !r.isEmpty
	r.mu.Unlock()
	return
}

func (r *RingBuffer) Length() (length int) {
	r.mu.Lock()
	if r.isEmpty {
		r.mu.Unlock()
		return 0
	}
	if r.w > r.r {
		length = r.w - r.r
	} else {
		length = r.size - r.r + r.w
	}
	r.mu.Unlock()
	return
}

func (r *RingBuffer) Free() (length int) {
	r.mu.Lock()
	if r.isEmpty {
		r.mu.Unlock()
		return r.size
	}
	if r.w > r.r {
		length = r.size - r.w + r.r
	} else {
		length = r.r - r.w
	}

	r.mu.Unlock()
	return
}

func (r *RingBuffer) Bytes() (buf []byte) {
	r.mu.Lock()
	if r.isEmpty {
		r.mu.Unlock()
		return nil
	}

	if r.w > r.r {
		buf = make([]byte, r.w-r.r)
		copy(buf, r.buf[r.r:r.w])
	} else {
		buf = make([]byte, r.size-r.r+r.w)
		r1 := r.size - r.r
		copy(buf, r.buf[r.r:])
		copy(buf[r1:], r.buf[0:r.w])
	}
	r.mu.Unlock()
	return
}

func (r *RingBuffer) Reset() {
	r.mu.Lock()
	r.r = 0
	r.w = 0
	r.isEmpty = true
	r.mu.Unlock()
	return
}
