package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticHLineNode is a generic node
type StaticHLineNode struct {
	nodes.Node

	color      []float32
	halfLength float32

	shape api.IAtlasShape
}

// NewStaticHLineNode constructs a generic shape node
func NewStaticHLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticHLineNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticHLineNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	r.shape = world.Atlas().GenerateShape("HLine", gl.LINES)

	// Populated shape
	r.populate()

	return nil
}

func (r *StaticHLineNode) populate() {
	r.halfLength = 0.5 // Or 1.0

	vertices := []float32{
		-r.halfLength, 0.0, 0.0,
		r.halfLength, 0.0, 0.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// HalfLength returns the scaled half length.
func (r *StaticHLineNode) HalfLength() float32 {
	return r.halfLength * r.Scale()
}

// SetColor sets line color
func (r *StaticHLineNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticHLineNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticHLineNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
