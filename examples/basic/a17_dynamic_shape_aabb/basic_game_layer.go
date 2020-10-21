package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	angle float64
	tri   api.INode

	line    api.INode
	dynAABB api.INode

	dynoTxt api.INode

	tVerts   []float32
	aabbRect api.IRectangle
	rotM4    api.IMatrix4
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

	// -------------------------------------------------------------
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	// -------------------------------------------------------------
	g.dynoTxt, err = shapes.NewDynamicPixelTextNode("Ranger", world, g)
	if err != nil {
		return err
	}
	g.dynoTxt.SetScale(2.0)
	g.dynoTxt.SetPosition(-10.0, 10.0)
	g.dynoTxt.SetRotation(maths.DegreeToRadians * 45.0)
	fmt.Println(g.dynoTxt.Rotation())

	gd := g.dynoTxt.(*shapes.DynamicPixelPixelTextNode)
	gd.SetText("Ranger is a Go!")
	gd.SetColor(color.NewPaletteInt64(color.LightPink).Array())
	gd.SetPixelSize(3.0)

	// -------------------------------------------------------------
	g.tri, err = shapes.NewMonoTriangleNode("Tri", api.FILLED, world, g)
	if err != nil {
		return err
	}
	g.tri.SetScale(200.0)
	g.tri.SetPosition(100.0, 100.0)
	gsq := g.tri.(*shapes.MonoTriangleNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.GoldYellow))
	gsq.SetFilledAlpha(0.5)

	// ---------------------------------------------------------------------
	// There are two ways we can set the parent of the AABB:
	// One is "g" and the other is "g.tri". If we set the parent to "g"
	// then we need to apply Scale and Translation manually to sync
	// with the triangle.
	// If we use "g.tri" then we can dispense with Scale and Translation
	// but we will need to remove the inherited rotation in the Update()
	// method.
	g.dynAABB, err = shapes.NewDynamicMonoSquareNode("DynoSquare", true, false, world, g)
	if err != nil {
		return err
	}
	// The AABB doesn't inherit from the triangle so we manually sync to
	// the triangle.
	g.dynAABB.SetScale(g.tri.Scale())
	g.dynAABB.SetPosition(g.tri.Position().X(), g.tri.Position().Y())

	gs := g.dynAABB.(*shapes.DynamicMonoSquareNode)
	gs.SetColor(color.NewPaletteInt64(color.White))
	gs.SetAlpha(0.5)

	// Note: local-space is being defined as the "Unit"-space but technically
	// it is model-space. I keep the two distinct, local-space (aka unit-space)
	// verse model-space.
	// These are the local-space transformed vertices of the triangle.
	g.tVerts = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}
	// This rectangle is resized based on the local-space transformed vertices.
	g.aabbRect = geometry.NewRectangle()

	// This is the rotation applied to the triangle
	g.rotM4 = maths.NewMatrix4()

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	rotAngle := maths.DegreeToRadians * g.angle / 5

	g.tri.SetRotation(rotAngle)

	angle := g.dynoTxt.Rotation() - (maths.DegreeToRadians * 1.5 / 10)
	g.dynoTxt.SetRotation(angle)

	// If the AABB inherited from the Triangle then we need to remove
	// the triangles rotation.
	// g.dynAABB.SetRotation(-rotAngle)

	gs := g.dynAABB.(*shapes.DynamicMonoSquareNode)

	// Calc local-space using the triangle's rotation.
	g.rotM4.SetRotation(rotAngle)

	gtri := g.tri.(*shapes.MonoTriangleNode)
	// Transform the original triangle vertices while remaining
	// in local-space.
	g.rotM4.TransformVertices3D(*gtri.Vertices(), g.tVerts)

	// Update AABB. This will adjust the rectangle to "fit"
	// the transformed vertices.
	g.aabbRect.SetBounds3D(g.tVerts)

	// Now update the visual to match the vertices.
	gs.SetLowerLeft(g.aabbRect.Left(), g.aabbRect.Bottom())
	gs.SetLowerRight(g.aabbRect.Right(), g.aabbRect.Bottom())
	gs.SetUpperRight(g.aabbRect.Right(), g.aabbRect.Top())
	gs.SetUpperLeft(g.aabbRect.Left(), g.aabbRect.Top())

	// And finally update the atlas so that the backing buffer
	// is copied to the GL buffer.
	dynoAtlas := g.World().GetAtlas(api.DynamicMonoAtlasName)
	dynoAtlas.(api.IDynamicAtlasX).Update()

	g.angle -= 1.5
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
}
