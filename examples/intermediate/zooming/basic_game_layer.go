package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64

	sqr        api.INode
	viewCoords api.INode
	line       api.INode
	plus       api.INode

	viewPoint api.IPoint

	dragSquare *draggableSquare

	// Zooming
	zoom api.INode
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

	g.zoom, err = custom.NewZoomNode("ZoomNode", world, g)

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(world, g.zoom)

	// ---------------------------------------------------------
	tri, err := custom.NewStaticTriangleNode("FilledTri", true, true, world, g.zoom)
	if err != nil {
		return err
	}
	tri.SetScale(100)
	tri.SetPosition(150.0, 0.0)
	gtr := tri.(*custom.StaticTriangleNode)
	gtr.SetColor(color.NewPaletteInt64WithAlpha(color.DeepPink, 0.5))

	// ---------------------------------------------------------
	g.viewCoords, err = custom.NewDynamicTextNode("ViewCoords", 200, world, g)
	if err != nil {
		return err
	}
	g.viewCoords.SetScale(2.0)
	g.viewCoords.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.viewCoords.(*custom.DynamicTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GreenYellow))
	gd.SetPixelSize(1.0)

	// ---------------------------------------------------------
	g.plus, err = custom.NewStaticPlusNode("Plus", world, g)
	if err != nil {
		return err
	}
	g.plus.SetScale(30.0)

	g.viewPoint = geometry.NewPoint()

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	text := fmt.Sprintf("(%d, %d)", int(g.viewPoint.X()), int(g.viewPoint.Y()))
	gd := g.viewCoords.(*custom.DynamicTextNode)
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
	handled := g.dragSquare.EventHandle(event)

	if !handled {
		if event.GetType() == api.IOTypeMouseMotion {
			mx, my := event.GetMousePosition()
			nodes.MapDeviceToView(g.World(), mx, my, g.viewPoint)

			g.plus.SetPosition(g.viewPoint.X(), g.viewPoint.Y())
		}
	}

	return false
}
