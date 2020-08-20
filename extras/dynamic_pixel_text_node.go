package extras

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicPixelTextNode is a generic node
type DynamicPixelTextNode struct {
	nodes.Node

	color []float32

	shape        api.IAtlasShape
	pixelCount   int
	pixelBufSize int

	text string

	pixelSize            float32
	bottomJustified      bool
	bottomVerticalOffset float32
	inverted             bool
	whiteSpaceDistance   float32

	pixelShifts []int
}

// NewDynamicTextNode constructs a generic shape node
func NewDynamicTextNode(name string, pixelBufSize int, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicPixelTextNode)
	o.Initialize(name)
	o.SetParent(parent)
	o.pixelBufSize = pixelBufSize
	if parent != nil {
		parent.AddChild(o)
	}

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *DynamicPixelTextNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.pixelShifts = []int{0, 1, 2, 3, 4, 5, 6, 7}

	r.color = color.NewPaletteInt64(color.White).Array()

	r.pixelSize = 5.0
	r.bottomJustified = true
	r.bottomVerticalOffset = 6.0
	r.inverted = false
	r.whiteSpaceDistance = 2.0

	r.shape = world.PixelAtlas().GenerateShape("PixelBuffer", gl.POINTS)

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *DynamicPixelTextNode) populate() {
	vertices := make([]float32, r.pixelBufSize*api.XYZComponentCount)

	r.shape.SetVertices(vertices)

	indices := []uint32{}
	for i := 0; i < r.pixelBufSize; i++ {
		indices = append(indices, uint32(i))
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Text returns current text
func (r *DynamicPixelTextNode) Text() string {
	return r.text
}

// SetText sets the text displayed
func (r *DynamicPixelTextNode) SetText(text string) {
	r.text = text
}

// SetInverted changes the style of rendering
func (r *DynamicPixelTextNode) SetInverted(inverted bool) {
	r.inverted = inverted
}

// SetPixelSize sets the size of the pixel not the spacing.
// Use the node's scale property to "push" the pixels apart.
func (r *DynamicPixelTextNode) SetPixelSize(size float32) {
	r.pixelSize = size
}

// SetBottomJustified forces the font's bottom to be at baseline.
// Otherwise the font's top is at the baseline.
func (r *DynamicPixelTextNode) SetBottomJustified(justified bool) {
	r.bottomJustified = justified
}

// SetVerticalOffset sets the offset from the baseline if bottom justified.
func (r *DynamicPixelTextNode) SetVerticalOffset(offset float32) {
	r.bottomVerticalOffset = offset
}

// SetColor sets color
func (r *DynamicPixelTextNode) SetColor(color []float32) {
	r.color = color
}

// Color gets the current color
func (r *DynamicPixelTextNode) Color() []float32 {
	return r.color
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *DynamicPixelTextNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

func (r *DynamicPixelTextNode) reGen() {
	// Update backing array
	// -------------------------------------------
	// Gen text
	// -------------------------------------------
	s := r.pixelSize

	rasterFont := r.World().RasterFont()
	cx := float32(0.0)
	rowWidth := float32(rasterFont.GlyphWidth())

	// Is the text colored or the space around it (aka inverted)
	bitInvert := uint8(1)
	if r.inverted {
		bitInvert = 0
	}

	r.pixelCount = 0

	for _, c := range r.text {
		if c == ' ' {
			cx += rowWidth * s / r.whiteSpaceDistance // move to next column/char/glyph
			continue
		}

		var gy float32
		// move y back to the "top/bottom" for each char
		if r.bottomJustified {
			gy = r.bottomVerticalOffset * s // bottom
		} else {
			gy = float32(0) // top
		}
		glyph := rasterFont.Glyph(byte(c))

		for _, g := range glyph {
			gx := cx // set to current column
			for _, shift := range r.pixelShifts {
				bit := (g >> shift) & 1
				if bit == bitInvert {
					r.shape.SetVertex2D(gx, gy, r.pixelCount)
					r.pixelCount++
				}
				gx += s
			}
			gy -= s // move to next pixel-row in char
		}
		cx += rowWidth * s // move to next column/char/glyph
	}
}

// Draw renders shape
func (r *DynamicPixelTextNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicPixelBufRenderGraphic)

	r.reGen()

	r.shape.SetCount(r.pixelCount)
	renG.Update(r.shape)

	renG.SetColor(r.color)

	gl.PointSize(r.pixelSize)

	renG.RenderElements(r.shape, r.pixelCount, 0, model)

	gl.PointSize(1)
}
