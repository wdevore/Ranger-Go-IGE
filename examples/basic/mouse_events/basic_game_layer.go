package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64

	sqr     api.INode
	dynoTxt api.INode
	line    api.INode

	viewPoint api.IPoint
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
	g.sqr, err = custom.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}
	g.sqr.SetScale(100.0)
	g.sqr.SetPosition(110.0, 100.0)
	gol2 := g.sqr.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightPurple))

	// ---------------------------------------------------------
	g.dynoTxt, err = custom.NewDynamicTextNode("Text", 500, world, g)
	if err != nil {
		return err
	}
	g.dynoTxt.SetScale(2.0)
	g.dynoTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.dynoTxt.(*custom.DynamicPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.White).Array())
	gd.SetPixelSize(1.0)

	g.viewPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.sqr.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 1.5

	text := fmt.Sprintf("(%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.dynoTxt.(*custom.DynamicPixelTextNode)
	gd.SetText(text)
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
	man.UnRegisterEventTarget(g)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (g *gameLayer) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeMouseMotion {
		mx, my := event.GetMousePosition()
		nodes.MapDeviceToView(g.World(), mx, my, g.viewPoint)
	}

	return false
}
