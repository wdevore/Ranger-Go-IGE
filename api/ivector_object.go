package api

// IVectorObject represents a vector object
type IVectorObject interface {
	Construct()
	UniformAtlas() IVectorAtlas
	Use()
	UnUse()
	Bind()
	Render(vs IVectorShape)
}
