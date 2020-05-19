package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(isStatic bool, atlas IAtlas)
	Use()
	UnUse()
	Update(offset, count int)
}
