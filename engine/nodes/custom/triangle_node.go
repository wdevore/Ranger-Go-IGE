package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// TriangleNode is a basic triangle
type TriangleNode struct {
	nodes.Node

	color []float32

	shape api.IVectorShape
}

// NewTriangleNode constructs a triangle shaped node
func NewTriangleNode(name string, world api.IWorld, parent api.INode) *TriangleNode {
	o := new(TriangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (t *TriangleNode) Build(world api.IWorld) error {
	t.Node.Build(world)

	t.color = rendering.NewPaletteInt64(rendering.White).Array()

	t.shape = world.Atlas().Shape("CenteredTriangle")

	return nil
}

// SetColor sets line color
func (t *TriangleNode) SetColor(color api.IPalette) {
	t.color = color.Array()
}

// Draw renders shape
func (t *TriangleNode) Draw(model api.IMatrix4) {
	// if t.IsDirty() {
	// 	t.SetDirty(false)
	// }

	w := t.World()
	gl.UniformMatrix4fv(w.ModelLoc(), 1, false, &model.Matrix()[0])

	gl.Uniform3fv(w.ColorLoc(), 1, &t.color[0])

	w.VecObj().Render(t.shape)
}
