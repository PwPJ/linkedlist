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
	tail   *Node
	middle *Node
	length uint
}

var (
	Nodes []*Node
	part  uint = 10
)

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) updateMiddleOnInsert(index uint) {
	if l.length == 1 {
		l.middle = l.head
		return
	}
	if l.middle == nil { // Add a nil check to ensure safe access
		l.middle = l.head
	}
	if index <= l.length/2 {
		if l.length%2 == 0 && l.middle.Next != nil { // Ensure l.middle.Next is not nil
			l.middle = l.middle.Next
		}
	}
}

func (l *LinkedList) updateMiddleOnRemove(index uint) {
	if l.length == 0 {
		l.middle = nil
		return
	}
	if index <= l.length/2 {
		if l.length%2 == 1 {
			current := l.head
			for i := uint(0); i < l.length/2-1; i++ {
				current = current.Next
			}
			l.middle = current
		}
	}
}

func (l *LinkedList) Insert(index uint, val int) bool {
	if index > l.length {
		return false
	}

	newNode := &Node{Value: val}

	if index == 0 {
		newNode.Next = l.head
		l.head = newNode
		if l.length == 0 {
			l.tail = newNode
			l.middle = newNode
		}
		l.length++
		if l.length >= 9 {
			l.updateCacheForInsert(index, newNode) // Maintain the cache
		}
		l.updateMiddleOnInsert(index) // Update the middle pointer
		return true
	}

	if index == l.length {
		if l.tail != nil {
			l.tail.Next = newNode
		}
		l.tail = newNode
		l.length++
		if l.length >= 9 {
			l.updateCacheForInsert(index, newNode) // Maintain the cache
		}
		l.updateMiddleOnInsert(index) // Update the middle pointer
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
	if l.length >= 9 {
		l.updateCacheForInsert(index, newNode) // Maintain the cache
	}
	l.updateMiddleOnInsert(index) // Update the middle pointer
	return true
}

func (l *LinkedList) Remove(index uint) bool {
	if index >= l.length {
		return false
	}

	if index == 0 {
		l.head = l.head.Next
		l.length--
		if l.length >= 9 {
			l.updateCacheForRemove(index)
			l.updateMiddleOnRemove(index)
		}
		return true
	}

	current := l.head
	for i := uint(0); i < index-1; i++ {
		current = current.Next
	}
	if current.Next == nil {
		return false
	}
	if current.Next == l.tail {
		l.tail = current
	}
	current.Next = current.Next.Next
	l.length--
	if l.length >= 9 {
		l.updateCacheForRemove(index)
		l.updateMiddleOnRemove(index)
	}
	return true
}

func (l *LinkedList) Find(val int) (uint, bool) {
	if l.length < 9 {
		current := l.head
		index := uint(0)
		for current != nil {
			if current.Value == val {
				return index, true
			}
			current = current.Next
			index++
		}
		return 0, false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	index, found := l.SearchConcurrently(ctx, cancel, val)
	return uint(index), found
}

func (l *LinkedList) Get(index uint) (int, bool) {
	if index >= l.length {
		return 0, false
	}

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

func (l *LinkedList) HandleList() []int {
	current := l.head
	var values []int
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	return values
}

func (l *LinkedList) updateCacheForInsert(index uint, newNode *Node) {
	partIndex := int(index / part)
	if partIndex >= len(Nodes) {
		Nodes = append(Nodes, newNode)
		return
	}

	counter := index
	current := newNode
	for current != nil {
		if counter%part == 0 {
			if partIndex >= len(Nodes) {
				Nodes = append(Nodes, current)
			} else {
				Nodes[partIndex] = current
			}
			partIndex++
		}
		counter++
		current = current.Next
	}
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

func (l *LinkedList) SearchConcurrently(ctx context.Context, cancel context.CancelFunc, target int) (int, bool) {
	var result int
	var isFound bool
	var mu sync.Mutex
	var wg sync.WaitGroup
	done := make(chan struct{})

	searchSegment := func(start *Node, limit int, step int, offset int) {
		defer wg.Done()
		current := start
		for i := 0; i < limit; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				if current == nil {
					return
				}
				if current.Value == target {
					mu.Lock()
					if !isFound {
						result = offset + i*step
						isFound = true
						cancel()
						close(done)
					}
					mu.Unlock()
					return
				}
				current = current.Next
			}
		}
	}

	quarter := int(l.length) / 4
	wg.Add(4)

	go searchSegment(l.head, quarter, 1, 0)
	go searchSegment(l.tail, quarter, -1, int(l.length)-1)
	go searchSegment(l.middle, quarter, 1, int(l.length/2))
	go searchSegment(l.middle, quarter, -1, int(l.length/2))

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
