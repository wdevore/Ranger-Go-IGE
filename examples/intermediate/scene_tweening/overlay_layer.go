package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type overlayLayer struct {
	nodes.Node

	timing     api.INode
	viewLocTxt api.INode
	viewPoint  api.IPoint
}

func newOverlayLayer(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(overlayLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

func (g *overlayLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	var err error
	dvr := world.Properties().Window.DeviceRes

	// ---------------------------------------------------------
	g.viewLocTxt, err = custom.NewDynamicTextNode("Text", 500, world, g)
	if err != nil {
		return err
	}
	g.viewLocTxt.SetScale(2.0)
	g.viewLocTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.viewLocTxt.(*custom.DynamicTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.White).Array())
	gd.SetPixelSize(1.0)

	if world.Properties().Engine.ShowTimingInfo {
		g.timing, err = custom.NewDynamicTextNode("TimingInfo", 500, world, g)
		if err != nil {
			return err
		}
		g.timing.SetScale(1.0)
		// Set position to lower-left corner
		dvr := world.Properties().Window.DeviceRes
		g.timing.SetPosition(float32(-dvr.Width/2+10.0), float32(-dvr.Height/2)+10.0)

		gt2 := g.timing.(*custom.DynamicTextNode)
		gt2.SetText("-")
		gt2.SetPixelSize(2.0)
		gt2.SetColor(color.NewPaletteInt64(color.LightOrange).Array())
	}

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

	g.viewPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *overlayLayer) Update(msPerUpdate, secPerUpdate float64) {
	text := fmt.Sprintf("(%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.viewLocTxt.(*custom.DynamicTextNode)
	gd.SetText(text)

	w := g.World()
	if w.Properties().Engine.ShowTimingInfo {
		gt := g.timing.(*custom.DynamicTextNode)
		s := fmt.Sprintf("f:%d u:%d r:%2.3f", w.Fps(), w.Ups(), w.AvgRender())
		gt.SetText(s)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *overlayLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *overlayLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
	man.UnRegisterEventTarget(g)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (g *overlayLayer) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeMouseMotion {
		mx, my := event.GetMousePosition()
		nodes.MapDeviceToView(g.World(), mx, my, g.viewPoint)
	}

	return false
}
