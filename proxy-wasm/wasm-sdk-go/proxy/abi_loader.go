package proxy

import (
	"github.com/mosn/wasm-sdk/proxy-wasm/wasm-sdk-go/proxy/types"
	"sync"
)

type loader struct {
	data map[types.ExtensionType]*Node
	mu   sync.RWMutex // protect data proxy
}

// register register extension loader
func (m *loader) register(eType types.ExtensionType, node *Node) {
	m.mu.Lock()

	if m.data == nil {
		m.data = make(map[types.ExtensionType]*Node, 4)
	}

	// we check head node first
	head := m.data[eType]
	if head == nil {
		head = &Node{
			count: 1,
		}
		// update head
		m.data[eType] = head
	}

	insert := &Node{
		Name:    node.Name,
		Context: node.Context,
		Order:   node.Order,
	}

	next := head.Next
	if next == nil {
		// fist insert, just insert to head
		head.Next = insert
		m.mu.Unlock()
		return
	}

	// we check already exist first
	for ; next != nil; next = next.Next {
		// we found it
		if next.Name == node.Name {
			// insert node has higher priority
			if node.Order > next.Order {
				next.Order = node.Order
				next.Context = node.Context
			}
			// release lock and do nothing
			m.mu.Unlock()
			return
		}
	}

	head.count++

	prev := head
	for next = head.Next; next != nil; next = next.Next {
		// we found insert position
		if node.Order > next.Order {
			break
		}
		prev = next
	}

	if prev.Next != nil {
		insert.Next = prev.Next.Next
	}

	prev.Next = insert

	m.mu.Unlock()
}

// find find extension instances.
func (m *loader) find(eType types.ExtensionType, name string) (node *Node, matched bool) {
	// we found nothing
	if m.data == nil {
		return nil, false
	}

	m.mu.RLocker().Lock()

	// we check head node first
	head := m.data[eType]
	if head == nil || head.count <= 0 {
		m.mu.RLocker().Unlock()
		return nil, false
	}

	node = head.Next

	// we expect find all eType Extension
	if len(name) == 0 {
		m.mu.RLocker().Unlock()
		return node, true
	}

	if head.count == 1 {
		m.mu.RLocker().Unlock()
		return node, name == node.Name
	}

	var count int
	var found *Node

	for ; node != nil; node = node.Next {
		if node.Name == name {
			if found == nil {
				found = node
			}
			count++
		}
	}

	m.mu.RLocker().Unlock()
	return found, count == 1
}

// clear all cache data
func (m *loader) clear() {
	m.mu.Lock()
	// help for gc
	m.data = nil
	m.mu.Unlock()
}

// Node pub or sub host metadata,
// Keep the structure as simple as possible.
type Node struct {
	Name    string      // extension name
	Context interface{} // extension context
	Order   int         // registered priority
	count   int         // number of node elements
	Next    *Node       // next node
}
