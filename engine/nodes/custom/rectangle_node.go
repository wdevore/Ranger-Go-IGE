package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// RectangleNode is a basic rectangle
type RectangleNode struct {
	nodes.Node

	color []float32

	shape api.IVectorShape
}

// NewRectangleNode constructs a rectangle shaped node
func NewRectangleNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(RectangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (r *RectangleNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = rendering.NewPaletteInt64(rendering.LightPurple).Array()

	r.shape = world.Atlas().Shape("Square")

	return nil
}

// SetColor sets line color
func (r *RectangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// Draw renders shape
func (r *RectangleNode) Draw(model api.IMatrix4) {
	// if r.IsDirty() {
	// 	r.SetDirty(false)
	// }

	w := r.World()
	gl.UniformMatrix4fv(w.ModelLoc(), 1, false, &model.Matrix()[0])

	gl.Uniform3fv(w.ColorLoc(), 1, &r.color[0])

	w.VecObj().Render(r.shape)
}
