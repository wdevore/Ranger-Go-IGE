package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticZBarNode is a generic node
type StaticZBarNode struct {
	nodes.Node

	color []float32

	shape  api.IAtlasShape
	filled bool
}

// NewStaticZBarNode constructs a generic shape node
func NewStaticZBarNode(name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticZBarNode)
	o.Initialize(name)
	o.SetParent(parent)
	o.filled = filled
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticZBarNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := "FilledZBar"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineZBar"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticZBarNode) populate() {
	var vertices []float32

	vertices = []float32{
		-0.1, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.5, -0.4, 0.0,
		0.1, -0.4, 0.0,
		0.1, 0.5, 0.0,
		-0.5, 0.5, 0.0,
		-0.5, 0.4, 0.0,
		-0.1, 0.4, 0.0,
	}

	r.shape.SetVertices(vertices)

	var indices []uint32

	if r.filled {
		indices = []uint32{
			0, 1, 2,
			2, 3, 0,
			0, 3, 4,
			0, 4, 7,
			6, 7, 4,
			6, 4, 5,
		}
	} else {
		indices = []uint32{
			0, 1, 2, 3, 4, 5, 6, 7,
		}
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// SetColor sets line color
func (r *StaticZBarNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticZBarNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticZBarNode) PointInside(p api.IPoint) bool {
	return false
}

// Draw renders shape
func (r *StaticZBarNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
