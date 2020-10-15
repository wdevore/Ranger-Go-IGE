package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// MonoVLineNode is a basic static VLine
type MonoVLineNode struct {
	nodes.Node

	shapeID int

	color []float32
}

// NewMonoVLineNode creates a basic static VLine.
// It comes with default colors, and will Add a shape to the MonoStatic
// Atlas IF its not present.
func NewMonoVLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoVLineNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.shapeID = -1

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoVLineNode) build(world api.IWorld) error {
	b.Node.Build(world)

	atl := world.GetAtlas(api.MonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atl)

	name := api.VLineShapeName
	atlas := atl.(api.IStaticAtlasX)

	b.shapeID = atlas.GetShapeByName(name)
	if b.shapeID < 0 {
		// Add shape
		vertices, indices, mode := generators.GenerateUnitVLineVectorShape()
		b.shapeID = atlas.AddShape(name, vertices, indices, mode)
	}

	// Default colors
	b.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetColor sets the color
func (b *MonoVLineNode) SetColor(color api.IPalette) {
	b.color = color.Array()
}

// SetAlpha overwrites the alpha value 0->1
func (b *MonoVLineNode) SetAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *MonoVLineNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	if b.shapeID > -1 {
		atlas.SetColor(b.color)
		atlas.Render(b.shapeID, model)
	}
}
