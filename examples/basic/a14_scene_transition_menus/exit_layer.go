package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameExitLayer struct {
	nodes.Node
}

func newBasicExitLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameExitLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *gameExitLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	square, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(100.0, 100.0)
	gsq := square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Lime))

	return nil
}
