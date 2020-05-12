package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticAtlasNode is a basic circle
type StaticAtlasNode struct {
	nodes.Node

	color []float32

	atlasName string
	shape     api.IVectorShape
}

// NewStaticAtlasNode constructs a rectangle shaped node
func NewStaticAtlasNode(name, atlasName string, world api.IWorld, parent api.INode) api.INode {
	o := new(StaticAtlasNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.atlasName = atlasName
	o.Build(world)
	return o
}

// Build configures the node
func (r *StaticAtlasNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.LightPurple).Array()

	r.shape = world.Atlas().Shape(r.atlasName)

	return nil
}

// SetColor sets line color
func (r *StaticAtlasNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// Draw renders shape
func (r *StaticAtlasNode) Draw(model api.IMatrix4) {
	w := r.World()
	gl.UniformMatrix4fv(w.ModelLoc(), 1, false, &model.Matrix()[0])

	gl.Uniform3fv(w.ColorLoc(), 1, &r.color[0])

	w.VecObj().Render(r.shape)
}
