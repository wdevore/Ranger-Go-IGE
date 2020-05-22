package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(isStatic bool, atlas IAtlas)

	Vertices() []float32

	Use()
	UnUse()

	Update(offset, size int)
	UpdatePreScaled(offset, size int)
	UpdatePreScaledUsing(offset, size int, vertices []float32)
}
