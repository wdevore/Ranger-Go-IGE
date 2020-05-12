package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type backgroundNode struct {
	nodes.Node

	background api.IVectorShape

	color []float32
}

func newBackgroundNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(backgroundNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

func (b *backgroundNode) Build(world api.IWorld) error {
	b.Node.Build(world)

	b.SetDirty(true)

	b.color = color.NewPaletteInt64(color.DarkerGray).Array()

	b.background = world.Atlas().Shape("CenteredSquare")

	return nil
}

// -----------------------------------------------------
// Visuals
// -----------------------------------------------------

func (b *backgroundNode) Draw(model api.IMatrix4) {
	// if b.IsDirty() {
	// 	b.SetDirty(false)
	// }

	w := b.World()
	gl.UniformMatrix4fv(w.ModelLoc(), 1, false, &model.Matrix()[0])

	gl.Uniform3fv(w.ColorLoc(), 1, &b.color[0])

	w.VecObj().Render(b.background)
}
