package queue

import (
	"sync"
)

// ConcurrentQueue struct keeps track of top element, size, last element of ConcurrentQueue
type ConcurrentQueue struct {
	top  *ConcurrentQueueElement
	size int
	lock sync.Mutex
	last *ConcurrentQueueElement
}

// ConcurrentQueueElement contains queue's element value, next, previous element
type ConcurrentQueueElement struct {
	value interface{}
	next  *ConcurrentQueueElement
	prev  *ConcurrentQueueElement
}

// Len returns the length of ConcurrentQueue
func (q *ConcurrentQueue) Len() int {
	return q.size
}

// Top function returns top elemetn of ConcurrentQueue
func (q *ConcurrentQueue) Top() interface{} {
	if q.Len() > 0 {
		return q.top.value
	}
	return nil
}

// Push function insert element at the last of ConcurrentQueue
func (q *ConcurrentQueue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Len() == 0 {
		q.top = &ConcurrentQueueElement{
			value: v,
			next:  nil,
		}
		q.last = q.top
		q.size++
		return
	}
	q.size++
	tmp := q.last
	q.last = &ConcurrentQueueElement{
		next:  q.last,
		value: v,
	}
	tmp.prev = q.last
}

// Pop function remove top element from the ConcurrentQueue
func (q *ConcurrentQueue) Pop() interface{} {
	if q.Len() > 0 {
		q.lock.Lock()
		defer q.lock.Unlock()
		ret := q.top.value
		q.top.next = nil
		q.top = q.top.prev
		q.size--
		return ret
	}
	return nil
}
