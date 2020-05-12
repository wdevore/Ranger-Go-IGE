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
	crow  api.INode
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

	square := custom.NewStaticAtlasNode("Square", "CenteredSquare", world, g)
	square.SetScale(100.0)
	square.SetPosition(110.0, 100.0)

	circle := custom.NewStaticAtlasNode("Circle", "Circle12Segments", world, g)
	circle.SetScale(100.0)
	circle.SetPosition(0.0, 100.0)
	gc := circle.(*custom.StaticAtlasNode)
	gc.SetColor(color.NewPaletteInt64(color.GoldYellow))

	tri := custom.NewStaticAtlasNode("Triangle", "CenteredTriangle", world, g)
	tri.SetScale(100.0)
	tri.SetPosition(-100.0, 100.0)
	gt := tri.(*custom.StaticAtlasNode)
	gt.SetColor(color.NewPaletteInt64(color.Pink))

	line := custom.NewStaticAtlasNode("Line", "Line", world, g)
	line.SetScale(100.0)
	line.SetPosition(-200.0, 100.0)
	line.SetRotation(maths.DegreeToRadians * 45.0)
	ge := line.(*custom.StaticAtlasNode)
	ge.SetColor(color.NewPaletteInt64(color.GreenYellow))

	cross := custom.NewStaticAtlasNode("Plus", "Cross", world, g)
	cross.SetScale(100.0)
	cross.SetPosition(-300.0, 100.0)
	// cross.SetRotation(maths.DegreeToRadians * 45.0)
	gs := cross.(*custom.StaticAtlasNode)
	gs.SetColor(color.NewPaletteInt64(color.PanSkin))

	g.crow = custom.NewStaticAtlasNode("Crow", "CrowBar", world, g)
	g.crow.SetScale(100.0)
	g.crow.SetPosition(300.0, 100.0)
	gb := g.crow.(*custom.StaticAtlasNode)
	gb.SetColor(color.NewPaletteInt64(color.LightOrange))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.crow.SetRotation(maths.DegreeToRadians * g.angle)
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
