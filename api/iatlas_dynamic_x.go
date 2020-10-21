package api

// IDynamicAtlasX is a container for dynamic shapes
type IDynamicAtlasX interface {
	AddShape(shapeName string, vertices []float32, indices []uint32, mode int) int
	GetShapeByName(shapeName string) int
	Update()
	SetData(vertices []float32, indices []uint32)
	SetPrimitiveMode(mode int)
	SetIndicesCount(count int)
	SetOffset(offset int)
	SetShapeVertex(x, y float32, index, shapdID int)
	SetVertex(x, y float32, index int)
}
