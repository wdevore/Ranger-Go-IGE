package rendering

// UniformAtlas defines a uniform colored atlas
type UniformAtlas struct {
	AtlasObject
}

// NewUniformAtlas creates a new uniform atlas
func NewUniformAtlas(isStatic bool) *UniformAtlas {
	vua := new(UniformAtlas)
	vua.Initialize(isStatic)
	return vua
}

// Add adds a vertex
func (vua *UniformAtlas) Add(x, y, z float32, index int) {
	vua.AddVertex(x, y, z)
	vua.AddIndex(index)
}

// Add2Component adds a vertex and auto generated index
func (vua *UniformAtlas) Add2Component(x, y float32) {
	vua.Add(x, y, 0.0, vua.idx)
}

// Add3Component adds a vertex and auto generated index
func (vua *UniformAtlas) Add3Component(x, y, z float32) {
	vua.Add(x, y, z, vua.idx)
}
