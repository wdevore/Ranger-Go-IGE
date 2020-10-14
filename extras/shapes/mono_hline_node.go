package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoHLineNode is a basic static HLine
type MonoHLineNode struct {
	nodes.Node

	shapeID int

	color []float32
}

// NewMonoHLineNode creates a basic static HLine.
// It comes with default colors, and will Add a shape to the MonoStatic
// Atlas IF its not present.
func NewMonoHLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoHLineNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.shapeID = -1

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoHLineNode) build(world api.IWorld) error {
	b.Node.Build(world)

	atlas := world.GetAtlas(api.MonoAtlasName)

	if atlas == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atlas)

	name := api.HLineShapeName

	b.shapeID = atlas.GetShapeByName(name)
	if b.shapeID < 0 {
		// Add shape
		vertices, indices, mode := generators.GenerateUnitHLineVectorShape()
		b.shapeID = atlas.AddShape(name, vertices, indices, mode)
	}

	// Default colors
	b.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetColor sets the color
func (b *MonoHLineNode) SetColor(color api.IPalette) {
	b.color = color.Array()
}

// SetAlpha overwrites the alpha value 0->1
func (b *MonoHLineNode) SetAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *MonoHLineNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	if b.shapeID > -1 {
		atlas.SetColor(b.color)
		atlas.Render(b.shapeID, model)
	}
}
