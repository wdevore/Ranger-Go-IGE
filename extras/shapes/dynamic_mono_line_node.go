package shapes

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// DynamicMonoLineNode is a basic dynamic Line
type DynamicMonoLineNode struct {
	nodes.Node

	shapeID int

	color []float32
}

// NewDynamicMonoLineNode creates a basic dynamic Line.
// It comes with default colors, and will Add a shape to the DynamicMono atlas
// Atlas IF its not present.
func NewDynamicMonoLineNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicMonoLineNode)

	o.Initialize(name)
	o.SetParent(parent)

	o.shapeID = -1

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *DynamicMonoLineNode) build(world api.IWorld) error {
	b.Node.Build(world)

	atl := world.GetAtlas(api.DynamicMonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find DynamicMono Atlas")
	}

	b.SetAtlas(atl)

	name := api.LineShapeName + b.Name()

	atlas := atl.(api.IDynamicAtlasX)

	b.shapeID = atlas.GetShapeByName(name)
	if b.shapeID < 0 {
		// Add shape
		vertices, indices, mode := generators.GenerateUnitHLineVectorShape()
		b.shapeID = atlas.AddShape(name, vertices, indices, mode)
	}

	// atlas.SetData(vertices, indices)
	// atlas.SetPrimitiveMode(gl.LINES)
	// atlas.SetIndicesCount(2)

	// Default colors
	b.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// SetVertex1 sets one of the points on the line
func (b *DynamicMonoLineNode) SetVertex1(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 0, b.shapeID)
}

// SetVertex2 sets one of the points on the line
func (b *DynamicMonoLineNode) SetVertex2(x, y float32) {
	atlas := b.Atlas().(api.IDynamicAtlasX)
	atlas.SetShapeVertex(x, y, 1, b.shapeID)
}

// SetColor sets the color
func (b *DynamicMonoLineNode) SetColor(color api.IPalette) {
	b.color = color.Array()
}

// SetAlpha overwrites the alpha value 0->1
func (b *DynamicMonoLineNode) SetAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *DynamicMonoLineNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()
	// atlas.(api.IDynamicAtlasX).Update()

	if b.shapeID > -1 {
		atlas.SetColor(b.color)
		atlas.Render(b.shapeID, model)
	}
}
