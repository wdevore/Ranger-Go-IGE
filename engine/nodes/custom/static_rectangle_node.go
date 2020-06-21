package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticRectangleNode is a generic node
type StaticRectangleNode struct {
	nodes.Node

	color                  []float32
	minX, minY, maxX, maxY float32

	shape  api.IAtlasShape
	filled bool
}

// NewStaticRectangleNode constructs a generic shape node
func NewStaticRectangleNode(minX, minY, maxX, maxY float32, name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticRectangleNode)
	o.Initialize(name)
	o.SetParent(parent)

	o.filled = filled
	o.minX = minX
	o.minY = minY
	o.maxX = maxX
	o.maxY = maxY

	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticRectangleNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := "FilledRectangle"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineRectangle"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate(r.minX, r.minY, r.maxX, r.maxY)

	return nil
}

func (r *StaticRectangleNode) populate(minX, minY, maxX, maxY float32) {
	var vertices []float32

	vertices = []float32{
		minX, maxY, 0.0,
		minX, minY, 0.0,
		maxX, minY, 0.0,
		maxX, maxY, 0.0,
	}

	r.shape.SetVertices(vertices)

	var indices []uint32

	if r.filled {
		indices = []uint32{
			0, 1, 2,
			0, 2, 3,
		}
	} else {
		indices = []uint32{
			0, 1, 2, 3,
		}
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Vertices returns the shape's vertices
func (r *StaticRectangleNode) Vertices() []float32 {
	return r.shape.Vertices()
}

// HorizontalLength returns the maxX-minX length.
func (r *StaticRectangleNode) HorizontalLength() float32 {
	return r.maxX - r.minX
}

// VerticalLength returns the maxY-minY length.
func (r *StaticRectangleNode) VerticalLength() float32 {
	return r.maxY - r.minY
}

// SetColor sets square color
func (r *StaticRectangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticRectangleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticRectangleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
