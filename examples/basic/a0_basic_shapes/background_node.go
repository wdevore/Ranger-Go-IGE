package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type backgroundNode struct {
	nodes.Node

	atlas api.IAtlasX
	shape int

	color []float32
}

func newBackgroundNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(backgroundNode)

	o.Initialize(name)
	o.SetParent(parent)

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *backgroundNode) build(world api.IWorld) error {
	b.Node.Build(world)

	b.atlas = world.GetAtlas(api.MonoAtlasName)

	b.shape = b.atlas.GetShapeByName(api.CenteredFilledSquareShapeName)

	// For simplicity we set the color here.
	b.color = color.NewPaletteInt64(color.DarkGray).Array()

	return nil
}

func (b *backgroundNode) setColor(color api.IPalette) {
	b.color = color.Array()
}

func (b *backgroundNode) setAlpha(alpha float32) {
	b.color[3] = alpha
}

func (b *backgroundNode) Atlas() api.IAtlasX {
	return b.atlas
}

// Draw renders shape
func (b *backgroundNode) Draw(model api.IMatrix4) {
	b.atlas.SetColor(b.color)
	b.atlas.Render(b.shape, model)
}