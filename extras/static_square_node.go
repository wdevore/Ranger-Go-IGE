package extras

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticSquareNode is a generic node
type StaticSquareNode struct {
	nodes.Node

	color    []float32
	halfSide float32

	shape    api.IAtlasShape
	centered bool
	filled   bool
}

// NewStaticSquareNode constructs a generic shape node
func NewStaticSquareNode(name string, centered, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticSquareNode)
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
func (r *StaticSquareNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := "FilledSquare"
		if r.centered {
			name = "FilledCenteredSquare"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineSquare"
		if r.centered {
			name = "OutlineCenteredSquare"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticSquareNode) populate() {
	var vertices []float32

	if r.centered {
		r.halfSide = 0.5
		vertices = []float32{
			-r.halfSide, -r.halfSide, 0.0,
			r.halfSide, -r.halfSide, 0.0,
			r.halfSide, r.halfSide, 0.0,
			-r.halfSide, r.halfSide, 0.0,
		}
	} else {
		r.halfSide = 1.0
		vertices = []float32{
			0.0, 0.0, 0.0,
			r.halfSide, 0.0, 0.0,
			r.halfSide, r.halfSide, 0.0,
			0.0, r.halfSide, 0.0,
		}
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

// HalfSide returns the scaled half side length.
func (r *StaticSquareNode) HalfSide() float32 {
	return r.halfSide * r.Scale()
}

// SetColor sets square color
func (r *StaticSquareNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticSquareNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticSquareNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor(r.color)
	renG.Render(r.shape, model)
}