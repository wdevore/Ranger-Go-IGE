package atlas

type shape struct {
	shapeName string

	// For static Atlases there vertices are discarded after construction.
	vertices []float32
	indices  []uint32

	indicesOffset int // Bytes
	indicesCount  int // EBO indices length

	primitiveMode uint32
}
