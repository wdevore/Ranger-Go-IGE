package extras

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicRectangleNode is a generic node
type DynamicRectangleNode struct {
	nodes.Node

	color []float32

	shape api.IAtlasShape
}

// NewDynamicRectangleNode constructs a generic shape node
func newDynamicRectangleNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicRectangleNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *DynamicRectangleNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.Red).Array()

	r.shape = world.DynoAtlas().GenerateShape(r.Name(), gl.LINE_LOOP)

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *DynamicRectangleNode) populate() {
	vertices := []float32{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1, 2, 3,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *DynamicRectangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *DynamicRectangleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// SetMinMax sets the lower-left and top-right corners of rectangle
func (r *DynamicRectangleNode) SetMinMax(minX, minY, maxX, maxY float32) {
	r.shape.SetVertex2D(minX, minY, 0)
	r.shape.SetVertex2D(maxX, minY, 1)
	r.shape.SetVertex2D(maxX, maxY, 2)
	r.shape.SetVertex2D(minX, maxY, 3)
	r.SetDirty(true)
}

// Draw renders shape
func (r *DynamicRectangleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicRenderGraphic)
	renG.SetColor(r.color)

	r.shape.SetCount(4)

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
