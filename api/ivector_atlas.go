package api

// IVectorAtlas represents a helper to meshes
type IVectorAtlas interface {
	Initialize(isStatic, hasColors bool)
	AddVertex(x, y, z float32) int
	AddIndex(index int)
	Begin() int
	End() int
	Mesh() IMesh
}
