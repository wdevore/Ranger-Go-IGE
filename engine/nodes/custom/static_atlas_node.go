package custom

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticAtlasNode is a generic node
type StaticAtlasNode struct {
	nodes.Node

	color []float32

	atlasName    string
	shape        api.IAtlasShape
	outlineShape api.IAtlasShape
}

// NewStaticAtlasNode constructs a generic shape node
func NewStaticAtlasNode(name, atlasName string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticAtlasNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.atlasName = atlasName
	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticAtlasNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.LightPurple).Array()

	r.shape = world.Atlas().Shape(r.atlasName)

	if r.shape == nil {
		return fmt.Errorf("StaticAtlasNode.Build: Shape '%s' not found", r.atlasName)
	}

	// FIXME i need access to outline shape

	if r.shape.IsOutlineType() {
		// Use the shape's data to setup polygon
		// 0 is used because static atlases have only 1 backing array
		v := r.shape.Vertices(0)

		// Subtract 1 because the offset always is 1 beyond
		// Subtract count*2 because each vertex is 6 floats 2(xyz) foreach
		// face of which there are 2.
		// o := (r.shape.Offset() - 1) - r.shape.Count()*2
		// o := r.shape.VertexIndex()

		polygon := r.shape.Polygon()
		max := r.shape.VertexIndex() + r.shape.Count()
		for i := r.shape.VertexIndex(); i < max; i += 3 {
			polygon.AddVertex(v[i], v[i+1])
		}
	}

	return nil
}

// SetColor sets line color
func (r *StaticAtlasNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticAtlasNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticAtlasNode) PointInside(p api.IPoint) bool {
	return r.shape.PointInside(p)
}

// Draw renders shape
func (r *StaticAtlasNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
