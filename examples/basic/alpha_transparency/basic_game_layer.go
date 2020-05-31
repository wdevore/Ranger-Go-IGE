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
	shline, err := custom.NewStaticHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*custom.StaticHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	svline, err := custom.NewStaticVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Width))
	gvl := svline.(*custom.StaticVLineNode)
	gvl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	square, err := custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(90.0, 80.0)
	gsq := square.(*custom.StaticSquareNode)
	gsq.SetColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	circle, err := custom.NewStaticCircleNode("FilledCirle", true, world, g)
	if err != nil {
		return err
	}

	circle.SetScale(100.0)
	circle.SetPosition(30.0, 80.0)
	gc := circle.(*custom.StaticCircleNode)
	gc.SetColor(color.NewPaletteInt64(color.GoldYellow))
	gc.SetAlpha(0.5)

	// ---------------------------------------------------------
	tri, err := custom.NewStaticTriangleNode("FilledTri", true, true, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(55.0, 45.0)
	gt := tri.(*custom.StaticTriangleNode)
	gt.SetColor(color.NewPaletteInt64WithAlpha(color.Pink, 0.5))

	// ---------------------------------------------------------
	g.zbar, err = custom.NewStaticZBarNode("FilledZBar", true, world, g)
	if err != nil {
		return err
	}
	g.zbar.SetScale(100)
	g.zbar.SetPosition(100.0, 150.0)
	gzr := g.zbar.(*custom.StaticZBarNode)
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
