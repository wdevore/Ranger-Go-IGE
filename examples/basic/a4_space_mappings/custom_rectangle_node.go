package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type customRectangleNode struct {
	nodes.Node

	visual api.INode
	angle  float64

	centered        bool
	drawStyle       int
	rotationEnabled bool

	color       api.IPalette
	insideColor api.IPalette

	mx, my        int32
	polygon       api.IPolygon
	localPosition api.IPoint
	pointInside   bool
}

func newCustomRectangleNode(name string, drawStyle int, centered bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(customRectangleNode)

	o.Initialize(name)
	o.SetParent(parent)
	o.centered = centered
	o.drawStyle = drawStyle

	parent.AddChild(o)
	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (c *customRectangleNode) build(world api.IWorld) error {
	c.Node.Build(world)

	c.SetDirty(true)

	c.color = color.NewPaletteInt64(color.LightOrange)
	c.insideColor = color.NewPaletteInt64(color.Red)

	var err error

	c.visual, err = shapes.NewMonoSquareNode("VisualSqr", c.drawStyle, c.centered, world, c)
	if err != nil {
		return err
	}
	c.visual.SetScale(100)
	c.visual.SetPosition(100.0, 100.0)
	c.angle = 35.0
	c.visual.SetRotation(maths.DegreeToRadians * c.angle)

	c.localPosition = geometry.NewPoint()
	c.polygon = geometry.NewPolygon()

	c.populate()

	return nil
}

func (c *customRectangleNode) populate() {
	vertices, _, _ := generators.GenerateUnitRectangleVectorShape(c.centered, true)

	// Populate polygon
	c.polygon.AddVertex(vertices[0], vertices[1])
	c.polygon.AddVertex(vertices[3], vertices[4])
	c.polygon.AddVertex(vertices[6], vertices[7])
	c.polygon.AddVertex(vertices[9], vertices[10])
}

// LocalPosition xxx
func (c *customRectangleNode) LocalPosition() api.IPoint {
	return c.localPosition
}

// EnterNode called when a node is entering the stage
func (c *customRectangleNode) EnterNode(man api.INodeManager) {
	man.RegisterTarget(c)
	man.RegisterEventTarget(c)
}

// ExitNode called when a node is exiting stage
func (c *customRectangleNode) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(c)
	man.UnRegisterEventTarget(c)
}

// -----------------------------------------------------
// Visuals
// -----------------------------------------------------

func (c *customRectangleNode) Update(msPerUpdate, secPerUpdate float64) {
	if c.rotationEnabled {
		c.visual.SetRotation(maths.DegreeToRadians * c.angle)
		c.angle -= 1.5
	}

	// This gets the local-space coords of the rectangle node.
	// Note: OpenGL's +Y axis is upward
	nodes.MapDeviceToNode(c.mx, c.my, c.visual, c.localPosition)

	c.pointInside = c.polygon.PointInside(c.localPosition)

	if c.pointInside {
		// fmt.Println("Inside: ", c.localPosition)
		gb := c.visual.(*shapes.MonoSquareNode)
		gb.SetOutlineColor(c.insideColor)
	} else {
		// fmt.Println(c.localPosition)
		gb := c.visual.(*shapes.MonoSquareNode)
		gb.SetOutlineColor(c.color)
	}
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

// Handle events from IO
func (c *customRectangleNode) Handle(event api.IEvent) bool {
	// fmt.Println(event)
	if event.GetType() == api.IOTypeMouseMotion {
		c.mx, c.my = event.GetMousePosition()
	} else if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		switch event.GetKeyScan() {
		case 82: // r
			if event.GetState() == 1 {
				c.rotationEnabled = !c.rotationEnabled
			}
		case 48: // 0
			if event.GetState() == 1 {
				c.angle = 0.0
				c.visual.SetRotation(maths.DegreeToRadians * c.angle)
			}
		}
	}

	return false
}
