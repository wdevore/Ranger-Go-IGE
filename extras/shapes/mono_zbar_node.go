package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoZBarNode is a basic static ZBar
type MonoZBarNode struct {
	nodes.Node

	filledShapeID   int
	outlinedShapeID int

	filledColor   []float32
	outlinedColor []float32
}

// NewMonoZBarNode creates a basic static ZBar.
// It comes with default colors, and will Add two shapes to the MonoStatic
// Atlas IF they are not present.
// drawStyle = FILLED, OUTLINED, FILLOUTLINED
func NewMonoZBarNode(name string, drawStyle int, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoZBarNode)

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

func (b *MonoZBarNode) build(drawStyle int, world api.IWorld) error {
	b.Node.Build(world)

	atlas := world.GetAtlas(api.MonoAtlasName)

	if atlas == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atlas)

	if drawStyle == api.FILLED {
		name := api.FilledZBarShapeName

		b.filledShapeID = atlas.GetShapeByName(name)
		if b.filledShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitZBarVectorShape(true)
			b.filledShapeID = atlas.AddShape(name, vertices, indices, mode)
		}
	} else if drawStyle == api.OUTLINED {
		name := api.OutlinedZBarShapeName

		b.outlinedShapeID = atlas.GetShapeByName(name)
		if b.outlinedShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitZBarVectorShape(false)
			b.outlinedShapeID = atlas.AddShape(name, vertices, indices, mode)
		}
	} else {
		nameF := api.FilledZBarShapeName
		nameO := api.OutlinedZBarShapeName

		b.filledShapeID = atlas.GetShapeByName(nameF)
		if b.filledShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitZBarVectorShape(true)
			b.filledShapeID = atlas.AddShape(nameF, vertices, indices, mode)
		}
		b.outlinedShapeID = atlas.GetShapeByName(nameO)
		if b.outlinedShapeID < 0 {
			// Add shape
			vertices, indices, mode := generators.GenerateUnitZBarVectorShape(false)
			b.outlinedShapeID = atlas.AddShape(nameO, vertices, indices, mode)
		}
	}

	// Default colors
	b.filledColor = color.NewPaletteInt64(color.DarkGray).Array()
	b.outlinedColor = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetFilledColor sets the fill color
func (b *MonoZBarNode) SetFilledColor(color api.IPalette) {
	b.filledColor = color.Array()
}

// SetFilledAlpha overwrites the filled alpha value 0->1
func (b *MonoZBarNode) SetFilledAlpha(alpha float32) {
	b.filledColor[3] = alpha
}

// SetOutlineColor sets the outline color
func (b *MonoZBarNode) SetOutlineColor(color api.IPalette) {
	b.outlinedColor = color.Array()
}

// SetOutlineAlpha overwrites the outline alpha value 0->1
func (b *MonoZBarNode) SetOutlineAlpha(alpha float32) {
	b.outlinedColor[3] = alpha
}

// Draw renders shape
func (b *MonoZBarNode) Draw(model api.IMatrix4) {
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
