package shapes

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

const (
	pixelBufferSize = 500
)

// DynamicPixelPixelTextNode is a generic node
type DynamicPixelPixelTextNode struct {
	nodes.Node

	color []float32

	shapeID int

	pixelBufSize int

	text string

	pixelSize            float32
	bottomJustified      bool
	bottomVerticalOffset float32
	inverted             bool
	whiteSpaceDistance   float32

	pixelShifts []int
}

// NewDynamicPixelTextNode constructs a generic shape node
func NewDynamicPixelTextNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicPixelPixelTextNode)
	o.Initialize(name)
	o.SetParent(parent)
	o.pixelBufSize = pixelBufferSize

	if parent != nil {
		parent.AddChild(o)
	}

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (m *DynamicPixelPixelTextNode) build(world api.IWorld) error {
	m.Node.Build(world)

	m.pixelShifts = []int{0, 1, 2, 3, 4, 5, 6, 7}

	m.color = color.NewPaletteInt64(color.White).Array()

	m.pixelSize = 5.0
	m.bottomJustified = true
	m.bottomVerticalOffset = 6.0
	m.inverted = false
	m.whiteSpaceDistance = 2.0

	dpAtlas := world.GetAtlas(api.DynamicPixelAtlasName)

	// This node requires the Dynamic Pixel Atlas. If not present then
	// create one.
	if dpAtlas == nil {
		dpAtlas = atlas.NewDynamicPixelAtlas(world)
		world.AddAtlas(api.DynamicPixelAtlasName, dpAtlas)

		vertices := make([]float32, m.pixelBufSize*api.XYZComponentCount)

		indices := []uint32{}
		for i := 0; i < m.pixelBufSize; i++ {
			indices = append(indices, uint32(i))
		}

		dpAtlas.(api.IDynamicPixelAtlasX).SetData(vertices, indices)

		err := dpAtlas.Burn()
		if err != nil {
			return err
		}
		fmt.Println("Dynamic Pixel Atlas created")
	}

	m.SetAtlas(dpAtlas)

	return nil
}

// Text returns current text
func (m *DynamicPixelPixelTextNode) Text() string {
	return m.text
}

// SetText sets the text displayed
func (m *DynamicPixelPixelTextNode) SetText(text string) {
	m.text = text
}

// SetInverted changes the style of rendering
func (m *DynamicPixelPixelTextNode) SetInverted(inverted bool) {
	m.inverted = inverted
}

// SetPixelSize sets the size of the pixel not the spacing.
// Use the node's scale property to "push" the pixels apart.
func (m *DynamicPixelPixelTextNode) SetPixelSize(size float32) {
	m.pixelSize = size
}

// SetBottomJustified forces the font's bottom to be at baseline.
// Otherwise the font's top is at the baseline.
func (m *DynamicPixelPixelTextNode) SetBottomJustified(justified bool) {
	m.bottomJustified = justified
}

// SetVerticalOffset sets the offset from the baseline if bottom justified.
func (m *DynamicPixelPixelTextNode) SetVerticalOffset(offset float32) {
	m.bottomVerticalOffset = offset
}

// SetColor sets color
func (m *DynamicPixelPixelTextNode) SetColor(color []float32) {
	m.color = color
}

// Color gets the current color
func (m *DynamicPixelPixelTextNode) Color() []float32 {
	return m.color
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (m *DynamicPixelPixelTextNode) SetAlpha(alpha float32) {
	m.color[3] = alpha
}

func (m *DynamicPixelPixelTextNode) reGen(atlas api.IDynamicPixelAtlasX) {
	// Update backing array
	// -------------------------------------------
	// Gen text
	// -------------------------------------------
	s := m.pixelSize

	rasterFont := m.World().RasterFont()
	cx := float32(0.0)
	rowWidth := float32(rasterFont.GlyphWidth())

	// Is the text colored or the space around it (aka inverted)
	bitInvert := uint8(1)
	if m.inverted {
		bitInvert = 0
	}

	pixelCount := 0

	for _, c := range m.text {
		if c == ' ' {
			cx += rowWidth * s / m.whiteSpaceDistance // move to next column/char/glyph
			continue
		}

		var gy float32
		// move y back to the "top/bottom" for each char
		if m.bottomJustified {
			gy = m.bottomVerticalOffset * s // bottom
		} else {
			gy = float32(0) // top
		}
		glyph := rasterFont.Glyph(byte(c))

		for _, g := range glyph {
			gx := cx // set to current column
			for _, shift := range m.pixelShifts {
				bit := (g >> shift) & 1
				if bit == bitInvert {
					atlas.SetVertex(gx, gy, pixelCount)
					pixelCount++

					if pixelCount > pixelBufferSize {
						panic("DynamicPixelText: text generated to many pixels for current buffer size.")
					}
				}
				gx += s
			}
			gy -= s // move to next pixel-row in char
		}
		cx += rowWidth * s // move to next column/char/glyph
	}

	atlas.SetPixelActiveCount(pixelCount)
}

// Draw renders shape
func (m *DynamicPixelPixelTextNode) Draw(model api.IMatrix4) {
	atlas := m.Atlas()
	dpAtlas := atlas.(api.IDynamicPixelAtlasX)

	m.reGen(dpAtlas)

	atlas.SetColor(m.color)

	dpAtlas.Update()

	gl.PointSize(m.pixelSize)

	atlas.Render(0, model)

	gl.PointSize(1)
}
