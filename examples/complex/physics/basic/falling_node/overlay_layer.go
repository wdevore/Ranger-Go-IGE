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

	underLay       api.INode
	underLayPoint  api.IPoint
	underLayLocTxt api.INode

	timing     api.INode
	viewLocTxt api.INode
	viewPoint  api.IPoint
}

func newOverlayLayer(name string, world api.IWorld, underLay, parent api.INode) (api.INode, error) {
	o := new(overlayLayer)
	o.Initialize(name)
	o.SetParent(parent)
	o.underLay = underLay
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}
	return o, nil
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
	g.viewLocTxt.SetScale(1.0)
	g.viewLocTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.viewLocTxt.(*custom.DynamicTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.PanSkin))
	gd.SetPixelSize(2.0)

	// ---------------------------------------------------------
	g.underLayLocTxt, err = custom.NewDynamicTextNode("Text", 500, world, g)
	if err != nil {
		return err
	}
	g.underLayLocTxt.SetScale(1.0)
	g.underLayLocTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-60.0))
	gud := g.underLayLocTxt.(*custom.DynamicTextNode)
	gud.SetText("(0,0)")
	gud.SetColor(color.NewPaletteInt64(color.PanSkin))
	gud.SetPixelSize(2.0)

	// ---------------------------------------------------------
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
		gt2.SetColor(color.NewPaletteInt64(color.LightOrange))
	}

	// ---------------------------------------------------------
	shline, err := custom.NewStaticHLineNode("HLine", world, g)
	if err != nil {
		return err
	}
	shline.SetScale(float32(dvr.Width))
	ghl := shline.(*custom.StaticHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.White))

	// ---------------------------------------------------------
	svline, err := custom.NewStaticVLineNode("VLine", world, g)
	if err != nil {
		return err
	}
	svline.SetScale(float32(float32(dvr.Height)))
	gvl := svline.(*custom.StaticVLineNode)
	gvl.SetColor(color.NewPaletteInt64(color.White))

	g.viewPoint = geometry.NewPoint()
	g.underLayPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *overlayLayer) Update(msPerUpdate, secPerUpdate float64) {
	text := fmt.Sprintf("V (%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.viewLocTxt.(*custom.DynamicTextNode)
	gd.SetText(text)

	text = fmt.Sprintf("G (%d, %d)", int(g.underLayPoint.X()), int(g.underLayPoint.Y()))
	gd = g.underLayLocTxt.(*custom.DynamicTextNode)
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

		nodes.MapDeviceToNode(mx, my, g.underLay, g.underLayPoint)
	}

	return false
}
