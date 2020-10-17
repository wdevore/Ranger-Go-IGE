package nodes

import "github.com/wdevore/Ranger-Go-IGE/api"

// Group holds the children properties and methods.
type Group struct {
	children []api.INode
}

func (g *Group) initializeGroup() {
	g.children = []api.INode{}
}

// Children returns the children of current node.
// Nodes should override this method for providing any child they contain.
func (g *Group) Children() []api.INode {
	return g.children
}

// AddChild adds a node to this node
func (g *Group) AddChild(child api.INode) {
	if child != nil {
		g.children = append(g.children, child)
	}
}

// PrependChild adds the give node to the start of the collection rather than the end.
func (g *Group) PrependChild(child api.INode) {
	if child != nil {
		front := []api.INode{child}
		g.children = append(front, g.children[1:]...)
	}
}

// GetChildByID finds an INode by ID.
func (g *Group) GetChildByID(id int) api.INode {
	if len(g.children) > 0 {
		for _, child := range g.children {
			if child.ID() == id {
				return child
			}
		}
	}

	return nil
}

// GetChildByName finds an INode by Name
func (g *Group) GetChildByName(name string) api.INode {
	if len(g.children) > 0 {
		for _, child := range g.children {
			if child.Name() == name {
				return child
			}
		}
	}

	return nil
}

// InsertAndShift adds "newNode" at the begining and shifts everything else to the End
// while returning the End.
// The `width` field determines where node drop-off the End.
func (g *Group) InsertAndShift(newNode api.INode, width int) api.INode {
	l := len(g.children)

	var r api.INode

	// Now move first element towards the End to make room for the new node
	switch l {
	case 0:
		// Just append
		g.children = append(g.children, newNode)
		r = nil // Indicates there was no End node already existing.
	case 1:
		// Grab End element for returning
		r = g.children[l-1]
		g.children = []api.INode{newNode, r}
	default:
		// Grab End element for returning
		r = g.children[l-1]

		// Shift all nodes
		shift := []api.INode{newNode}
		g.children = append(shift, g.children[1:width-1]...)
	}

	return r
}

// RemoveLast removes the last Node and returns it.
func (g *Group) RemoveLast() api.INode {
	l := len(g.children)
	var n api.INode

	if l > 0 {
		n = g.children[l-1]
		g.children = g.children[:l-1]
	}

	return n
}
