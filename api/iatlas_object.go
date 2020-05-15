package api

// IAtlasObject represents a helper to meshes
type IAtlasObject interface {
	Initialize(isStatic bool)
	AddVertex(x, y, z float32) int
	SetVertex(x, y, z float32, index int)
	AddIndex(index int)
	Begin() int
	End() int
	Mesh() IMesh
}
