package custom

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticArcNode is a generic node
type StaticArcNode struct {
	nodes.Node

	color      []float32
	radius     float64
	startAngle float64
	endAngle   float64
	segments   float64

	shape  api.IAtlasShape
	filled bool
}

// NewStaticArcNode constructs a generic shape node
func NewStaticArcNode(name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticArcNode)
	o.Initialize(name)
	o.SetParent(parent)

	o.filled = filled
	o.startAngle = 0.0
	o.endAngle = maths.DegreeToRadians * 45.0
	o.segments = 5
	o.radius = 0.5 // diameter of 1.0
	o.color = color.NewPaletteInt64(color.White).Array()

	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticArcNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	if r.filled {
		name := "FilledArc"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLE_FAN)
	} else {
		name := "OutlineArc"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticArcNode) populate() {
	var vertices []float32
	var indices []uint32

	step := (r.endAngle - r.startAngle) / float64(r.segments)

	index := uint32(0)

	vertices = append(vertices, 0.0, 0.0, 0.0)

	// Reference the center point
	indices = append(indices, 0)

	index++

	for i := r.startAngle; i <= r.startAngle+r.endAngle; i += step {
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

// Vertices returns the shape's vertices
func (r *StaticArcNode) Vertices() []float32 {
	return r.shape.Vertices()
}

// Radius gets scaled radius, typically used by box2d.
func (r *StaticArcNode) Radius() float32 {
	return float32(r.radius) * r.Scale()
}

// SetStartAngle sets the starting angle of the Arc
func (r *StaticArcNode) SetStartAngle(angle float64) {
	r.startAngle = angle
}

// SetEndAngle sets the starting angle of the Arc
func (r *StaticArcNode) SetEndAngle(angle float64) {
	r.endAngle = angle
}

// SetSegmentCount sets how many segments between start and end angle
func (r *StaticArcNode) SetSegmentCount(count float64) {
	r.segments = count
}

// SetColor sets line color
func (r *StaticArcNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticArcNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticArcNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor4(r.color)
	renG.Render(r.shape, model)
}
