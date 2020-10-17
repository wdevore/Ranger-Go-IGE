package api

// INodeManager manages node on a stack and forms a SceneGraph
type INodeManager interface {
	Configure(IWorld) error

	ClearEnabled(bool)

	SetRoot(INode)

	Begin() error
	End()
	Visit(interpolation float64) bool

	Update(msPerUpdate, secPerUpdate float64)

	PushNode(INode)
	PopNode() INode
	ReplaceNode(INode)

	RouteEvents(IEvent)

	RegisterTarget(target INode)
	UnRegisterTarget(target INode)

	RegisterEventTarget(target INode)
	UnRegisterEventTarget(target INode)

	Debug()
}
