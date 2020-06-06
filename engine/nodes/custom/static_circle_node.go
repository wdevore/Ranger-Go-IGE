package custom

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticCircleNode is a generic node
type StaticCircleNode struct {
	nodes.Node

	color  []float32
	radius float64

	shape  api.IAtlasShape
	filled bool
}

// NewStaticCircleNode constructs a generic shape node
func NewStaticCircleNode(name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticCircleNode)
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
func (r *StaticCircleNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := "FilledCircle"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLE_FAN)
	} else {
		name := "OutlineCircle"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticCircleNode) populate() {
	var vertices []float32
	var indices []uint32

	segments := 12
	r.radius = 0.5 //1.0 // diameter of 1.0
	step := math.Pi / float64(segments)

	index := uint32(0)

	if r.filled {
		// Filled circles have a center point for the Fan fill algorithm
		vertices = append(vertices, 0.0, 0.0, 0.0)

		// Reference the center point
		indices = append(indices, 0)

		index++
	}

	for i := 0.0; i < 2.0*math.Pi; i += step {
		x := math.Cos(i) * r.radius
		y := math.Sin(i) * r.radius
		vertices = append(vertices, float32(x), float32(y), 0.0)
		indices = append(indices, index)
		index++
	}

	r.shape.SetVertices(vertices)

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Radius gets scaled radius, typically used by box2d.
func (r *StaticCircleNode) Radius() float32 {
	return float32(r.radius) * r.Scale()
}

// SetColor sets line color
func (r *StaticCircleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticCircleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticCircleNode) PointInside(p api.IPoint) bool {
	return false
}

// Draw renders shape
func (r *StaticCircleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
