package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// BitmapFont9x9Node is a font node
type BitmapFont9x9Node struct {
	nodes.Node

	shape api.IAtlasShape

	textureMan   api.ITextureManager
	textureAtlas api.ITextureAtlas

	verticesAndTexture []float32
	textureIndexes     []int

	text  string
	color []float32
	model api.IMatrix4
}

// NewBitmapFont9x9Node constructs a bitmap font node
func NewBitmapFont9x9Node(name, textureAtlas string, textureMan api.ITextureManager, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(BitmapFont9x9Node)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	o.textureMan = textureMan

	o.model = maths.NewMatrix4()
	o.textureAtlas = textureMan.GetAtlasByName(textureAtlas)
	o.color = color.NewPaletteInt64(color.Transparent).Array()

	return o, nil
}

// Build configures the node
func (d *BitmapFont9x9Node) Build(world api.IWorld) error {
	d.Node.Build(world)

	d.shape = world.ShapeAtlas().GenerateShape("DynamicTextureQuad", gl.TRIANGLES)

	d.textureIndexes = []int{}

	return nil
}

// Populate ...
func (d *BitmapFont9x9Node) Populate() {
	// These 2D vertices are interleaved with 2D texture coords
	// The s,t coords are sourced by the manifest based on index
	idx := d.textureIndexes[0]
	coords := d.textureMan.GetSTCoords(0, idx)
	c := *coords
	d.verticesAndTexture = []float32{
		// Pos           Tex
		//x  y       z/s    w/t
		-0.5, -0.5, c[0], c[1], // CCW
		0.5, -0.5, c[2], c[3],
		0.5, 0.5, c[4], c[5],
		-0.5, 0.5, c[6], c[7],
	}

	d.shape.SetVertices(d.verticesAndTexture)

	indices := []uint32{
		0, 1, 2,
		0, 2, 3,
	}

	d.shape.SetIndices(indices)

	d.shape.SetElementCount(len(indices))
}

// SetText updates the text rendered.
func (d *BitmapFont9x9Node) SetText(text string) {
	d.text = text
	d.textureIndexes = make([]int, len(text))

	for i, c := range text {
		sc := string(c)
		idx := d.textureAtlas.GetIndex(sc)
		d.textureIndexes[i] = idx
	}
}

// SetColor ...
func (d *BitmapFont9x9Node) SetColor(colr []float32) {
	d.color = colr
}

// Draw renders shape
func (d *BitmapFont9x9Node) Draw(model api.IMatrix4) {
	renG := d.World().UseRenderGraphic(api.TextureRenderGraphic)
	renG.SetColor4(d.color)

	d.model.Set(model)
	for _, i := range d.textureIndexes {
		if i < 0 {
			d.model.TranslateBy2Comps(0.65, 0.0) // Simulate " "
		} else {
			coords := d.textureMan.GetSTCoords(0, i)
			renG.UpdateTexture(coords)

			d.model.TranslateBy2Comps(0.85, 0.0)
			renG.Render(d.shape, d.model)
		}
	}
}
