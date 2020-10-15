package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

type backgroundNode struct {
	nodes.Node

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

	atlas := world.GetAtlas(api.MonoAtlasName)

	b.Node.SetAtlas(atlas)

	mAtlas := atlas.(api.IStaticAtlasX)

	vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(true, true)
	b.shape = mAtlas.AddShape(api.CenteredFilledSquareShapeName, vertices, indices, mode)

	// For simplicity we set the color here.
	b.color = color.NewPaletteInt64(color.DarkGray).Array()

	return nil
}

// Draw renders shape
func (b *backgroundNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()
	atlas.SetColor(b.color)
	atlas.Render(b.shape, model)
}
