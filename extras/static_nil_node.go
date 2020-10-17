package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// StaticNilNode is mostly for development and testing in non visual
// environments.
type StaticNilNode struct {
	nodes.Node
}

// NewStaticNilNode constructs shapeless node
func NewStaticNilNode(name string) (api.INode, error) {
	o := new(StaticNilNode)
	o.Initialize(name)

	return o, nil
}
