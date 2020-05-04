package rendering

// VectorUniformAtlas defines a uniform colored atlas
type VectorUniformAtlas struct {
	VectorAtlas
}

// NewVectorUniformAtlas creates a new uniform atlas
func NewVectorUniformAtlas(isStatic bool) *VectorUniformAtlas {
	vua := new(VectorUniformAtlas)
	vua.Initialize(isStatic, false)
	return vua
}

// Add adds a vertex and auto generated index
func (vua *VectorUniformAtlas) Add(x, y, z float32, index int) {
	vua.AddVertex(x, y, z)
	vua.AddIndex(index)
}

// Add2Component adds a vertex and auto generated index
func (vua *VectorUniformAtlas) Add2Component(x, y float32) {
	vua.Add(x, y, 0.0, vua.Idx)
}

// Add3Component adds a vertex and auto generated index
func (vua *VectorUniformAtlas) Add3Component(x, y, z float32) {
	vua.Add(x, y, z, vua.Idx)
}
