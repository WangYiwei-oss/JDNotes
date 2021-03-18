package stack

import "testing"

func Test(t *testing.T) {
	a := NewStack()
	a.Push("1")
	a.Push("2")
	a.Print()
	a.Pop()
	a.Print()
}
