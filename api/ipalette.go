package api

// IPalette represents colors
type IPalette interface {
	// Components return each component
	Components() (r, g, b, a float32)
	Array() []float32

	// R is red component
	R() float32
	// G is green component
	G() float32
	// B is blus component
	B() float32
	// A is alpha component
	A() float32

	// SetColor
	SetColor(r, g, b, a float32)

	// SetRed
	SetRed(r float32)
	// SetGreen
	SetGreen(r float32)
	// SetBlue
	SetBlue(r float32)
	// SetAlpha
	SetAlpha(r float32)
}
