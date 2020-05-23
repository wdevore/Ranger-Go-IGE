package api

// IMesh represents 2D polygon
type IMesh interface {
	Vertices() []float32
	Indices() []uint32

	AddVertex(x, y, z float32)

	SetVertex(x, y, z float32, index int)

	AddIndex(index int)

	VboID() uint32
	BindVBO()
	Discard()

	BindEBO()

	Update(offset, vertexCount int)
	UpdatePreScaled(offset, vertexCount int)
}
