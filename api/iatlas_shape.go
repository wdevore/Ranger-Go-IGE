package api

// IAtlasShape represents a vector shape
type IAtlasShape interface {
	SetOffset(offset int)
	SetElementOffset(offset int)
	Offset() int
	SetBufferSize(size int)
	BufferSize() int

	SetIndices(indices []uint32)
	Indices() []uint32

	SetVertex2D(x, y float32, index int)

	SetVertices(vertices []float32)
	Vertices() *[]float32

	SetElementCount(count int)
	ElementCount() int

	Name() string
	SetName(string)

	PrimitiveMode() uint32
	SetPrimitiveMode(uint32)

	Count() int
	SetCount(int)
}
