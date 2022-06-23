package stack

import (
	"github.com/muzin/go_rt/collection/vector"
)

type Stack struct {
	vector *vector.Vector
}

func NewStack() *Stack {
	vector := vector.NewVector()
	return &Stack{vector}
}

func (stack *Stack) Push(value interface{}) {
	stack.vector.Add(value)
}

func (stack *Stack) Pop() interface{} {
	lastIndex := stack.vector.Size() - 1
	if (lastIndex) >= 0 {
		lastElement := stack.vector.Remove(lastIndex)
		return lastElement
	}
	return nil
}

func (stack *Stack) Shift() interface{} {
	firstIndex := 0
	if (firstIndex) >= 0 {
		firstElement := stack.vector.Remove(firstIndex)
		return firstElement
	}
	return nil
}

func (stack *Stack) Get(i int) interface{} {
	obj := stack.vector.Get(i)
	return obj
}

func (stack *Stack) Size() int {
	return stack.vector.Size()
}

func (stack *Stack) IsEmpty() bool {
	return stack.vector.Size() == 0
}
