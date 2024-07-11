package linkedlist

type Node struct {
	Value int
	Next  *Node
}

type LinkedList struct {
	head   *Node
	length uint
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) Find(n int) (index uint, found bool) {
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
	return true
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
