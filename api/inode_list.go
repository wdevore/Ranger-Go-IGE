package api

// INodeList is a simple list collection
type INodeList interface {
	Items() *[]INode
	DeleteAt(i int)
	FindFirstElement(node INode) int
	Add(node INode)
	Remove(node INode)
}
