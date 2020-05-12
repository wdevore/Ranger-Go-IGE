package color

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type palette struct {
	r, g, b, a float32
}

// NewPalette constructs an IPalette black color object
func NewPalette() api.IPalette {
	o := new(palette)
	o.r = 0.0
	o.g = 0.0
	o.b = 0.0
	o.a = 1.0
	return o
}

// NewPaletteRGB constructs an RGB color object
func NewPaletteRGB(r, g, b uint8) api.IPalette {
	o := new(palette)
	o.r = float32(r) / 255.0
	o.g = float32(g) / 255.0
	o.b = float32(b) / 255.0
	o.a = 1.0
	return o
}

// NewPaletteRGBA constructs an RGB color object with alpha
func NewPaletteRGBA(r, g, b, a uint8) api.IPalette {
	o := NewPaletteRGB(r, g, b)
	o.SetAlpha(float32(a) / 255.0)
	return o
}

// NewPaletteInt64 constructs an RGB color object a single 64bit int
func NewPaletteInt64(c uint64) api.IPalette {
	o := NewPaletteRGBA(
		uint8((c&0xff000000)>>24),
		uint8((c&0x00ff0000)>>16),
		uint8((c&0x0000ff00)>>8),
		uint8(c&0x000000ff))
	return o
}

func (p *palette) Components() (r, g, b, a float32) {
	return r, g, b, a
}

func (p *palette) Array() []float32 {
	return []float32{p.r, p.g, p.b, p.a}
}

func (p *palette) R() float32 {
	return p.r
}

func (p *palette) G() float32 {
	return p.g
}

func (p *palette) B() float32 {
	return p.b
}

func (p *palette) A() float32 {
	return p.a
}

func (p *palette) SetColor(r, g, b, a float32) {
	p.r = r
	p.g = g
	p.b = b
	p.a = a
}

func (p *palette) SetRed(c float32) {
	p.r = c
}

func (p *palette) SetGreen(c float32) {
	p.g = c
}

func (p *palette) SetBlue(c float32) {
	p.b = c
}

func (p *palette) SetAlpha(c float32) {
	p.a = c
}

func (p palette) String() string {
	return fmt.Sprintf("{%3.0f,%3.0f,%3.0f,%3.0f}", p.r, p.g, p.b, p.a)
}
