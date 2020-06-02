package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type backgroundNode struct {
	nodes.Node

	background api.INode

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

	var err error
	b.background, err = custom.NewStaticSquareNode("CenteredSquare", true, true, world, b)
	if err != nil {
		return err
	}
	gol := b.background.(*custom.StaticSquareNode)
	gol.SetColor(color.NewPaletteInt64(color.DarkerGray))

	return nil
}
