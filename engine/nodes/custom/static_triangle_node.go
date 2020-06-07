package custom

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticTriangleNode is a generic node
type StaticTriangleNode struct {
	nodes.Node

	color    []float32
	halfSide float32

	vertices []float32

	shape    api.IAtlasShape
	centered bool
	filled   bool
}

// NewStaticTriangleNode constructs a generic shape node
func NewStaticTriangleNode(name string, centered, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticTriangleNode)
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

// Build configures the node
func (r *StaticTriangleNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := "FilledTriangle"
		if r.centered {
			name = "FilledCenteredTriangle"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineTriangle"
		if r.centered {
			name = "OutlineCenteredTriangle"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticTriangleNode) populate() {
	const centerOffset = float32(math.Pi / 4 / 10)

	const top = float32(math.Pi / 10)

	if r.centered {
		r.halfSide = 0.5
		r.vertices = []float32{
			-r.halfSide, -r.halfSide + centerOffset, 0.0,
			r.halfSide, -r.halfSide + centerOffset, 0.0,
			0.0, 0.314 + centerOffset, 0.0,
		}
	} else {
		r.halfSide = 1.0
		r.vertices = []float32{
			0.0, 0.0 + centerOffset, 0.0,
			r.halfSide, 0.0 + centerOffset, 0.0,
			0.0, top + centerOffset, 0.0,
		}
	}

	r.shape.SetVertices(r.vertices)

	var indices []uint32

	indices = []uint32{
		0, 1, 2,
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Vertices returns the shape's vertices
func (r *StaticTriangleNode) Vertices() []float32 {
	return r.vertices
}

// SideLength returns the scale length
func (r *StaticTriangleNode) SideLength() float32 {
	return r.halfSide * r.Scale() * 2
}

// HalfSide returns the scaled half side length.
func (r *StaticTriangleNode) HalfSide() float32 {
	return r.halfSide * r.Scale()
}

// SetColor sets line color
func (r *StaticTriangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticTriangleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticTriangleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
