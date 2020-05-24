package api

// IAtlasShape represents a vector shape
type IAtlasShape interface {
	SetOffset(offset int)
	SetElementOffset(offset int)
	Offset() int

	BackingArrayIdx() int
	SetBackingArrayIdx(idx int)

	Vertices(backingArrayIdx int) []float32

	InUse() bool
	SetInUse(inuse bool)

	SetElementCount(count int)
	ElementCount() int

	SetMaxSize(size int)
	MaxSize() int

	SetVertex3D(x, y, z float32, index int)
	SetVertex2D(x, y float32, index int)

	Name() string
	SetName(string)

	PrimitiveMode() uint32
	SetPrimitiveMode(uint32)

	Count() int
	SetCount(int)
}
