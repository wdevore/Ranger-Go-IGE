package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/quadtree"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	dragSquare *draggableSquare

	// Structual tree
	tree api.IQuadTree

	quadTreeNode api.INode
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

	g.buildQuadtree(world)

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(10.0, world, g)

	err := g.buildOriginAxies(world)
	if err != nil {
		return err
	}

	g.tree.Add(g.dragSquare.BaseNode())

	return nil
}

func (g *gameLayer) buildQuadtree(world api.IWorld) error {
	var err error
	scale := float32(400.0)

	g.tree = quadtree.NewQuadTree()
	g.tree.SetBoundaryByMinMax(-scale, -scale, scale, scale)
	g.tree.SetMaxDepth(6)

	g.quadTreeNode, err = NewQTreeNode("Quadtree", world, g)
	if err != nil {
		return err
	}

	gqt := g.quadTreeNode.(*QTreeNode)
	gqt.SetColor(color.NewPaletteInt64(color.LightOrange))
	gqt.SetTree(g.tree)

	return nil
}

func (g *gameLayer) buildOriginAxies(world api.IWorld) error {
	var err error

	dvr := world.Properties().Window.DeviceRes

	// ---------------------------------------------------------
	// Instead of using two node: vline and hline, I'm using one "+ node.
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

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
	handled := g.dragSquare.EventHandle(event)

	if handled {
		// fmt.Println(g.dragSquare.BaseNode().Position())
		g.tree.Remove(g.dragSquare.BaseNode())
		// TODO Clean() should work better but it doesn't. Research it!
		g.tree.Clear()
		g.tree.Add(g.dragSquare.BaseNode())
	}

	// TODO fix this. I shouldn't need to set this to false
	// Something isn't right with the event system.
	// The drag pauses and trips and it shouldn't
	handled = false

	return handled
}
