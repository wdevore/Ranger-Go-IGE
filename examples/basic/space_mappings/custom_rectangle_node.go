package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type customRectangleNode struct {
	nodes.Node

	background api.IAtlasShape

	color       []float32
	insideColor []float32

	mx, my        int32
	polygon       api.IPolygon
	localPosition api.IPoint
	pointInside   bool
}

func newCustomRectangleNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(customRectangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (c *customRectangleNode) Build(world api.IWorld) error {
	c.Node.Build(world)

	c.SetDirty(true)

	c.color = color.NewPaletteInt64(color.DarkerGray).Array()
	c.insideColor = color.NewPaletteInt64(color.Red).Array()

	c.background = world.Atlas().Shape("CenteredSquare")

	c.polygon = geometry.NewPolygon()

	// Use the shape's data to setup polygon
	v := c.background.Vertices(0)

	// Subtract 1 because the offset always is 1 beyond
	// Subtract count*2 because each vertex is 6 floats 2(xyz) foreach
	// face of which there are 2.
	o := (c.background.Offset() - 1) - c.background.Count()*2

	c.polygon.AddVertex(v[0+o], v[1+o])
	c.polygon.AddVertex(v[3+o], v[4+o])
	c.polygon.AddVertex(v[6+o], v[7+o])
	c.polygon.AddVertex(v[9+o], v[10+o])

	// c.polygon.AddVertex(-0.5, -0.5)
	// c.polygon.AddVertex(-0.5, 0.5)
	// c.polygon.AddVertex(0.5, 0.5)
	// c.polygon.AddVertex(0.5, -0.5)

	c.localPosition = geometry.NewPoint()

	return nil
}

// LocalPosition xxx
func (c *customRectangleNode) LocalPosition() api.IPoint {
	return c.localPosition
}

// SetColor sets line color
func (c *customRectangleNode) SetColor(color api.IPalette) {
	c.color = color.Array()
}

// EnterNode called when a node is entering the stage
func (c *customRectangleNode) EnterNode(man api.INodeManager) {
	man.RegisterEventTarget(c)
}

// ExitNode called when a node is exiting stage
func (c *customRectangleNode) ExitNode(man api.INodeManager) {
	man.UnRegisterEventTarget(c)
}

// -----------------------------------------------------
// Visuals
// -----------------------------------------------------

func (c *customRectangleNode) Draw(model api.IMatrix4) {
	if c.IsDirty() {
		// p := maths.NewVector3With2Components(-0.5, -0.5)
		// p.Mul(model)
		// c.polygon.SetVertex(p.X(), p.Y(), 0)

		// p.Set2Components(-0.5, 0.5)
		// p.Mul(model)
		// c.polygon.SetVertex(p.X(), p.Y(), 1)

		// p.Set2Components(0.5, 0.5)
		// p.Mul(model)
		// c.polygon.SetVertex(p.X(), p.Y(), 2)

		// p.Set2Components(0.5, -0.5)
		// p.Mul(model)
		// c.polygon.SetVertex(p.X(), p.Y(), 3)
		// fmt.Println(c.polygon)
		c.SetDirty(false)
	}

	// This gets the local-space coords of the rectangle node.
	// Note: OpenGL's +Y axis is upward
	nodes.MapDeviceToNode(c.mx, c.my, c, c.localPosition)

	c.pointInside = c.polygon.PointInside(c.localPosition)

	renG := c.World().UseRenderGraphic(api.StaticRenderGraphic)
	if c.pointInside {
		// fmt.Println("Inside: ", c.localPosition)
		renG.SetColor(c.insideColor)
	} else {
		// fmt.Println(c.localPosition)
		renG.SetColor(c.color)
	}

	renG.Render(c.background, model)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

// Handle events from IO
func (c *customRectangleNode) Handle(event api.IEvent) bool {
	// fmt.Println(event)
	if event.GetType() == api.IOTypeMouseMotion {
		c.mx, c.my = event.GetMousePosition()
	}

	return false
}
