package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	atlas api.IAtlasX

	sqr api.INode

	angle float64
}

func newGameLayer(name string, atlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.atlas = atlas

	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *gameLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	g.addLine("Type 'r' to return", -100.0, 0.0, 25, color.NewPaletteInt64(color.GoldYellow), world)

	var err error

	g.sqr, err = shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	g.sqr.SetScale(100.0)
	g.sqr.SetPosition(100.0, 100.0)
	gsq := g.sqr.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightOrange))

	return nil
}

func (g *gameLayer) addLine(text string, x, y, s float32, textColor api.IPalette, world api.IWorld) error {
	textureNode, err := shapes.NewBitmapFont9x9Node("SettingsLine", g.atlas, world, g)
	if err != nil {
		return err
	}
	textureNode.SetScale(s)
	textureNode.SetPosition(x, y)
	bf := textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(textColor.Array())
	bf.SetText(text)

	return nil
}

func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.sqr.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 1.5
}

func (g *gameLayer) EnterNode(man api.INodeManager) {
	// fmt.Println("gameLayer EnterNode")
	man.RegisterTarget(g)
}

// ExitScene called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	// fmt.Println("gameLayer ExitNode")
	man.UnRegisterTarget(g)
}
