package shapes

import (
	"errors"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// MonoPolygonNode is a basic static polygon
type MonoPolygonNode struct {
	nodes.Node

	shapeID    int
	halfLength float32

	color []float32

	vertices *[]float32
}

// NewMonoPolygonNode creates a basic static polygon.
// It comes with default colors, and will Add a shape to the MonoStatic
// Atlas IF its not present.
func NewMonoPolygonNode(name string, vertices *[]float32, indices *[]uint32, drawStyle int, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(MonoPolygonNode)

	o.Initialize(name)
	o.SetParent(parent)
	o.vertices = vertices

	o.shapeID = -1

	parent.AddChild(o)

	if err := o.build(vertices, indices, drawStyle, world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *MonoPolygonNode) build(vertices *[]float32, indices *[]uint32, drawStyle int, world api.IWorld) error {
	b.Node.Build(world)

	b.halfLength = 0.5

	atl := world.GetAtlas(api.MonoAtlasName)

	if atl == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atl)

	name := api.PolygonShapeName + "_" + b.Name()

	atlas := atl.(api.IStaticAtlasX)

	mode := gl.TRIANGLES
	if drawStyle == api.OUTLINED {
		mode = gl.LINE_LOOP
	}
	b.shapeID = atlas.GetShapeByName(name)
	if b.shapeID < 0 {
		// Add shape
		b.shapeID = atlas.AddShape(name, *vertices, *indices, mode)
	}

	// Default colors
	b.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// Vertices returns shape's vertices
func (b *MonoPolygonNode) Vertices() *[]float32 {
	return b.vertices
}

// SetColor sets the color
func (b *MonoPolygonNode) SetColor(color api.IPalette) {
	b.color = color.Array()
}

// SetAlpha overwrites the alpha value 0->1
func (b *MonoPolygonNode) SetAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *MonoPolygonNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	if b.shapeID > -1 {
		atlas.SetColor(b.color)
		atlas.Render(b.shapeID, model)
	}
}
