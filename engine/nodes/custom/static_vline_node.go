package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticVLineNode is a generic node
type StaticVLineNode struct {
	nodes.Node

	color []float32

	shape api.IAtlasShape
}

// NewStaticVLineNode constructs a generic shape node
func NewStaticVLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticVLineNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticVLineNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	r.shape = world.Atlas().GenerateShape("VLine", gl.LINES)

	r.populate()

	return nil
}

func (r *StaticVLineNode) populate() {
	vertices := []float32{
		0.0, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *StaticVLineNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticVLineNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticVLineNode) PointInside(p api.IPoint) bool {
	return false
}

// Draw renders shape
func (r *StaticVLineNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
