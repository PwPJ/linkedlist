package linkedlist

import (
	"context"
	"sync"
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
	l.updateCacheForRemove(index)
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
		l.updateCacheForInsert(index, newNode)
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

	l.updateCacheForInsert(index, newNode)

	return true
}

func (l *LinkedList) updateCacheForRemove(index uint) {

	indexNodes := int(index/part) + 1

	if index == 0 {
		indexNodes = 0
	}

	for i := indexNodes; i < len(Nodes); i++ {
		if Nodes[i].Next != nil {
			Nodes[i] = Nodes[i].Next
		} else {
			Nodes = Nodes[:i]
		}

	}

}

func (l *LinkedList) updateCacheForInsert(index uint, newNode *Node) {
	partIndex := int(index / part)
	if partIndex >= len(Nodes) {
		Nodes = append(Nodes, newNode)
		return
	}

	var counter = index
	newNode1 := newNode
	for newNode1 != nil {
		if counter%part == 0 {
			if partIndex >= len(Nodes) {
				Nodes = append(Nodes, newNode1)
			} else {
				Nodes[partIndex] = newNode1
			}
			partIndex++
		}
		counter++
		newNode1 = newNode1.Next
	}
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

func (l *LinkedList) SearchConcurrently(ctx context.Context, cancel context.CancelFunc, find int) (int, bool) {
	var result int
	var isFound bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < len(Nodes); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			index := Nodes[i]
			for j := 0; j < int(part); j++ {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if index == nil {
					return
				}

				if index.Value == find {
					mu.Lock()
					result = (i * int(part)) + j
					isFound = true
					mu.Unlock()
					cancel()
					return
				}

				if index.Next == nil {
					return
				}

				select {
				case <-ctx.Done():
					return
				default:
					index = index.Next
				}
			}
		}(i)
	}

	wg.Wait()
	return result, isFound
}

func (l *LinkedList) SearchInSegmentedNodes(ctx context.Context, index int) (int, bool) {
	indexStart := index / 10
	if indexStart >= len(Nodes) {
		return 0, false
	}

	if index == 0 && len(Nodes) != 0 && Nodes[0].Value >= 0 {
		return Nodes[0].Value, true
	}

	nodes := Nodes[indexStart]
	targetIndex := index % 10

	for i := 0; i <= targetIndex; i++ {
		if i == targetIndex && nodes.Value >= 0 {
			return nodes.Value, true
		}
		if nodes.Next == nil {
			return 0, false
		}
		nodes = nodes.Next
	}

	return 0, false

}
