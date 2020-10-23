package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoTriangleNode is a basic static Triangle
type MonoTriangleNode struct {
	nodes.Node

	vertices []float32

	halfSide float32

	filledShapeID   int
	outlinedShapeID int

	filledColor   []float32
	outlinedColor []float32
}

// NewMonoTriangleNode creates a basic static Triangle.
// It comes with default colors, and will Add two shapes to the MonoStatic
// Atlas IF they are not present.
// drawStyle = FILLED, OUTLINED, FILLOUTLINED
func NewMonoTriangleNode(name string, drawStyle int, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoTriangleNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.filledShapeID = -1
	o.outlinedShapeID = -1

	parent.AddChild(o)

	if err := o.build(drawStyle, world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoTriangleNode) build(drawStyle int, world api.IWorld) error {
	b.Node.Build(world)

	b.halfSide = 0.5

	atl := world.GetAtlas(api.MonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atl)
	atlas := atl.(api.IStaticAtlasX)

	var indices []uint32
	var mode int

	if drawStyle == api.FILLED {
		name := api.FilledTriangleShapeName

		b.filledShapeID = atlas.GetShapeByName(name)
		if b.filledShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(true)
			b.filledShapeID = atlas.AddShape(name, b.vertices, indices, mode)
		}
	} else if drawStyle == api.OUTLINED {
		name := api.OutlinedTriangleShapeName

		b.outlinedShapeID = atlas.GetShapeByName(name)
		if b.outlinedShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(false)
			b.outlinedShapeID = atlas.AddShape(name, b.vertices, indices, mode)
		}
	} else {
		nameF := api.FilledTriangleShapeName
		nameO := api.OutlinedTriangleShapeName

		b.filledShapeID = atlas.GetShapeByName(nameF)
		if b.filledShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(true)
			b.filledShapeID = atlas.AddShape(nameF, b.vertices, indices, mode)
		}
		b.outlinedShapeID = atlas.GetShapeByName(nameO)
		if b.outlinedShapeID < 0 {
			// Add shape
			b.vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(false)
			b.outlinedShapeID = atlas.AddShape(nameO, b.vertices, indices, mode)
		}
	}

	// Default colors
	b.filledColor = color.NewPaletteInt64(color.DarkGray).Array()
	b.outlinedColor = color.NewPaletteInt64(color.White).Array()

	return nil
}

// Vertices returns local-space vertices
func (b *MonoTriangleNode) Vertices() *[]float32 {
	return &b.vertices
}

// SideLength returns the scale length
func (b *MonoTriangleNode) SideLength() float32 {
	return b.halfSide * b.Scale() * 2
}

// HalfSide returns the scaled half side length.
func (b *MonoTriangleNode) HalfSide() float32 {
	return b.halfSide * b.Scale()
}

// SetFilledColor sets the fill color
func (b *MonoTriangleNode) SetFilledColor(color api.IPalette) {
	b.filledColor = color.Array()
}

// SetFilledAlpha overwrites the filled alpha value 0->1
func (b *MonoTriangleNode) SetFilledAlpha(alpha float32) {
	b.filledColor[3] = alpha
}

// SetOutlineColor sets the outline color
func (b *MonoTriangleNode) SetOutlineColor(color api.IPalette) {
	b.outlinedColor = color.Array()
}

// SetOutlineAlpha overwrites the outline alpha value 0->1
func (b *MonoTriangleNode) SetOutlineAlpha(alpha float32) {
	b.outlinedColor[3] = alpha
}

// Draw renders shape
func (b *MonoTriangleNode) Draw(model api.IMatrix4) {
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
