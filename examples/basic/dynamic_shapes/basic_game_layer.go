package main

import (
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	sqr   api.INode
	line  api.INode
	line2 api.INode

	dynoTxt api.INode
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

	// -------------------------------------------------------------
	shline, err := custom.NewStaticHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*custom.StaticHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// -------------------------------------------------------------
	svline, err := custom.NewStaticVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Width))
	gvl2 := svline.(*custom.StaticVLineNode)
	gvl2.SetColor(color.NewPaletteInt64(color.LightGray))

	// -------------------------------------------------------------
	g.sqr, err = custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	g.sqr.SetScale(100)
	g.sqr.SetPosition(100.0, 100.0)
	gol2 := g.sqr.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	// -------------------------------------------------------------
	g.dynoTxt, err = custom.NewDynamicTextNode("DynoTxt", 500, world, g)
	if err != nil {
		return err
	}
	g.dynoTxt.SetScale(2.0)
	g.dynoTxt.SetPosition(-100.0, 100.0)
	gd := g.dynoTxt.(*custom.DynamicTextNode)
	gd.SetText("Ranger is a Go!")
	gd.SetColor(color.NewPaletteInt64(color.LightPink).Array())
	gd.SetPixelSize(3.0)

	// ---------------------------------------------------------------------
	g.line, err = newDynamicLineNode("DynoLin", world, g)
	if err != nil {
		return err
	}
	glc := g.line.(*DynamicLineNode)
	glc.SetColor(color.NewPaletteInt64(color.White))
	glc.SetPoint1(100.0, -100.0)
	glc.SetPoint2(50.0, -50.0)

	// ---------------------------------------------------------------------
	g.line2, err = newDynamicLineNode("DynoLin2", world, g)
	if err != nil {
		return err
	}
	glc2 := g.line2.(*DynamicLineNode)
	glc2.SetColor(color.NewPaletteInt64(color.GreenYellow))
	glc2.SetPoint1(-100.0, -100.0)
	glc2.SetPoint2(-200.0, -200.0)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.sqr.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 1.5

	g.dynoTxt.SetRotation(maths.DegreeToRadians * -g.angle / 10)

	glc := g.line.(*DynamicLineNode)
	x := math.Cos(maths.DegreeToRadians * g.angle)
	y := math.Sin(maths.DegreeToRadians * g.angle)
	glc.SetPoint2(100.0*float32(x), 100.0*float32(y))

	glc2 := g.line2.(*DynamicLineNode)
	x2 := math.Cos(maths.DegreeToRadians * -g.angle / 2)
	y2 := math.Sin(maths.DegreeToRadians * -g.angle / 2)
	glc2.SetPoint1(-150.0*float32(x2), -150.0*float32(y2))
	glc2.SetPoint2(-50.0*float32(x2), -50.0*float32(y2))
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
