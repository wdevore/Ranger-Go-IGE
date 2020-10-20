package api

// IDynamicPixelAtlasX is a container for pixels
type IDynamicPixelAtlasX interface {
	Update()
	SetVertex(x, y float32, index int)
	SetData(vertices []float32, indices []uint32)
	SetPixelActiveCount(count int)
}
