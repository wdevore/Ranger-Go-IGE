package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

type draggable struct {
	mx, my        int32
	polygon       api.IPolygon
	localPosition api.IPoint
	pointInside   bool
}

func newDraggable(centered bool) *draggable {
	o := new(draggable)

	if err := o.build(centered); err != nil {
		return nil
	}

	return o
}

func (c *draggable) build(centered bool) error {
	c.localPosition = geometry.NewPoint()
	c.polygon = geometry.NewPolygon()

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	c.populate(centered)

	return nil
}

func (c *draggable) populate(centered bool) {
	vertices, _, _ := generators.GenerateUnitRectangleVectorShape(centered, true)

	// Populate polygon
	c.polygon.AddVertex(vertices[0], vertices[1])
	c.polygon.AddVertex(vertices[3], vertices[4])
	c.polygon.AddVertex(vertices[6], vertices[7])
	c.polygon.AddVertex(vertices[9], vertices[10])
}

// LocalPosition --
func (c *draggable) LocalPosition() api.IPoint {
	return c.localPosition
}

// PointInside --
func (c *draggable) PointInside() bool {
	return c.pointInside
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

// Handle events from IO
func (c *draggable) Handle(node api.INode, event api.IEvent) bool {
	// fmt.Println(event)
	if event.GetType() == api.IOTypeMouseMotion {
		c.mx, c.my = event.GetMousePosition()

		// This gets the local-space coords of the rectangle node.
		// Note: OpenGL's +Y axis is upward
		nodes.MapDeviceToNode(c.mx, c.my, node, c.localPosition)

		c.pointInside = c.polygon.PointInside(c.localPosition)

		// if c.pointInside {
		// 	fmt.Println("Inside: ", c.localPosition)
		// } else {
		// 	fmt.Println(c.localPosition)
		// }
	}

	return false
}
