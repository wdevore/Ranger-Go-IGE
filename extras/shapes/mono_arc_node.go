package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoArcNode is a basic static Arc
type MonoArcNode struct {
	nodes.Node

	filledShapeID   int
	outlinedShapeID int

	filledColor   []float32
	outlinedColor []float32

	vertices []float32
}

// NewMonoArcNode creates a basic static Arc.
// It comes with default colors, and will Add two shapes to the MonoStatic
// Atlas IF they are not present.
// drawStyle = FILLED, OUTLINED, FILLOUTLINED
func NewMonoArcNode(name string, drawStyle, segments int, startAngle, endAngle float64, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoArcNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.filledShapeID = -1
	o.outlinedShapeID = -1

	parent.AddChild(o)

	if err := o.build(drawStyle, segments, startAngle, endAngle, world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoArcNode) build(drawStyle, segments int, startAngle, endAngle float64, world api.IWorld) error {
	b.Node.Build(world)

	atl := world.GetAtlas(api.MonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atl)

	atlas := atl.(api.IStaticAtlasX)
	var indices []uint32
	var mode int

	if drawStyle == api.FILLED {
		name := api.FilledArcShapeName

		b.filledShapeID = atlas.GetShapeByName(name)
		if b.filledShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitArcVectorShape(startAngle, endAngle, segments, true)
			b.filledShapeID = atlas.AddShape(name, b.vertices, indices, mode)
		}
	} else if drawStyle == api.OUTLINED {
		name := api.OutlinedArcShapeName

		b.outlinedShapeID = atlas.GetShapeByName(name)
		if b.outlinedShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitArcVectorShape(startAngle, endAngle, segments, false)
			b.outlinedShapeID = atlas.AddShape(name, b.vertices, indices, mode)
		}
	} else {
		nameF := api.FilledArcShapeName
		nameO := api.OutlinedArcShapeName

		b.filledShapeID = atlas.GetShapeByName(nameF)
		if b.filledShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitArcVectorShape(startAngle, endAngle, segments, true)
			b.filledShapeID = atlas.AddShape(nameF, b.vertices, indices, mode)
		}
		b.outlinedShapeID = atlas.GetShapeByName(nameO)
		if b.outlinedShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitArcVectorShape(startAngle, endAngle, segments, false)
			b.outlinedShapeID = atlas.AddShape(nameO, b.vertices, indices, mode)
		}
	}

	// Default colors
	b.filledColor = color.NewPaletteInt64(color.DarkGray).Array()
	b.outlinedColor = color.NewPaletteInt64(color.White).Array()

	return nil
}

// Vertices returns shape's vertices
func (b *MonoArcNode) Vertices() *[]float32 {
	return &b.vertices
}

// SetFilledColor sets the fill color
func (b *MonoArcNode) SetFilledColor(color api.IPalette) {
	b.filledColor = color.Array()
}

// SetFilledAlpha overwrites the filled alpha value 0->1
func (b *MonoArcNode) SetFilledAlpha(alpha float32) {
	b.filledColor[3] = alpha
}

// SetOutlineColor sets the outline color
func (b *MonoArcNode) SetOutlineColor(color api.IPalette) {
	b.outlinedColor = color.Array()
}

// SetOutlineAlpha overwrites the outline alpha value 0->1
func (b *MonoArcNode) SetOutlineAlpha(alpha float32) {
	b.outlinedColor[3] = alpha
}

// Draw renders shape
func (b *MonoArcNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	if b.filledShapeID > -1 {
		atlas.SetColor(b.filledColor)
		atlas.Render(b.filledShapeID, model)
	}

	if b.outlinedShapeID > -1 {
		atlas.SetColor(b.outlinedColor)
		atlas.Render(b.outlinedShapeID, model)
	}
}
