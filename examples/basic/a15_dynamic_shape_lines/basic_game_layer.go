package main

import (
	"fmt"
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	angle float64
	sqr   api.INode

	line  api.INode
	line2 api.INode
	line3 api.INode

	dynoTxt api.INode

	debug    float64
	debugCnt float64
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)

	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	o.debug = 100
	if err := o.build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	dvr := world.Properties().Window.DeviceRes

	// -------------------------------------------------------------
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// -------------------------------------------------------------
	g.dynoTxt, err = shapes.NewDynamicPixelTextNode("Ranger", world, g)
	if err != nil {
		return err
	}
	g.dynoTxt.SetScale(2.0)
	g.dynoTxt.SetPosition(-10.0, 10.0)
	g.dynoTxt.SetRotation(maths.DegreeToRadians * -45.0)
	fmt.Println(g.dynoTxt.Rotation())

	gd := g.dynoTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("Ranger is a Go!")
	gd.SetColor(color.NewPaletteInt64(color.LightPink).Array())
	gd.SetPixelSize(3.0)

	// -------------------------------------------------------------
	g.sqr, err = shapes.NewMonoSquareNode("Square", api.FILLOUTLINED, true, world, g)
	if err != nil {
		return err
	}
	g.sqr.SetScale(100.0)
	g.sqr.SetPosition(100.0, 0.0)
	gsq := g.sqr.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightOrange))
	gsq.SetFilledAlpha(0.5)

	// ---------------------------------------------------------------------
	g.line, err = shapes.NewDynamicMonoLineNode("DynoLine", world, g)
	if err != nil {
		return err
	}
	// g.line.SetScale(2.0)
	gl := g.line.(*shapes.DynamicMonoLineNode)
	gl.SetColor(color.NewPaletteInt64(color.LightOrange))
	gl.SetVertex1(100.0, 100.0)
	gl.SetVertex2(100.0, -100.0)

	// ---------------------------------------------------------------------
	g.line2, err = shapes.NewDynamicMonoLineNode("DynoLine2", world, g)
	if err != nil {
		return err
	}
	gl2 := g.line2.(*shapes.DynamicMonoLineNode)
	gl2.SetColor(color.NewPaletteInt64(color.Lime))

	g.line3, err = shapes.NewDynamicMonoLineNode("DynoLine3", world, g)
	if err != nil {
		return err
	}
	gl3 := g.line3.(*shapes.DynamicMonoLineNode)
	gl3.SetColor(color.NewPaletteInt64(color.Red))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {

	g.sqr.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 1.5

	angle := g.dynoTxt.Rotation() + (maths.DegreeToRadians * 1.5 / 10)
	g.dynoTxt.SetRotation(angle)

	if g.debugCnt > g.debug {
		g.debugCnt = 0
		gl := g.line.(*shapes.DynamicMonoLineNode)
		x := 100.0 + rand.Float32()*-20
		y := -100.0 + rand.Float32()*-20
		gl.SetVertex2(x, y)

		x = 100.0 + rand.Float32()*-50
		y = 100.0 + rand.Float32()*50
		gl.SetVertex1(x, y)

		gl3 := g.line3.(*shapes.DynamicMonoLineNode)
		x = 150.0 + rand.Float32()*100
		y = 150.0 + rand.Float32()*100
		gl3.SetVertex2(x, y)

		x = 150.0 + rand.Float32()*100
		y = 150.0 + rand.Float32()*100
		gl3.SetVertex1(x, y)
	}
	g.debugCnt += msPerUpdate

	gl2 := g.line2.(*shapes.DynamicMonoLineNode)
	x := -150.0 + rand.Float32()*100
	y := -150.0 + rand.Float32()*100
	gl2.SetVertex2(x, y)

	x = -150.0 + rand.Float32()*100
	y = -150.0 + rand.Float32()*100
	gl2.SetVertex1(x, y)

	dynoAtlas := g.World().GetAtlas(api.DynamicMonoAtlasName)
	dynoAtlas.(api.IDynamicAtlasX).Update()

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
