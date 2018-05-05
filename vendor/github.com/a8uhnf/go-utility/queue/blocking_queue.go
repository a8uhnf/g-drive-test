package queue

import (
	"sync"
)

// BlockingQueue struct keeps track of top element, size, last element of BlockingQueue
type BlockingQueue struct {
	top            *BlockingQueueElement
	size           int
	maxSize        int
	popLock        sync.Mutex
	pushLock       sync.Mutex
	popBlock       chan int
	pushBlock      chan int
	pushBlockState bool
	popBlockState  bool
	last           *BlockingQueueElement
}

// BlockingQueueElement contains queue's element value, next, previous element
type BlockingQueueElement struct {
	value interface{}
	next  *BlockingQueueElement
	prev  *BlockingQueueElement
}

// NewBlockingQueue initialize blocking queue.
func NewBlockingQueue() *BlockingQueue {
	ret := &BlockingQueue{
		popBlock:  make(chan int),
		pushBlock: make(chan int),
	}
	return ret
}

// SetSize sets the maximum size of blocking queue
func (q *BlockingQueue) SetSize(sz int) {
	q.maxSize = sz
}

// Len returns the length of BlockingQueue
func (q *BlockingQueue) Len() int {
	return q.size
}

// Top function returns top elemetn of BlockingQueue
func (q *BlockingQueue) Top() interface{} {
	if q.Len() > 0 {
		return q.top.value
	}
	return nil
}

// Push function insert element at the last of BlockingQueue
func (q *BlockingQueue) Push(v interface{}) {
	q.pushLock.Lock()
	defer q.pushLock.Unlock()
	if q.Len() >= q.maxSize {
		q.pushBlockState = true
		q.pushBlock <- 1
		q.pushBlockState = false
	}
	if q.Len() == 0 {
		q.top = &BlockingQueueElement{
			value: v,
			next:  nil,
		}
		q.last = q.top
		q.size++
		return
	}
	q.size++
	tmp := q.last
	q.last = &BlockingQueueElement{
		next:  q.last,
		value: v,
	}
	tmp.prev = q.last
	if q.popBlockState {
		<-q.popBlock
	}
}

// Pop function remove top element from the BlockingQueue
func (q *BlockingQueue) Pop() interface{} {
	q.popLock.Lock()
	defer q.popLock.Unlock()
	if q.Len() == 0 {
		q.popBlockState = true
		q.popBlock <- 1
		q.popBlockState = false
	}
	ret := q.top.value
	q.top.next = nil
	q.top = q.top.prev
	q.size--
	if q.pushBlockState {
		<-q.pushBlock
	}
	return ret
}
