package api

// IAtlas represents a collection of shapes
type IAtlas interface {
	Initialize(vo IBufferObject)
	Shape(string) IAtlasShape
	Shapes() map[string]IAtlasShape
	AddShape(vs IAtlasShape)
}
