package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func regularStackOps(t *testing.T) {
	s := new(ConcurrentStack)
	assert.Equal(t, s.Len(), 0, "Stack length should be 0")
	s.Push(2)
	assert.Equal(t, s.Top().(int), 2, fmt.Sprintf("Top element should be 2. Instead found %d", s.Top().(int)))
	s.Push(3.2)
	assert.Equal(t, s.Top().(float64), 3.2, fmt.Sprintf("Top element should be 3.2. Instead found %f", s.Top().(float64)))
	p := s.Pop()
	assert.Equal(t, p, 3.2, fmt.Errorf("Poped element should be 3.2. Instead found %f", p))
	p = s.Pop()
	assert.Equal(t, p, 2, fmt.Errorf("Poped element should be 2. Instead found %f", p))
	assert.Equal(t, s.Len(), 0, fmt.Sprint("Stack size should be zero."))
}

func TestStack(t *testing.T) {
	regularStackOps(t)
}

func regularStackOpsBenchmark() error {
	s := new(ConcurrentStack)
	if s.Len() != 0 {
		return fmt.Errorf("Stack length should be 0")
	}
	s.Push(2)
	if s.Top().(int) != 2 {
		return fmt.Errorf("Top element should be 2. Instead found %d", s.Top().(int))
	}
	s.Push(3.2)
	if int(s.Top().(float64)) != 3 {
		return fmt.Errorf("Poped element should be 3.2. Instead found %f", s.Top().(float64))
	}
	return nil
}

func BenchmarkStackRegularOps(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := regularStackOpsBenchmark()
		assert.Equal(b, err, nil, fmt.Sprintf("Error should be nil"))
	}
}
