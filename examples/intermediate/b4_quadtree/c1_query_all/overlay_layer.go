package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type overlayLayer struct {
	nodes.Node

	underLay       api.INode
	underLayPoint  api.IPoint
	underLayLocTxt api.INode

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
	g.viewLocTxt, err = shapes.NewDynamicPixelTextNode("MosPos", world, world.Overlay())
	if err != nil {
		return err
	}
	g.viewLocTxt.SetScale(1.0)
	g.viewLocTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.viewLocTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())
	gd.SetPixelSize(2.0)

	g.underLayLocTxt, err = shapes.NewDynamicPixelTextNode("LocPos", world, world.Overlay())
	if err != nil {
		return err
	}
	g.underLayLocTxt.SetScale(1.0)
	g.underLayLocTxt.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-60.0))
	gd = g.underLayLocTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())
	gd.SetPixelSize(2.0)

	g.viewPoint = geometry.NewPoint()
	g.underLayPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *overlayLayer) Update(msPerUpdate, secPerUpdate float64) {
	text := fmt.Sprintf("V (%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.viewLocTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText(text)

	text = fmt.Sprintf("G (%d, %d)", int(g.underLayPoint.X()), int(g.underLayPoint.Y()))
	gd = g.underLayLocTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText(text)
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
