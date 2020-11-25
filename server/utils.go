package main

import (
	"sync"
)

// SLLNode holds the data of a SLList
type SLLNode struct {
	data interface{}
	next *SLLNode
}

// SLList is a thread-safe linked list
type SLList struct {
	head, tail *SLLNode
	sync.RWMutex
}

// Push adds a node, data, or another linked list to the beginning of the list
func (ll *SLList) Push(data interface{}) *SLLNode {
	ll.Lock()
	defer ll.Unlock()
	if data == nil {
		return ll.head
	}
	switch data.(type) {
	case *SLLNode:
		node := data.(*SLLNode)
		node.next = ll.head
		ll.head = node
		if ll.tail == nil {
			ll.tail = node
		}
	case *SLList:
		list := data.(*SLList)
		list.tail = ll.head
		ll.head = list.head
		if ll.tail == nil {
			ll.tail = list.tail
		}
	default:
		node := new(SLLNode)
		node.data = data
		node.next = ll.head
		ll.head = node
		if ll.tail == nil {
			ll.tail = node
		}
	}
	return ll.head
}

// Append adds a node, data, or another linked list to the end of the list
func (ll *SLList) Append(data interface{}) *SLLNode {
	ll.Lock()
	defer ll.Unlock()
	if data == nil {
		return ll.head
	}
	switch data.(type) {
	case *SLLNode:
		ll.tail = data.(*SLLNode)
		if ll.head == nil {
			ll.head = ll.tail
		}
	case *SLList:
		ll.tail = data.(*SLList).head
		if ll.head == nil {
			ll.head = ll.tail
		}
	default:
		ll.tail = new(SLLNode)
		ll.tail.data = data
		if ll.head == nil {
			ll.head = ll.tail
		}
	}
	return ll.head
}

// Remove removes a node from the linked list; takes either a node or data
func (ll *SLList) Remove(data interface{}) *SLLNode {
	ll.Lock()
	defer ll.Unlock()
	var curr, prev, node *SLLNode
	curr = ll.head
	prev = ll.head
	switch data.(type) {
	case *SLLNode:
		n := data.(*SLLNode)
		for curr != nil {
			if curr == n {
				node := curr
				if curr == ll.head {
					ll.head = nil
					ll.tail = nil
				} else if curr == ll.tail {
					ll.tail = prev
				}
				prev.next = curr.next
				node.next = nil
				break
			}
			prev = curr
			curr = curr.next
		}
	default:
		for curr != nil {
			if curr.data == data {
				if curr == ll.head {
					node = curr
					ll.head = nil
				} else {
					node = curr
					prev.next = curr.next
				}
				break
			}
			prev = curr
			curr = curr.next
		}
	}
	return node
}

// Pop removes the first item of the list
func (ll *SLList) PopFirst() (node *SLLNode) {
	ll.Lock()
	defer ll.Unlock()
	node = ll.head
	if node == nil {
		return
	}
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		ll.head = node.next
	}
	return
}

// Pop removes the first item of the list
func (ll *SLList) PopLast() (node *SLLNode) {
	ll.Lock()
	defer ll.Unlock()
	node = ll.tail
	if node == nil {
		return
	}
	if ll.head == ll.tail {
		ll.head = nil
		ll.tail = nil
	} else {
		n := ll.head
		for n.next != nil {
			n = n.next
		}
		ll.tail = n
	}
	return
}

// Length returns the length of the linked list
func (ll *SLList) Length() int {
	ll.RLock()
	defer ll.RUnlock()
	curr := ll.head
	len := 0
	for curr != nil {
		len++
		curr = curr.next
	}
	return len
}

// ForEach iterates over each list element and applies a function to each
func (ll *SLList) ForEach(f func(*SLLNode, ...interface{}), args ...interface{}) {
	ll.RLock()
	defer ll.RUnlock()
	node := ll.head
	for node != nil {
		f(node, args)
		node = node.next
	}
}
