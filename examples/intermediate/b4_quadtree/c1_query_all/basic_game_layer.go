package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/quadtree"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	dragSquare *draggableSquare

	tree         api.IQuadTree
	quadTreeNode api.INode
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

	g.buildQuadtree(world)

	// Square ----------------------------------------------------
	g.dragSquare = newDraggableSquare()
	g.dragSquare.Build(50.0, world, g)

	g.buildTriangles(world)

	err := g.buildOriginAxies(world)
	if err != nil {
		return err
	}

	// g.tree.Add(g.dragSquare.BaseNode())

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
	gqt.SetColor(color.NewPaletteInt64(color.Gray))
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

func (g *gameLayer) buildTriangles(world api.IWorld) error {
	poss := []float32{
		// 0.0, 0.0,
		150.0, 150.0,
		15.0, 15.0,
		-125.0, 150.0,
		-145.0, 130.0,
	}

	// ---------------------------------------------------------
	scale := float32(10.0)
	for i := 0; i < len(poss); i += 2 {
		tri, err := shapes.NewMonoTriangleNode("::Tri", api.FILLED, world, g)
		if err != nil {
			return err
		}
		tri.SetScale(scale)
		tri.SetPosition(poss[i], poss[i+1])
		tri.SetBoundByMinMax(
			tri.Position().X()-scale/2.0,
			tri.Position().Y()-scale/2.0,
			tri.Position().X()+scale/2.0,
			tri.Position().Y()+scale/2.0,
		)

		gsq := tri.(*shapes.MonoTriangleNode)
		gsq.SetFilledColor(color.NewPaletteInt64(color.LightOrange))
		gsq.SetFilledAlpha(0.5)

		g.tree.Add(tri)
	}

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
		bounds := g.dragSquare.BaseNode().Bounds()
		nodes := []api.INode{}
		g.tree.Query(bounds, &nodes)
		if len(nodes) > 0 {
			fmt.Println(nodes)
		}
	}

	// TODO fix this. I shouldn't need to set this to false
	// Something isn't right with the event system.
	// The drag pauses and trips and it shouldn't
	handled = false

	return handled
}
