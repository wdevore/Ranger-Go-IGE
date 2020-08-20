package extras

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticPlusNode is a generic node
type StaticPlusNode struct {
	nodes.Node

	color []float32

	shape api.IAtlasShape
}

// NewStaticPlusNode constructs a generic shape node
func NewStaticPlusNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticPlusNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticPlusNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	r.shape = world.Atlas().GenerateShape("Plus", gl.LINES)

	// Populated shape
	r.populate()

	return nil
}

func (r *StaticPlusNode) populate() {
	vertices := []float32{
		-0.5, 0.0, 0.0,
		0.5, 0.0, 0.0,
		0.0, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1, 2, 3,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *StaticPlusNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticPlusNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticPlusNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor(r.color)
	renG.Render(r.shape, model)
}
