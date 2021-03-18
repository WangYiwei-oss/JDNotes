package stack

import "fmt"

type Node struct {
	data interface{}
	next *Node
}

type Stack struct {
	topeNode *Node
}

func NewStack() *Stack {
	return &Stack{
		nil,
	}
}

func (this *Stack) Push(value interface{}) {
	this.topeNode = &Node{data: value, next: this.topeNode}
}

func (this *Stack) IsEmpty() bool {
	if this.topeNode == nil {
		return true
	}

	return false
}

func (this *Stack) Pop() interface{} {
	if this.IsEmpty() {
		return nil
	}

	value := this.topeNode.data
	this.topeNode = this.topeNode.next
	return value
}

func (this *Stack) Top() interface{} {
	if this.IsEmpty() {
		return nil
	}

	return this.topeNode.data
}

func (this *Stack) Print() {
	if this.IsEmpty() {
		fmt.Println("empty Stack")
	}

	cur := this.topeNode
	for nil != cur {
		fmt.Println(cur.data)
		cur = cur.next
	}

}
