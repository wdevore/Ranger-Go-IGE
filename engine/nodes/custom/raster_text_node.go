package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// RasterTextNode is a basic raster pixel text node
type RasterTextNode struct {
	nodes.Node

	color []float32

	text                 string
	pixelSize            float32
	bottomJustified      bool
	bottomVerticalOffset float32
	inverted             bool
	whiteSpaceDistance   float32

	shape api.IAtlasShape

	m4 api.IMatrix4
}

// NewRasterTextNode constructs a rectangle shaped node
func NewRasterTextNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(RasterTextNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (r *RasterTextNode) Build(world api.IWorld) error {
	r.Node.Build(world)
	r.m4 = maths.NewMatrix4()

	r.color = color.NewPaletteInt64(color.White).Array()
	r.pixelSize = 5.0
	r.bottomJustified = true
	r.bottomVerticalOffset = 6.0
	r.inverted = false
	r.whiteSpaceDistance = 2.0

	r.shape = world.Atlas().Shape("Pixel")

	return nil
}

// SetColor sets line color
func (r *RasterTextNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetText sets the text displayed
func (r *RasterTextNode) SetText(text string) {
	r.text = text
}

// SetInverted changes the style of rendering
func (r *RasterTextNode) SetInverted(inverted bool) {
	r.inverted = inverted
}

// SetPixelSize sets the size of the pixel not the spacing.
// Use the node's scale property to "push" the pixels apart.
func (r *RasterTextNode) SetPixelSize(size float32) {
	r.pixelSize = size
}

// SetBottomJustified forces the font's bottom to be at baseline.
// Otherwise the font's top is at the baseline.
func (r *RasterTextNode) SetBottomJustified(justified bool) {
	r.bottomJustified = justified
}

// SetVerticalOffset sets the offset from the baseline if bottom justified.
func (r *RasterTextNode) SetVerticalOffset(offset float32) {
	r.bottomVerticalOffset = offset
}

var pixelShifts = []int{0, 1, 2, 3, 4, 5, 6, 7}

// Draw renders shape
func (r *RasterTextNode) Draw(model api.IMatrix4) {
	// if r.IsDirty() {
	// 	r.SetDirty(false)
	// }
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)

	renG.SetColor(r.color)

	// -------------------------------------------
	// Draw text
	// -------------------------------------------
	rasterFont := r.World().RasterFont()
	cx := float32(0.0)
	s := r.pixelSize
	rowWidth := float32(rasterFont.GlyphWidth())

	// Is the text colored or the space around it (aka inverted)
	bitInvert := uint8(1)
	if r.inverted {
		bitInvert = 0
	}

	gl.PointSize(s)

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
			for _, shift := range pixelShifts {
				bit := (g >> shift) & 1
				if bit == bitInvert {
					r.m4.Set(model) // Reset for this pixel
					// r.m4.ScaleByComp(1.0, -1.0, 1.0) // This is slower. Use "gy -= s" as below.
					r.m4.TranslateBy2Comps(gx, gy)
					renG.Render(r.shape, r.m4)
				}
				gx += s
			}
			gy -= s // move to next pixel-row in char
		}
		cx += rowWidth * s // move to next column/char/glyph
	}

	gl.PointSize(1)
}
