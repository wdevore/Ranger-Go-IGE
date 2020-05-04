package rendering

// VectorObject associates an Atlas with a VAO
type VectorObject struct {
	UniAtlas *VectorUniformAtlas
	vao      *VAO
}

// NewVectorObject creates a new vector object with an associated Mesh
func NewVectorObject() *VectorObject {
	vo := new(VectorObject)
	return vo
}

// Construct configures a vector object
func (vo *VectorObject) Construct() {
	vo.UniAtlas = NewVectorUniformAtlas(true)
	vo.vao = NewVAO(&vo.UniAtlas.mesh)
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
func (vo *VectorObject) Render(vs *VectorShape) {
	vo.vao.Render(vs)
}

func ConstructBasicShapes(vo *VectorObject) {

}
