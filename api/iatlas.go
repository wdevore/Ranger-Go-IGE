package api

// IAtlas represents a collection of shapes
type IAtlas interface {
	Build(vo IBufferObject)
	Shape(string) IAtlasShape
	Shapes() map[string]IAtlasShape
	AddShape(vs IAtlasShape)
}
