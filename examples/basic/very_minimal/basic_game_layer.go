package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *gameLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	osql, err := custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}

	osql.SetScale(100.0)
	osql.SetPosition(100.0, 100.0)
	gol2 := osql.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightPurple))

	return nil
}
