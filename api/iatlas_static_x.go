package api

// IStaticAtlasX is a container for static shapes
type IStaticAtlasX interface {
	AddShape(shapeName string, vertices []float32, indices []uint32, mode int) int
	GetShapeByName(shapeName string) int
	FetchVerticesByName(shapeName string) *[]float32
}
