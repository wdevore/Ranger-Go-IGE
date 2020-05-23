package api

// IAtlasObject represents a helper to meshes
type IAtlasObject interface {
	Initialize(meshType int)
	AddArray() int
	AddVertex(x, y, z float32) int
	SetVertex(x, y, z float32, index int)
	AddIndex(index int)
	Begin() int
	End() int
	Mesh() IMesh
	Update(offset, vertexCount int)
}
