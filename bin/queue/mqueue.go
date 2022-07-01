package queue

import (
	"sync"
)

type element struct {
	next  *element
	value []byte
}

func newElement(v []byte) *element {
	return &element{value: v}
}

type queue struct {
	head *element
	tail *element
	l    *sync.RWMutex
	size int64
	ch   chan bool
}

func NewQueue() *queue {
	return &queue{head: nil, tail: nil, l: new(sync.RWMutex), size: 0, ch: make(chan bool, 100)}
}

func (q *queue) Size() int64 {
	q.l.RLock()
	defer q.l.RUnlock()
	return q.size
}

// Push Insert message to the end of the queue
func (q *queue) Push(v []byte) {
	e := newElement(v)
	q.l.Lock()
	defer q.l.Unlock()
	if q.size == 0 {
		q.head = e
		q.tail = e
	} else {
		q.tail.next = e
		q.tail = e
	}
	q.ch <- true
	q.size++
}

// Take get header message
func (q *queue) Take() []byte {
	q.l.RLocker()
	defer q.l.RUnlock()
	<-q.ch
	return q.head.value
}

func (q *queue) DeleteHead() {
	q.l.Lock()
	defer q.l.Unlock()
	q.head = q.head.next
	q.size--
	if q.size == 0 {
		q.tail = nil
	}
}
