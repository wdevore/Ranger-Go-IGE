package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicLineNode is a generic node
type DynamicLineNode struct {
	nodes.Node

	color []float32

	shape api.IAtlasShape
}

// NewDynamicLineNode constructs a generic shape node
func newDynamicLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicLineNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *DynamicLineNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.Red).Array()

	r.shape = world.DynoAtlas().GenerateShape(r.Name(), gl.LINES)

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *DynamicLineNode) populate() {
	vertices := []float32{
		0.0, 0.0, 0.0, // Begin point
		0.0, 0.0, 0.0, // End point
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *DynamicLineNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *DynamicLineNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// SetPoint1 sets point 1 end position
func (r *DynamicLineNode) SetPoint1(x, y float32) {
	r.shape.SetVertex2D(x, y, 0)
	r.SetDirty(true)
}

// SetPoint2 sets point 2 end position
func (r *DynamicLineNode) SetPoint2(x, y float32) {
	r.shape.SetVertex2D(x, y, 1)
	r.SetDirty(true)
}

// SetPoints sets both the begin and end points
func (r *DynamicLineNode) SetPoints(x1, y1, x2, y2 float32) {
	r.shape.SetVertex2D(x1, y1, 0)
	r.shape.SetVertex2D(x2, y2, 1)
	r.SetDirty(true)
}

// Draw renders shape
func (r *DynamicLineNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicRenderGraphic)
	renG.SetColor(r.color)

	r.shape.SetCount(2)

	// We need to override the element offset because this
	// buffer is shared between multiple line shapes, so
	// the offset is effectively 0 for each update. If the buffer
	// containes more than 1 line then an offset for the other
	// lines would be valid. But in this version all lines are
	// sharing.
	// TODO in the future, when we move to batch processing, things
	// will be different.
	r.shape.SetElementOffset(0)
	renG.Update(r.shape)

	renG.RenderElements(r.shape, r.shape.ElementCount(), 0, model)
}
