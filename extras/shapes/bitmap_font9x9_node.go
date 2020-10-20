package shapes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// BitmapFont9x9Node is a font node that uses a SingleTextureAtlas and
// Font9x9SpriteSheet for rendering.
type BitmapFont9x9Node struct {
	nodes.Node

	// Sprite sheet containing font.
	spriteSheet api.ISpriteSheet

	// font indices
	textureIndexes []int

	text  string
	color []float32
	model api.IMatrix4 // Char spacing
}

// NewBitmapFont9x9Node constructs a node for rendering text.
func NewBitmapFont9x9Node(name string, singleTextureAtlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(BitmapFont9x9Node)

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
func (d *BitmapFont9x9Node) build(world api.IWorld) error {
	d.Node.Build(world)

	d.model = maths.NewMatrix4()
	d.color = color.NewPaletteInt64(color.White).Array()

	d.textureIndexes = []int{}

	return nil
}

// SetText updates the text rendered.
func (d *BitmapFont9x9Node) SetText(text string) {
	d.text = text
	d.textureIndexes = make([]int, len(text))

	for i, c := range text {
		idx := d.spriteSheet.GetIndex(string(c))
		d.textureIndexes[i] = idx
	}
}

// SetIndex sets the active sub texture index
func (d *BitmapFont9x9Node) SetIndex(index int) {
	d.textureIndexes = make([]int, 1)
	d.textureIndexes[0] = index
}

// SetColor sets the mixin color. default is white = no effect.
func (d *BitmapFont9x9Node) SetColor(color []float32) {
	d.color = color
}

// Draw renders the font indices
func (d *BitmapFont9x9Node) Draw(model api.IMatrix4) {
	atlas := d.Atlas()

	atlas.SetColor(d.color)

	d.model.Set(model)
	stAtlas := atlas.(api.ISingleTextureAtlasX)

	for _, i := range d.textureIndexes {
		if i < 0 {
			d.model.TranslateBy2Comps(0.65, 0.0) // Simulate " "
		} else {
			stAtlas.SelectCoordsByIndex(i)

			d.model.TranslateBy2Comps(0.85, 0.0)
			atlas.Render(0, d.model)
		}
	}
}
