package custom

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicAtlasNode is a generic node
type DynamicAtlasNode struct {
	nodes.Node

	color []float32

	atlasName string
	shape     api.IAtlasShape
}

// NewDynamicAtlasNode constructs a generic shape node
func NewDynamicAtlasNode(name, atlasName string, world api.IWorld, parent api.INode) api.INode {
	o := new(DynamicAtlasNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.atlasName = atlasName
	o.Build(world)
	return o
}

// Build configures the node
func (r *DynamicAtlasNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.LightPurple).Array()

	r.shape = world.DynoAtlas().Shape(r.atlasName)

	return nil
}

// SetColor sets line color
func (r *DynamicAtlasNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// Draw renders shape
func (r *DynamicAtlasNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
