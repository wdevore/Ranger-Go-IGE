package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicLineNode is a generic node
type DynamicLineNode struct {
	nodes.Node

	color []float32

	atlasName string
	shape     api.IAtlasShape
}

// NewDynamicLineNode constructs a generic shape node
func newDynamicLineNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(DynamicLineNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (r *DynamicLineNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.LightPurple).Array()

	r.shape = world.DynoAtlas().Shape("Line")

	return nil
}

// SetColor sets line color
func (r *DynamicLineNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetPoint1 sets point 1 end position
func (r *DynamicLineNode) SetPoint1(x, y float32) {
	r.shape.SetVertex2D(x, y, 0)
	r.SetDirty(true)
}

// SetPoint2 sets point 2 end position
func (r *DynamicLineNode) SetPoint2(x, y float32) {
	r.shape.SetVertex2D(x, y, 1)
	r.SetDirty(true)
}

// Draw renders shape
func (r *DynamicLineNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicRenderGraphic)

	if r.IsDirty() {
		// Update buffer
		renG.Update(0, r.shape.Count())
		r.SetDirty(false)
	}

	renG.SetColor(r.color)

	renG.Render(r.shape, model)
}
