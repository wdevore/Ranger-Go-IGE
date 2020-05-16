package api

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(isStatic bool)
	UniformAtlas() IAtlasObject
	Atlas() IAtlas
	Use()
	UnUse()
	Render(vs IAtlasShape)
}
