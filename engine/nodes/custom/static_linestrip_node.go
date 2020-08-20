package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// StaticLineStripNode is a generic node
type StaticLineStripNode struct {
	nodes.Node

	color []float32

	vertices []float32
	indices  []uint32

	shape api.IAtlasShape
}

// NewStaticLineStripNode constructs a generic shape node
func NewStaticLineStripNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticLineStripNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticLineStripNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.shape = world.Atlas().GenerateShape(r.Name(), gl.LINE_STRIP)

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.

	return nil
}

// Populate loads and sets up the node's shape
func (r *StaticLineStripNode) Populate(vertices []float32, indices []uint32) {
	r.shape.SetVertices(vertices)

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Vertices returns the shape's vertices
func (r *StaticLineStripNode) Vertices() *[]float32 {
	return r.shape.Vertices()
}

// SetColor sets line color
func (r *StaticLineStripNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticLineStripNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticLineStripNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor(r.color)
	renG.Render(r.shape, model)
}
