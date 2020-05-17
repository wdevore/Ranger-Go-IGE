package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(isStatic bool, atlas IAtlas)
	Use()
	UnUse()
	Render(vs IAtlasShape)
	Update(offset, vertexCount int)
}
