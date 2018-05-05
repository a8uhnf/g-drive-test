package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func regularQueueOps(t *testing.T) {
	q := new(ConcurrentQueue)
	assert.Equal(t, q.Len(), 0, fmt.Sprint("Queue length should be zero."))
	q.Push(2)
	assert.Equal(t, q.Top().(int), 2, fmt.Sprintf("Top element should be 2. Instead found %v", q.Top()))
	q.Push(3)
	assert.Equal(t, q.Top().(int), 2, fmt.Sprintf("Top element should be 2. Instead found %v", q.Top()))
	top := q.Pop()
	assert.Equal(t, top.(int), 2, fmt.Sprintf("Poped element should be 2. Instead found %v", top))
	top = q.Pop()
	assert.Equal(t, top, 3, fmt.Sprintf("Poped element should be 3, Instead found %v", top))
	assert.Equal(t, q.Len(), 0, fmt.Sprint("Queue length should be zero."))
}

func regularQueueOpsBenchmark() error {
	q := new(ConcurrentQueue)
	if q.Len() != 0 {
		return fmt.Errorf("Queue length should be zero")
	}
	q.Push(2)
	if q.Top().(int) != 2 {
		return fmt.Errorf("Top element should be 2. Instead found %v", q.Top())
	}
	q.Push(3)
	if q.Top() != 2 {
		return fmt.Errorf("Top element should be 2. Instead found %v", q.Top())
	}
	top := q.Pop()
	if top != 2 {
		return fmt.Errorf("Poped element should be 2. Instead found %v", top)
	}
	top = q.Pop()
	if top != 3 {
		return fmt.Errorf("Poped element should be 3, Instead found %v", top)
	}
	if q.Len() != 0 {
		return fmt.Errorf("Queue length should be zero")
	}
	return nil
}

func TestQueue(t *testing.T) {
	regularQueueOps(t)
}

func BenchmarkQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := regularQueueOpsBenchmark()
		assert.Equal(b, err, nil, fmt.Sprintf("Error should be nil"))
	}
}
