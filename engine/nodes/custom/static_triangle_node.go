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

	color []float32

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
	var vertices []float32

	const centerOffset = float32(math.Pi / 4 / 10)
	const top = float32(math.Pi / 10)

	if r.centered {
		vertices = []float32{
			-0.5, -0.5 + centerOffset, 0.0,
			0.5, -0.5 + centerOffset, 0.0,
			0.0, 0.314 + centerOffset, 0.0,
		}
	} else {
		vertices = []float32{
			0.0, 0.0 + centerOffset, 0.0,
			1.0, 0.0 + centerOffset, 0.0,
			0.0, top + centerOffset, 0.0,
		}
	}

	r.shape.SetVertices(vertices)

	var indices []uint32

	if r.filled {
		indices = []uint32{
			0, 1, 2,
		}
	} else {
		indices = []uint32{
			0, 1, 2,
		}
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// SetColor sets line color
func (r *StaticTriangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticTriangleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticTriangleNode) PointInside(p api.IPoint) bool {
	return false
}

// Draw renders shape
func (r *StaticTriangleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
