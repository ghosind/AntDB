package core

type LinkedListNode struct {
	Value string
	Prev  *LinkedListNode
	Next  *LinkedListNode
}

type LinkedList struct {
	Head *LinkedListNode
	Tail *LinkedListNode
	Size int
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (ll *LinkedList) LPush(value string) {
	node := &LinkedListNode{Value: value}
	if ll.Size == 0 {
		ll.Head = node
		ll.Tail = node
	} else {
		node.Next = ll.Head
		ll.Head.Prev = node
		ll.Head = node
	}
	ll.Size++
}

func (ll *LinkedList) RPush(value string) {
	node := &LinkedListNode{Value: value}
	if ll.Size == 0 {
		ll.Head = node
		ll.Tail = node
	} else {
		node.Prev = ll.Tail
		ll.Tail.Next = node
		ll.Tail = node
	}
	ll.Size++
}

func (ll *LinkedList) LPop() (string, bool) {
	if ll.Size == 0 {
		return "", false
	}
	node := ll.Head
	ll.RemoveNode(node)
	return node.Value, true
}

func (ll *LinkedList) RPop() (string, bool) {
	if ll.Size == 0 {
		return "", false
	}

	node := ll.Tail
	ll.RemoveNode(node)
	return node.Value, true
}

func (ll *LinkedList) IndexAt(index int) *LinkedListNode {
	if index < 0 {
		index = ll.Size + index
	}
	if index >= ll.Size || index < 0 {
		return nil
	}

	current := ll.Head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	return current
}

func (ll *LinkedList) Set(index int, value string) error {
	if index < 0 {
		index = ll.Size + index
	}
	if index >= ll.Size || index < 0 {
		return ErrOutOfRange
	}

	current := ll.Head
	for i := 0; i < index; i++ {
		current = current.Next
	}
	current.Value = value
	return nil
}

func (ll *LinkedList) RemoveNode(node *LinkedListNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		ll.Head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		ll.Tail = node.Prev
	}
	ll.Size--
}
