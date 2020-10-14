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
	zbar  api.INode
	plus  api.INode
	arc   api.INode
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

	dvr := world.Properties().Window.DeviceRes

	var err error

	// ---------------------------------------------------------
	shline, err := shapes.NewMonoHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*shapes.MonoHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	svline, err := shapes.NewMonoVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Width))
	gvl := svline.(*shapes.MonoVLineNode)
	gvl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	g.plus, err = shapes.NewMonoPlusNode("Plus", world, g)
	if err != nil {
		return err
	}
	g.plus.SetScale(float32(dvr.Width))
	g.plus.SetScale(500.0)
	gpl := g.plus.(*shapes.MonoPlusNode)
	gpl.SetColor(color.NewPaletteInt64(color.Lime))

	// ---------------------------------------------------------
	square, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(90.0, 80.0)
	gsq := square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	g.zbar, err = shapes.NewMonoZBarNode("ZBar", api.FILLOUTLINED, world, g)
	if err != nil {
		return err
	}
	g.zbar.SetScale(150.0)
	g.zbar.SetPosition(100.0, 150.0)
	gzr := g.zbar.(*shapes.MonoZBarNode)
	gzr.SetFilledColor(color.NewPaletteInt64(color.DarkBlue))

	// ---------------------------------------------------------
	circle, err := shapes.NewMonoCircleNode("Circle", api.FILLOUTLINED, 12, world, g)
	if err != nil {
		return err
	}
	circle.SetScale(100.0)
	circle.SetPosition(30.0, 80.0)
	gc := circle.(*shapes.MonoCircleNode)
	gc.SetFilledColor(color.NewPaletteInt64(color.GoldYellow))
	gc.SetOutlineColor(color.NewPaletteInt64(color.DarkOrange))
	gc.SetFilledAlpha(0.5)

	// ---------------------------------------------------------
	startAngle := -20 * maths.DegreeToRadians
	endAngle := 20 * maths.DegreeToRadians
	g.arc, err = shapes.NewMonoArcNode("Arc", api.FILLOUTLINED, 6, startAngle, endAngle, world, g)
	if err != nil {
		return err
	}
	g.arc.SetScale(400.0)
	g.arc.SetPosition(30.0, -80.0)
	gac := g.arc.(*shapes.MonoArcNode)
	gac.SetFilledColor(color.NewPaletteInt64(color.PanSkin))
	gac.SetOutlineColor(color.NewPaletteInt64(color.GreenYellow))
	gac.SetFilledAlpha(0.5)

	// ---------------------------------------------------------
	tri, err := shapes.NewMonoTriangleNode("Triangle", api.FILLED, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(55.0, 45.0)
	gt := tri.(*shapes.MonoTriangleNode)
	gt.SetFilledColor(color.NewPaletteInt64(color.Pink))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.zbar.SetRotation(maths.DegreeToRadians * g.angle * 2)
	g.plus.SetRotation(maths.DegreeToRadians * -g.angle)
	g.arc.SetRotation(maths.DegreeToRadians * -g.angle / 5)
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
