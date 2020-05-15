package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(isStatic bool)
	UniformAtlas() IAtlasObject
	Use()
	UnUse()
	Bind()
	Render(vs IAtlasShape)
}
