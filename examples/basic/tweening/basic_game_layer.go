package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	// Tanema's framework
	tween *gween.Tween

	square api.INode
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	var err error

	// ---------------------------------------------------------
	g.square, err = custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	g.square.SetScale(100.0)
	g.square.SetPosition(100.0, 100.0)
	gol2 := g.square.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	// 5s = 5000ms
	g.tween = gween.New(float32(g.square.Position().X()), -600.0, 5000, ease.OutExpo)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	value, isFinished := g.tween.Update(float32(msPerUpdate))

	if !isFinished {
		g.square.SetPosition(value, g.square.Position().Y())
	} else {
		g.tween.Reset()
	}
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
