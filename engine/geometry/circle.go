package geometry

import (
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Circle is a circle with behaviours
type Circle struct {
	radius float32

	center api.IPoint
}

// NewCircle creates a circle
func NewCircle() *Circle {
	o := new(Circle)
	o.radius = 1.0
	o.center = NewPoint()
	return o
}

// SetRadius of circle
func (c *Circle) SetRadius(radius float32) {
	c.radius = radius
}

// Radius of circle
func (c *Circle) Radius() float32 {
	return c.radius
}

// SetCenter of circle
func (c *Circle) SetCenter(x, y float32) {
	c.center.SetByComp(x, y)
}

// Center of circle
func (c *Circle) Center() api.IPoint {
	return c.center
}

// PointInside checks point to radius
func (c *Circle) PointInside(point api.IPoint) bool {
	distance := c.DistanceFromCenter(point)
	return distance < c.radius
}

// DistanceFromCenter returns distance from point to circle center
func (c *Circle) DistanceFromCenter(point api.IPoint) float32 {
	dx := c.center.X() - point.X()
	dy := c.center.Y() - point.Y()
	return float32(math.Sqrt(float64(dx*dx) + float64(dy*dy)))
}

// DistanceFromEdge returns distance from point to circle edge
// if <= 0 then on inside edge.
func (c *Circle) DistanceFromEdge(point api.IPoint) float32 {
	dx := c.center.X() - point.X()
	dy := c.center.Y() - point.Y()
	distance := math.Sqrt(float64(dx*dx) + float64(dy*dy))
	return float32(distance) - c.radius
}
