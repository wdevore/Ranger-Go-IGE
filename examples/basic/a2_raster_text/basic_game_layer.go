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

	angle float64
	text  api.INode
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

	// ---------------------------------------------------------
	square, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(110.0, 100.0)
	gsq := square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightPurple))

	// ---------------------------------------------------------
	g.text, err = shapes.NewDynamicPixelTextNode("Text", world, g)
	if err != nil {
		return err
	}
	g.text.SetScale(2.0)
	gt := g.text.(*shapes.DynamicPixelPixelTextNode)
	gt.SetText("Ranger Go!")
	gt.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.text.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 0.25
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
}
