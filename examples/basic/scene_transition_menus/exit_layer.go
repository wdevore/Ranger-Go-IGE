package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
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

	osql, err := extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}

	osql.SetScale(100.0)
	osql.SetPosition(100.0, 100.0)
	gol2 := osql.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.Lime))

	return nil
}
