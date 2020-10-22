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

	// ---------------------------------------------------------
	// Instead of using two node: vline and hline, I'm using one "+ node.
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, g)
	if err != nil {
		return err
	}
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(world, g)

	// ---------------------------------------------------------
	tri, err := shapes.NewMonoTriangleNode("Tri", api.FILLED, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(150.0, 0.0)
	gsq := tri.(*shapes.MonoTriangleNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.DeepPink))
	gsq.SetFilledAlpha(0.5)

	// ---------------------------------------------------------
	g.viewCoords, err = shapes.NewDynamicPixelTextNode("ViewCoords", world, g)
	if err != nil {
		return err
	}
	g.viewCoords.SetScale(2.0)
	g.viewCoords.SetPosition(-float32(dvr.Width/2)+20.0, float32(dvr.Height/2-30.0))
	gd := g.viewCoords.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("(0,0)")
	gd.SetColor(color.NewPaletteInt64(color.GreenYellow).Array())
	gd.SetPixelSize(1.0)

	// ---------------------------------------------------------
	g.plus, err = shapes.NewMonoPlusNode("Plus", world, g)
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
	gd := g.viewCoords.(*shapes.DynamicPixelPixelTextNode)
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
