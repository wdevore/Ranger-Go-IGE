package atlas

type shape struct {
	id        int
	shapeName string
	dirty     bool

	// For static Atlases there vertices are discarded after construction.
	vertices      []float32
	vertexOffset  int
	indices       []uint32
	indicesOffset int // Bytes
	indicesCount  int // EBO indices length

	primitiveMode uint32
}
