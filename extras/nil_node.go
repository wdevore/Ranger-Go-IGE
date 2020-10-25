package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// NilNode is mostly for development and testing in non visual
// environments.
type NilNode struct {
	nodes.Node
}

// NewNilNode constructs shapeless node
func NewNilNode(name string) (api.INode, error) {
	o := new(NilNode)
	o.Initialize(name)

	return o, nil
}
