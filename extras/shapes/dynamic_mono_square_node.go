package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// DynamicMonoSquareNode is a basic dynamic square
type DynamicMonoSquareNode struct {
	nodes.Node

	shapeID int

	color []float32
}

// NewDynamicMonoSquareNode creates a basic dynamic square.
// It comes with default colors, and will Add a shape to the DynamicMono atlas
// Atlas IF its not present.
func NewDynamicMonoSquareNode(name string, centered bool, forFilling bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicMonoSquareNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.shapeID = -1

	parent.AddChild(o)

	if err := o.build(centered, forFilling, world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *DynamicMonoSquareNode) build(centered bool, forFilling bool, world api.IWorld) error {
	b.Node.Build(world)

	atl := world.GetAtlas(api.DynamicMonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find DynamicMono Atlas")
	}

	b.SetAtlas(atl)

	name := api.SquareShapeName + b.Name()

	atlas := atl.(api.IDynamicAtlasX)

	b.shapeID = atlas.GetShapeByName(name)
	if b.shapeID < 0 {
		// Add shape
		vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(centered, forFilling)
		b.shapeID = atlas.AddShape(name, vertices, indices, mode)
	}

	// Default colors
	b.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetLowerLeft --
func (b *DynamicMonoSquareNode) SetLowerLeft(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 0, b.shapeID)
}

// SetLowerRight --
func (b *DynamicMonoSquareNode) SetLowerRight(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 1, b.shapeID)
}

// SetUpperRight --
func (b *DynamicMonoSquareNode) SetUpperRight(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 2, b.shapeID)
}

// SetUpperLeft --
func (b *DynamicMonoSquareNode) SetUpperLeft(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 3, b.shapeID)
}

// SetColor sets the color
func (b *DynamicMonoSquareNode) SetColor(color api.IPalette) {
	b.color = color.Array()
}

// SetAlpha overwrites the alpha value 0->1
func (b *DynamicMonoSquareNode) SetAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *DynamicMonoSquareNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	if b.shapeID > -1 {
		atlas.SetColor(b.color)
		atlas.Render(b.shapeID, model)
	}
}
