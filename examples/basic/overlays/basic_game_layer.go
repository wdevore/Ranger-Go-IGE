package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type gameLayer struct {
	nodes.Node
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

func (g *gameLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	// ---------------------------------------------------------
	osql, err := extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	osql.SetScale(100.0)
	osql.SetPosition(100.0, 100.0)
	gol2 := osql.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	return nil
}
