package api

const (
	// GlobalRenderGraphic generic global render setup
	GlobalRenderGraphic = 0
)

// IRender is the visual interface for drawing
type IRender interface {
	Draw(IMatrix4)
}
