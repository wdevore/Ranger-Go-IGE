package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	dragSquare *draggableSquare
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

	var err error

	dvr := world.Properties().Window.DeviceRes

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
	tri, err := custom.NewStaticTriangleNode("FilledTri", true, true, world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100)
	tri.SetPosition(150.0, 0.0)
	gtr := tri.(*custom.StaticTriangleNode)
	gtr.SetColor(color.NewPaletteInt64(color.DeepPink))

	// ---------------------------------------------------------
	otri, err := custom.NewStaticTriangleNode("OutlineTri", true, false, world, g)
	if err != nil {
		return err
	}
	otri.SetScale(100)
	otri.SetPosition(150.0, 0.0)
	gotr := otri.(*custom.StaticTriangleNode)
	gotr.SetColor(color.NewPaletteInt64(color.White))

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(world, g)

	return nil
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	// We want the mouse events so the node can track the mouse.
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterEventTarget(g)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (g *gameLayer) Handle(event api.IEvent) bool {
	handled := g.dragSquare.SquareEventHandle(event)

	return handled
}
