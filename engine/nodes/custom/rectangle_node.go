package custom

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// RectangleNode is a basic rectangle
type RectangleNode struct {
	nodes.Node

	color api.IPalette

	polygon api.IPolygon
}

// NewRectangleNode constructs a triangle shaped node
func NewRectangleNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(RectangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (t *RectangleNode) Build(world api.IWorld) error {
	t.Node.Build(world)

	t.polygon = geometry.NewPolygon()
	t.polygon.AddVertex(-0.5, 0.5)
	t.polygon.AddVertex(0.5, 0.5)
	t.polygon.AddVertex(0.0, -0.314)

	t.polygon.Build()

	t.color = rendering.NewPaletteInt64(rendering.White)

	return nil
}

// Polygon returns the internal polygon mesh
func (t *RectangleNode) Polygon() api.IPolygon {
	return t.polygon
}

// SetColor sets line color
func (t *RectangleNode) SetColor(color api.IPalette) {
	t.color = color
}

// SetPoints sets the edge points of the triangle
func (t *RectangleNode) SetPoints(x1, y1, x2, y2, x3, y3 float32) {
	t.polygon.SetVertex(x1, y1, 0)
	t.polygon.SetVertex(x2, y2, 1)
	t.polygon.SetVertex(x3, y3, 2)
}

// Draw renders shape
func (t *RectangleNode) Draw(m4 api.IMatrix4) {
	if t.IsDirty() {
		// context.TransformPolygon(t.polygon)

		t.SetDirty(false)
	}

	// context.SetDrawColor(t.color)
	// context.RenderPolygon(t.polygon, api.CLOSED)
}
