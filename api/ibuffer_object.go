package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(meshType int, atlas IAtlas)

	Use()
	UnUse()

	Update(shape IAtlasShape)
}
