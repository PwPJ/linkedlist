package main

import (
	"fmt"
)

// Node represents an element in the linked list
type Node struct {
	Value int
	Next  *Node
}

// L represents a linked list
type L struct {
	head   *Node
	length uint
}

// new creates and returns a new instance of a linked list
func new() *L {
	return &L{}
}

func (l *L) Find(n int) (index uint, found bool) {
	current := l.head
	index = 0
	for current != nil {
		if current.Value == n {
			return index, true
		}
		current = current.Next
		index++
	}
	return 0, false
}

func (l *L) Get(index uint) (int, bool) {
	current := l.head

	for i := uint(0); i < index; i++ {
		if current == nil {
			return 0, false // return 0 and false if index is out of range
		}
		current = current.Next // move to the next node
	}

	if current == nil {
		return 0, false // return 0 and false if index is exactly one past the last node
	}

	return current.Value, true // return the value at the node and true indicating success
}

func (l *L) Insert(index uint, val int) bool {
	if index > l.length {
		return false
	}

	if index == 0 {
		l.head = &Node{Value: val, Next: l.head}
		l.length++
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current == nil {
			return false
		}
		current = current.Next
	}

	current.Next = &Node{Value: val, Next: current.Next}
	l.length++
	return true //successful insertion
}

func (l *L) Remove(index uint) bool {
	if index >= l.length { // Check if index is out of range
		return false
	}

	if index == 0 { // Remove the first element
		l.head = l.head.Next
		l.length--
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current.Next == nil {
			return false // If there's no next element, fail
		}
		current = current.Next
	}

	if current.Next == nil {
		return false // Safety check for the last element
	}

	current.Next = current.Next.Next // Adjust the pointer to skip over the removed node
	l.length--
	return true
}

func main() {
	newList := &L{}
	newList.Insert(0, 100)
	newList.Insert(1, 200)
	newList.Insert(2, 300)

	value, ok := newList.Get(1) // value at index 1
	if !ok {
		fmt.Println("failed, index out of range")
	} else {
		fmt.Println("value at index 1:", value)
	}

	value, ok = newList.Get(8) // value at an out-of-range index
	if !ok {
		fmt.Println("Retrieval failed, index out of range")
	} else {
		fmt.Println("Value at index 8:", value)
	}
}
