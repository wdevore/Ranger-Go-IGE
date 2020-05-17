package api

const (
	// StaticRenderGraphic generic static render setup
	StaticRenderGraphic = 0
	// DynamicRenderGraphic generic dynamic render setup
	DynamicRenderGraphic = 1
)

// IRender is the visual interface for drawing
type IRender interface {
	Draw(IMatrix4)
}
