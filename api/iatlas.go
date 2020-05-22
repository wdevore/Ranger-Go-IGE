package api

// IAtlas represents a collection of shapes
type IAtlas interface {
	Populate(uAtlas IAtlasObject)
	Shape(string) IAtlasShape
	Shapes() map[string]IAtlasShape
	AddShape(vs IAtlasShape)

	// OpenGL Object type specific behaviours
	GetNextIndex(glType int) int
}
