package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// groupNode is a non-rendering node meant to contain children only.
type groupNode struct {
	nodes.Node
}

// NewGroupNode constructs group node container.
func NewGroupNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(groupNode)

	o.Initialize(name)
	o.SetParent(parent)

	if parent != nil {
		parent.AddChild(o)
	}

	return o, nil
}
