package api

// IAtlas represents a collection of shapes
type IAtlas interface {
	Initialize(vo IVectorObject)
	Shape(string) IVectorShape
	Shapes() map[string]IVectorShape
	AddShape(vs IVectorShape)
}
