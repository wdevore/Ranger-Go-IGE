package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type gameLayer struct {
	nodes.Node

	angle float64
	zbar  api.INode
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

	dvr := world.Properties().Window.DeviceRes

	var err error

	// ---------------------------------------------------------
	shline, err := extras.NewStaticHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*extras.StaticHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	svline, err := extras.NewStaticVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Width))
	gvl := svline.(*extras.StaticVLineNode)
	gvl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	square, err := extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(90.0, 80.0)
	gsq := square.(*extras.StaticSquareNode)
	gsq.SetColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	circle, err := extras.NewStaticCircleNode("FilledCirle", true, world, g)
	if err != nil {
		return err
	}

	circle.SetScale(100.0)
	circle.SetPosition(30.0, 80.0)
	gc := circle.(*extras.StaticCircleNode)
	gc.SetColor(color.NewPaletteInt64(color.GoldYellow))
	gc.SetAlpha(0.5)

	// ---------------------------------------------------------
	tri, err := extras.NewStaticTriangleNode("FilledTri", true, true, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(55.0, 45.0)
	gt := tri.(*extras.StaticTriangleNode)
	gt.SetColor(color.NewPaletteInt64WithAlpha(color.Pink, 0.5))

	// ---------------------------------------------------------
	g.zbar, err = extras.NewStaticZBarNode("FilledZBar", true, world, g)
	if err != nil {
		return err
	}
	g.zbar.SetScale(100)
	g.zbar.SetPosition(100.0, 150.0)
	gzr := g.zbar.(*extras.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64WithAlpha(color.DarkBlue, 0.75))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.zbar.SetRotation(maths.DegreeToRadians * g.angle)
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
