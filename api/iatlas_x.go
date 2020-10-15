package api

// IAtlasX is a container for shapes and renders them.
type IAtlasX interface {
	Configure() error

	Burn() error
	Shake()
	Bake() error

	Use()
	UnUse()

	SetColor(color []float32)
	Render(shapeID int, model IMatrix4)
}
