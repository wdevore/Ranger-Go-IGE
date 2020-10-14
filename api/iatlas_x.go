package api

// IAtlasX is a container for shapes and renders them.
type IAtlasX interface {
	// Attaches shaders
	Configure() error
	AddShape(shapeName string, vertices []float32, indices []uint32, mode int) int
	GetShapeByName(shapeName string) int

	Burn() error
	Shake()
	Bake() error

	Use()
	UnUse()
	SetColor(color []float32)
	Render(shapeID int, model IMatrix4)
}
