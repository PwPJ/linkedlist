package linkedlist

import (
	"sync"
	"sync/atomic"
)

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	head   *Node
	length uint
}

var (
	Nodes []*Node
	part  uint = 10
)

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) Find(val int) (index uint, found bool) {
	current := l.head
	index = 0
	for current != nil {
		if current.Value == val {
			return index, true
		}
		current = current.Next
		index++
	}
	return 0, false
}

func (l *LinkedList) Remove(index uint) bool {
	if index >= l.length {
		return false
	}

	if index == 0 {
		l.head = l.head.Next
		l.length--
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current.Next == nil {
			return false
		}
		current = current.Next
	}

	if current.Next == nil {
		return false
	}

	current.Next = current.Next.Next
	l.length--

	partIndex := int(index / part)
	if index%part == 0 && partIndex < len(Nodes) {
		Nodes[partIndex] = current
	}

	return true
}

func (l *LinkedList) Get(index uint) (int, bool) {
	current := l.head
	for i := uint(0); i < index; i++ {
		if current == nil {
			return 0, false
		}
		current = current.Next
	}

	if current == nil {
		return 0, false
	}

	return current.Value, true
}

func (l *LinkedList) Insert(index uint, val int) bool {
	if index > l.length {
		return false
	}

	newNode := &Node{Value: val}

	if index == 0 {
		newNode.Next = l.head
		l.head = newNode
		l.length++
		if len(Nodes) == 0 {
			Nodes = append(Nodes, newNode)
		}
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		if current == nil {
			return false
		}
		current = current.Next
	}

	newNode.Next = current.Next
	current.Next = newNode
	l.length++

	partIndex := int(index / part)
	if partIndex >= len(Nodes) {
		// resize the Nodes array
		for len(Nodes) <= partIndex {
			Nodes = append(Nodes, nil)
		}
	}
	if index%part == 0 {
		Nodes[partIndex] = newNode
	}

	return true
}

func (l *LinkedList) HandleList() []int {
	current := l.head
	var values []int
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	return values
}

func (l *LinkedList) SearchConcurrently(wg *sync.WaitGroup, found *int32, find int) (int, bool) {
	var result int
	var isFound bool
	var mu sync.Mutex

	for i := 0; i < len(Nodes); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			index := Nodes[i]
			for j := 0; j < int(part); j++ {
				if atomic.LoadInt32(found) == 1 {
					return
				}

				if index == nil {
					return
				}

				if index.Value == find {
					atomic.StoreInt32(found, 1)
					mu.Lock()
					result = j + i
					isFound = true
					mu.Unlock()
					return
				}

				if index.Next == nil {
					return
				}
				index = index.Next
			}
		}(i)
	}

	wg.Wait()
	return result, isFound
}
