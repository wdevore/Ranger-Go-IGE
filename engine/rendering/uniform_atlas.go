package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// UniformAtlas defines a uniform colored atlas
type UniformAtlas struct {
	AtlasObject
}

// NewUniformAtlas creates a new uniform atlas
func NewUniformAtlas(isStatic bool) api.IAtlasObject {
	ua := new(UniformAtlas)
	ua.Initialize(isStatic)
	return ua
}

// BindAndBufferVBO calls the atlas's mesh VBO bind
func (ua *UniformAtlas) BindAndBufferVBO() {
	ua.Mesh().BindVBO()
}

// BindAndBufferEBO calls the atlas's mesh EBO bind
func (ua *UniformAtlas) BindAndBufferEBO() {
	ua.Mesh().BindEBO()
}

// Add adds a vertex
func (ua *UniformAtlas) Add(x, y, z float32, index int) {
	ua.AddVertex(x, y, z)
	ua.AddIndex(index)
}

// Add2Component adds a vertex and auto generated index
func (ua *UniformAtlas) Add2Component(x, y float32) {
	ua.Add(x, y, 0.0, ua.idx)
}

// Add3Component adds a vertex and auto generated index
func (ua *UniformAtlas) Add3Component(x, y, z float32) {
	ua.Add(x, y, z, ua.idx)
}
