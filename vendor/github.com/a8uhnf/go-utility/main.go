package main

import (
	"fmt"
	"time"

	"github.com/a8uhnf/go-utility/queue"
	"github.com/a8uhnf/go-utility/stack"
)

func main() {
	fmt.Println("Hello World!!!")
	s := new(stack.ConcurrentStack)
	s.Push(2)
	s.Push(3)
	fmt.Println(s.Pop())
	fmt.Println(s.Top())

	fmt.Println("----------------------")

	q := new(queue.ConcurrentQueue)
	q.Push(2)
	q.Push(3)
	fmt.Println(q.Pop())
	fmt.Println(q.Top())
	fmt.Println("----------------------")

	// blocking stack example
	bs := stack.NewBlockingStack()
	bs.SetSize(2)
	bs.Push(2)
	bs.Push(3)
	go bs.Push(6)
	fmt.Println(bs.Top())
	var chk int
	go func() {
		chk = bs.Pop().(int)
		fmt.Println("***", chk)
	}()
	go bs.Push(4)
	go func() {
		chk = bs.Pop().(int)
		fmt.Println("***", chk)
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("xxx", bs.Top())
	fmt.Println("xxx", bs.Len())

	fmt.Println("----------------------: Blocking queue...")
	// blocking queue example
	qs := queue.NewBlockingQueue()
	qs.SetSize(2)
	qs.Push(2)
	qs.Push(3)
	go qs.Push(6)
	fmt.Println(qs.Top())
	var chk1 int
	go func() {
		chk1 = qs.Pop().(int)
		fmt.Println("***", chk1)
	}()
	go qs.Push(4)
	go func() {
		chk1 = qs.Pop().(int)
		fmt.Println("***", chk1)
	}()

	time.Sleep(2 * time.Second)
	fmt.Println("xxx", qs.Top())
	fmt.Println("xxx", qs.Len())

	fmt.Scanf("%d", &chk)
}
