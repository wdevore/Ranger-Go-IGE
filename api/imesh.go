package api

const (
	// MeshStatic represents static VBO buffers
	MeshStatic = 0
	// MeshDynamic represents dynamic single mesh buffers,
	// for example, PixelBuffer
	MeshDynamic = 1
	// MeshDynamicMulti represent dynamic multi mesh buffers,
	// for example, lines
	MeshDynamicMulti = 2
)

// IMesh represents 2D polygon
type IMesh interface {
	Vertices() []float32
	VerticesUsing(backingArrayIdx int) []float32

	GenNextBackingIndex() int
	AddArray() int
	ActivateArray(backingIdx int)

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
