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
	ozbar api.INode
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
	osql, err := extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	osql.SetScale(100.0)
	osql.SetPosition(110.0, 100.0)
	gol2 := osql.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightPurple))

	// ---------------------------------------------------------
	tri, err := extras.NewStaticTriangleNode("FilledTri", true, true, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100)
	tri.SetPosition(-100.0, 100.0)
	gtr := tri.(*extras.StaticTriangleNode)
	gtr.SetColor(color.NewPaletteInt64(color.Pink))

	// ---------------------------------------------------------
	circle, err := extras.NewStaticCircleNode("FilledCirle", true, world, g)
	if err != nil {
		return err
	}
	circle.SetScale(100)
	circle.SetPosition(0.0, 100.0)
	gcr := circle.(*extras.StaticCircleNode)
	gcr.SetColor(color.NewPaletteInt64(color.GoldYellow))

	ocircle, err := extras.NewStaticCircleNode("OutlineCirle", false, world, g)
	if err != nil {
		return err
	}
	ocircle.SetScale(100)
	ocircle.SetPosition(0.0, 100.0)
	gocr := ocircle.(*extras.StaticCircleNode)
	gocr.SetColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	circle2, err := extras.NewStaticCircleNode("FilledCirle", true, world, g)
	if err != nil {
		return err
	}
	circle2.SetScale(100)
	circle2.SetPosition(0.0, -150.0)
	gcr2 := circle2.(*extras.StaticCircleNode)
	gcr2.SetColor(color.NewPaletteInt64(color.SoftGreen))

	// ---------------------------------------------------------
	g.zbar, err = extras.NewStaticZBarNode("FilledZBar", true, world, g)
	if err != nil {
		return err
	}
	g.zbar.SetScale(100)
	g.zbar.SetPosition(300.0, 100.0)
	gzr := g.zbar.(*extras.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.LightNavyBlue))

	g.ozbar, err = extras.NewStaticZBarNode("OutlineZBar", false, world, g)
	if err != nil {
		return err
	}
	g.ozbar.SetScale(100)
	g.ozbar.SetPosition(300.0, 100.0)
	gzr = g.ozbar.(*extras.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	point, err := extras.NewStaticPointNode("Point", world, g)
	if err != nil {
		return err
	}
	point.SetScale(100)
	point.SetPosition(100.0, -100.0)
	gp := point.(*extras.StaticPointNode)
	gp.SetColor(color.NewPaletteInt64(color.Aqua))
	gp.SetSize(10.0)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.zbar.SetRotation(maths.DegreeToRadians * g.angle)
	g.ozbar.SetRotation(maths.DegreeToRadians * g.angle)
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
