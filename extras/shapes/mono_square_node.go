package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoSquareNode is a basic static square
type MonoSquareNode struct {
	nodes.Node

	filledShapeID   int
	outlinedShapeID int

	filledColor   []float32
	outlinedColor []float32
}

// NewMonoSquareNode creates a basic static square.
// It comes with default colors, and will Add two shapes to the MonoStatic
// Atlas IF they are not present.
// drawStyle = FILLED, OUTLINED, FILLOUTLINED
func NewMonoSquareNode(name string, drawStyle int, centered bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoSquareNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.filledShapeID = -1
	o.outlinedShapeID = -1

	parent.AddChild(o)

	if err := o.build(drawStyle, centered, world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoSquareNode) build(drawStyle int, centered bool, world api.IWorld) error {
	b.Node.Build(world)

	atlas := world.GetAtlas(api.MonoAtlasName)

	if atlas == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atlas)

	if drawStyle == api.FILLED {
		name := ""
		if centered {
			name = api.CenteredFilledSquareShapeName
		} else {
			name = api.UnCenteredFilledSquareShapeName
		}

		b.filledShapeID = atlas.GetShapeByName(name)
		if b.filledShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(centered, true)
			b.filledShapeID = atlas.AddShape(name, vertices, indices, mode)
		}
	} else if drawStyle == api.OUTLINED {
		name := ""
		if centered {
			name = api.CenteredOutlinedSquareShapeName
		} else {
			name = api.UnCenteredOutlinedSquareShapeName
		}

		b.outlinedShapeID = atlas.GetShapeByName(name)
		if b.outlinedShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(centered, false)
			b.outlinedShapeID = atlas.AddShape(name, vertices, indices, mode)
		}
	} else {
		nameF := ""
		nameO := ""
		if centered {
			nameF = api.CenteredFilledSquareShapeName
			nameO = api.CenteredOutlinedSquareShapeName
		} else {
			nameF = api.UnCenteredFilledSquareShapeName
			nameO = api.UnCenteredOutlinedSquareShapeName
		}

		b.filledShapeID = atlas.GetShapeByName(nameF)
		if b.filledShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(centered, true)
			b.filledShapeID = atlas.AddShape(nameF, vertices, indices, mode)
		}
		b.outlinedShapeID = atlas.GetShapeByName(nameO)
		if b.outlinedShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(centered, false)
			b.outlinedShapeID = atlas.AddShape(nameO, vertices, indices, mode)
		}
	}

	// Default colors
	b.filledColor = color.NewPaletteInt64(color.DarkGray).Array()
	b.outlinedColor = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetFilledColor sets the fill color
func (b *MonoSquareNode) SetFilledColor(color api.IPalette) {
	b.filledColor = color.Array()
}

// SetFilledAlpha overwrites the filled alpha value 0->1
func (b *MonoSquareNode) SetFilledAlpha(alpha float32) {
	b.filledColor[3] = alpha
}

// SetOutlineColor sets the outline color
func (b *MonoSquareNode) SetOutlineColor(color api.IPalette) {
	b.outlinedColor = color.Array()
}

// SetOutlineAlpha overwrites the outline alpha value 0->1
func (b *MonoSquareNode) SetOutlineAlpha(alpha float32) {
	b.outlinedColor[3] = alpha
}

// Draw renders shape
func (b *MonoSquareNode) Draw(model api.IMatrix4) {
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
