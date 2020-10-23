package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	sqr            api.INode
	mousePosTxt    api.INode
	localRecPosTxt api.INode

	viewPoint api.IPoint
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
	shline, err := shapes.NewMonoHLineNode("XAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*shapes.MonoHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	svline, err := shapes.NewMonoVLineNode("YAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	svline.SetScale(float32(dvr.Height))
	ghv := svline.(*shapes.MonoVLineNode)
	ghv.SetColor(color.NewPaletteInt64(color.LightGray))

	// ---------------------------------------------------------
	g.sqr, err = newCustomRectangleNode("CustomSqr", api.OUTLINED, true, world, g)

	// ---------------------------------------------------------
	g.mousePosTxt, err = shapes.NewDynamicPixelTextNode("MosPos", world, world.Overlay())
	if err != nil {
		return err
	}
	g.mousePosTxt.SetScale(1.0)
	g.mousePosTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.mousePosTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())
	gd.SetPixelSize(2.0)

	g.localRecPosTxt, err = shapes.NewDynamicPixelTextNode("LocPos", world, world.Overlay())
	if err != nil {
		return err
	}
	g.localRecPosTxt.SetScale(1.0)
	g.localRecPosTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-60.0))
	gd = g.localRecPosTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())
	gd.SetPixelSize(2.0)

	g.viewPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	text := fmt.Sprintf("(%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.mousePosTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText(text)

	glp := g.sqr.(*customRectangleNode)
	lp := glp.LocalPosition()
	text = fmt.Sprintf("(%7.3f, %7.3f)", lp.X(), lp.Y())
	gd = g.localRecPosTxt.(*shapes.DynamicPixelPixelTextNode)
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
