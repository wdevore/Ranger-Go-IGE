package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	zbar  api.INode
	ozbar api.INode

	squareShape int
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

	var err error

	// ---------------------------------------------------------
	osql, err := newMySquareNode("MySquareNode", world, g)
	if err != nil {
		return err
	}
	osql.SetScale(100.0)
	osql.SetPosition(100.0, 100.0)
	gol2 := osql.(*mySquareNode)
	gol2.setColor(color.NewPaletteInt64(color.Pink))

	osql, err = newMyTriangleNode("MyTriangeNode", world, g)
	if err != nil {
		return err
	}
	osql.SetScale(100.0)
	osql.SetPosition(-100.0, 100.0)
	golt := osql.(*myTriangleNode)
	golt.setColor(color.NewPaletteInt64(color.GoldYellow))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.angle += 1.25
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
