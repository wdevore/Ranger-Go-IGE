package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// BitmapFont9x9Node is a font node
type BitmapFont9x9Node struct {
	nodes.Node

	textureRenderer api.ITextureRenderer
	textureAtlas    api.ITextureAtlas

	// font indices
	textureIndexes []int

	text  string
	color []float32
	model api.IMatrix4 // Char spacing
}

// NewBitmapFont9x9Node constructs a bitmap font node
func NewBitmapFont9x9Node(name string, textureAtlas api.ITextureAtlas, textureRenderer api.ITextureRenderer, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(BitmapFont9x9Node)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	o.textureRenderer = textureRenderer
	o.textureAtlas = textureAtlas

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
		idx := d.textureAtlas.GetIndex(string(c))
		d.textureIndexes[i] = idx
	}
}

// SetColor sets the mixin color. default is white = no effect.
func (d *BitmapFont9x9Node) SetColor(color []float32) {
	d.color = color
}

// Draw renders the font indices
func (d *BitmapFont9x9Node) Draw(model api.IMatrix4) {
	d.World().SetRenderGraphic(api.TextureRenderGraphic)

	d.textureRenderer.Use()
	d.textureRenderer.SetColor(d.color)

	d.model.Set(model)

	for _, i := range d.textureIndexes {
		if i < 0 {
			d.model.TranslateBy2Comps(0.65, 0.0) // Simulate " "
		} else {
			d.textureRenderer.SelectCoordsByIndex(i)

			d.model.TranslateBy2Comps(0.85, 0.0)
			d.textureRenderer.Draw(d.model)
		}
	}
}
