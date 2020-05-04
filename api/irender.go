package api

// IRender is the visual interface for drawing
type IRender interface {
	Draw(IMatrix4)
}
