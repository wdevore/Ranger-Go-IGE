package api

// IMesh represents 2D polygon
type IMesh interface {
	Vertices() []float32
	Indices() []uint32

	AddVertex(x, y, z float32)

	SetVertex(x, y, z float32, index int)

	AddIndex(index int)

	GenBuffers()
	Bind()
}
