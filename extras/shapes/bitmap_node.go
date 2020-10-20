package shapes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// BitmapNode is a basic node that uses a SingleTextureAtlas
type BitmapNode struct {
	nodes.Node

	// Sprite sheet containing font.
	spriteSheet api.ISpriteSheet

	// active index
	textureIndex int

	text  string
	color []float32
}

// NewBitmapNode constructs a node for rendering text.
func NewBitmapNode(name string, singleTextureAtlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(BitmapNode)

	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	o.SetAtlas(singleTextureAtlas)

	stAtlas := singleTextureAtlas.(api.ISingleTextureAtlasX)
	o.spriteSheet = stAtlas.SpriteSheet()

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (d *BitmapNode) build(world api.IWorld) error {
	d.Node.Build(world)

	d.color = color.NewPaletteInt64(color.White).Array()

	d.textureIndex = 0

	return nil
}

// SetIndex sets the active sub texture index
func (d *BitmapNode) SetIndex(index int) {
	d.textureIndex = index
}

// SetColor sets the mixin color. default is white = no effect.
func (d *BitmapNode) SetColor(color []float32) {
	d.color = color
}

// Draw renders the font indices
func (d *BitmapNode) Draw(model api.IMatrix4) {
	atlas := d.Atlas()

	atlas.SetColor(d.color)

	stAtlas := atlas.(api.ISingleTextureAtlasX)

	stAtlas.SelectCoordsByIndex(d.textureIndex)

	atlas.Render(0, model)
}
