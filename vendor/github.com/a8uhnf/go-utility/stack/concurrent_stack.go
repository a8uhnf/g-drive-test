package stack

import (
	"sync"
)

const (
	maxedoutConcurrentStackError = "ConcurrentStack maxed out. Unable to push"
	emptyConcurrentStackError    = "ConcurrentStack empty. Unable to pop."
)

// ConcurrentStack struct keep tracks of top element, Concurrentstack size and lock for concurrent safe
type ConcurrentStack struct {
	top  *ConcurrentstackElement
	size int
	lock sync.Mutex
}

// ConcurrentstackElement represent stack element. Contains value and previous element.
type ConcurrentstackElement struct {
	value interface{}
	prev  *ConcurrentstackElement
}

// Len returns the length of Concurrentstack.
func (s *ConcurrentStack) Len() int {
	return s.size
}

// Top returns the top element of Concurrentstack.
func (s *ConcurrentStack) Top() interface{} {
	if s.Len() > 0 {
		return s.top.value
	}
	return nil
}

// Pop remove top element of Concurrentstack and returns it.
func (s *ConcurrentStack) Pop() interface{} {
	if s.Len() > 0 {
		s.lock.Lock()
		defer s.lock.Unlock()
		ret := s.top.value
		s.top = s.top.prev
		s.size--
		return ret
	}
	return nil
}

// Push add element to Concurrentstack
func (s *ConcurrentStack) Push(v interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.top = &ConcurrentstackElement{
		value: v,
		prev:  s.top,
	}
	s.size++
}
