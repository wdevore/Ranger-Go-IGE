package api

// IAtlasShape represents a vector shape
type IAtlasShape interface {
	SetOffset(offset int)
	Offset() int

	SetVertex3D(x, y, z float32, index int)
	SetVertex2D(x, y float32, index int)

	Name() string
	SetName(string)

	PrimitiveMode() uint32
	SetPrimitiveMode(uint32)

	Count() int
	SetCount(int)
}
