package api

// IAtlas represents a collection of shapes
type IAtlas interface {
	GenerateShape(atlasName string, primitiveMode uint32, bufferType bool) IAtlasShape
	HasShapes() bool
	Shapes() []IAtlasShape
}
