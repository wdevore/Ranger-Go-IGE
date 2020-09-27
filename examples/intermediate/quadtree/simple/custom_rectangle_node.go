package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type customRectangleNode struct {
	nodes.Node

	shape    api.IAtlasShape
	centered bool
	filled   bool

	color       []float32
	insideColor []float32

	mx, my        int32
	polygon       api.IPolygon
	localPosition api.IPoint
	pointInside   bool
}

func newCustomRectangleNode(name string, centered, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(customRectangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	o.centered = centered
	o.filled = filled
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
	c.insideColor = color.NewPaletteInt64(color.Orange).Array()

	if c.filled {
		name := "FilledSquare"
		if c.centered {
			name = "FilledCenteredSquare"
		}
		c.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineSquare"
		if c.centered {
			name = "OutlineCenteredSquare"
		}
		c.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	c.localPosition = geometry.NewPoint()
	c.polygon = geometry.NewPolygon()

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	c.populate()

	return nil
}

func (c *customRectangleNode) populate() {
	var vertices []float32

	if c.centered {
		vertices = []float32{
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0,
			0.5, 0.5, 0.0,
			-0.5, 0.5, 0.0,
		}
	} else {
		vertices = []float32{
			0.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			1.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
		}
	}

	// Populate polygon
	c.polygon.AddVertex(vertices[0], vertices[1])
	c.polygon.AddVertex(vertices[3], vertices[4])
	c.polygon.AddVertex(vertices[6], vertices[7])
	c.polygon.AddVertex(vertices[9], vertices[10])

	c.shape.SetVertices(vertices)

	var indices []uint32

	if c.filled {
		indices = []uint32{
			0, 1, 2,
			0, 2, 3,
		}
	} else {
		indices = []uint32{
			0, 1, 2, 3,
		}
	}

	c.shape.SetElementCount(len(indices))

	c.shape.SetIndices(indices)
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

// PointInside xxx
func (c *customRectangleNode) PointInside() bool {
	return c.pointInside
}

// -----------------------------------------------------
// Visuals
// -----------------------------------------------------

func (c *customRectangleNode) Draw(model api.IMatrix4) {
	if c.IsDirty() {
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

	renG.Render(c.shape, model)
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
