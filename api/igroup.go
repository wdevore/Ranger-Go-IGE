package api

// IGroup is a collection of nodes. Group nodes can't be leafs.
type IGroup interface {
	// Children returns the children of current node.
	// Nodes should override this method for providing any child they contain.
	Children() []INode

	AddChild(INode)
	PrependChild(INode)

	GetChildByID(id int) INode
	GetChildByName(name string) INode

	InsertAndShift(newNode INode, width int) INode
	RemoveLast() INode
}
