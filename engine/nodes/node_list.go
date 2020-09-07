package nodes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// NodeList is a simple list collection
// It can be used as a simple Pool.
type NodeList struct {
	items []api.INode
}

// NewNodeList create a new list collection
func NewNodeList() api.INodeList {
	o := new(NodeList)
	return o
}

// Items returns the internal collection
func (l *NodeList) Items() *[]api.INode {
	return &l.items
}

// Add adds item
func (l *NodeList) Add(node api.INode) {
	l.items = append(l.items, node)
}

// Remove removes item
func (l *NodeList) Remove(node api.INode) {
	idx := l.FindFirstElement(node)

	if idx >= 0 {
		l.DeleteAt(idx)
	} else {
		// fmt.Println("NodeManager: Unable to remove ", node, " node")
	}
}

// DeleteAt removes an item from the slice
func (l *NodeList) DeleteAt(i int) {
	// Remove the element at index i from slice.
	copy(l.items[i:], l.items[i+1:])   // Shift a[i+1:] left one index.
	l.items[len(l.items)-1] = nil      // Erase last element (write zero value).
	l.items = l.items[:len(l.items)-1] // Truncate slice.
}

// FindFirstElement finds the first item in the slice
func (l *NodeList) FindFirstElement(node api.INode) int {
	for idx, item := range l.items {
		if item.ID() == node.ID() {
			return idx
		}
	}

	return -1
}
