package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	text  api.INode
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
	osql, err := custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	osql.SetScale(100.0)
	osql.SetPosition(110.0, 100.0)
	gol2 := osql.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightPurple))

	g.text, err = custom.NewDynamicTextNode("Text", 500, world, g)
	g.text.SetScale(2.0)
	gt := g.text.(*custom.DynamicTextNode)
	gt.SetText("Ranger Go!")

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
