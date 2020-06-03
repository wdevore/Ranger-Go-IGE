package api

// INodeManager manages node on a stack and forms a SceneGraph
type INodeManager interface {
	Configure() error

	ClearEnabled(bool)

	SetPreNode(INode)
	PreVisit()
	Visit(interpolation float64) bool
	SetPostNode(INode)
	PostVisit()

	Update(msPerUpdate, secPerUpdate float64)

	PushNode(INode)
	PopNode() INode
	ReplaceNode(INode)

	RouteEvents(IEvent)

	RegisterTarget(target INode)
	UnRegisterTarget(target INode)

	RegisterEventTarget(target INode)
	UnRegisterEventTarget(target INode)

	End()

	Debug()
}
