package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// VectorObject associates an Atlas with a VAO
type VectorObject struct {
	uniformAtlas api.IVectorAtlas // VectorAtlas
	vao          *VAO
}

// NewVectorObject creates a new vector object with an associated Mesh
func NewVectorObject() api.IVectorObject {
	vo := new(VectorObject)
	return vo
}

// Construct configures a vector object
func (vo *VectorObject) Construct() {
	vo.uniformAtlas = NewVectorUniformAtlas(true)
	vo.vao = NewVAO(vo.uniformAtlas.Mesh())
}

// Use activates the VAO
func (vo *VectorObject) Use() {
	vo.vao.Use()
}

// UnUse deactivates the VAO
func (vo *VectorObject) UnUse() {
	vo.vao.UnUse()
}

// Bind binds the VAO
func (vo *VectorObject) Bind() {
	vo.vao.Bind()
}

// Render renders the given shape using the currently activated VAO
func (vo *VectorObject) Render(vs api.IVectorShape) {
	vo.vao.Render(vs)
}

// UniformAtlas returns atlas
func (vo *VectorObject) UniformAtlas() api.IVectorAtlas {
	return vo.uniformAtlas
}
