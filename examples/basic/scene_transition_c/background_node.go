package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type backgroundNode struct {
	nodes.Node

	background api.INode

	color []float32
}

func newBackgroundNode(name string, world api.IWorld, parent api.INode, bgColor api.IPalette) api.INode {
	o := new(backgroundNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.build(world, bgColor)
	return o
}

func (b *backgroundNode) build(world api.IWorld, bgColor api.IPalette) error {
	b.Node.Build(world)

	var err error
	b.background, err = extras.NewStaticSquareNode("CenteredSquare", true, true, world, b)
	if err != nil {
		return err
	}
	gol := b.background.(*extras.StaticSquareNode)
	gol.SetColor(bgColor)

	return nil
}
