package extras

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticPointNode is a generic node
type StaticPointNode struct {
	nodes.Node

	color []float32
	size  float32

	shape api.IAtlasShape
}

// NewStaticPointNode constructs a generic shape node
func NewStaticPointNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticPointNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticPointNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()
	r.size = 1.0

	r.shape = world.Atlas().GenerateShape("Point", gl.POINTS)

	r.populate()

	return nil
}

func (r *StaticPointNode) populate() {
	vertices := []float32{
		0.0, 0.0, 0.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *StaticPointNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetSize sets the size of the point
func (r *StaticPointNode) SetSize(size float32) {
	r.size = size
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticPointNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticPointNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor(r.color)
	gl.PointSize(r.size)
	renG.Render(r.shape, model)
	gl.PointSize(1)
}
