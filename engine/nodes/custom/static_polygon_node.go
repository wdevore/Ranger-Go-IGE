package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// StaticPolygonNode is a generic node
type StaticPolygonNode struct {
	nodes.Node

	color    []float32
	halfSide float32

	vertices []float32
	indices  []uint32

	shape    api.IAtlasShape
	centered bool
	filled   bool
}

// NewStaticPolygonNode constructs a generic shape node
func NewStaticPolygonNode(name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticPolygonNode)
	o.Initialize(name)
	o.SetParent(parent)
	o.filled = filled
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticPolygonNode) build(world api.IWorld) error {
	r.Node.Build(world)

	if r.filled {
		name := "FilledPolygon"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlinePolygon"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.

	return nil
}

// Populate loads and sets up the node's shape
func (r *StaticPolygonNode) Populate(vertices []float32, indices []uint32) {
	r.shape.SetVertices(vertices)

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Vertices returns the shape's vertices
func (r *StaticPolygonNode) Vertices() *[]float32 {
	return r.shape.Vertices()
}

// SetColor sets line color
func (r *StaticPolygonNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticPolygonNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticPolygonNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
