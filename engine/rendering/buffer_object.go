package rendering

import "github.com/wdevore/Ranger-Go-IGE/api"

// BufferObject associates an Atlas with a VAO
type BufferObject struct {
	uniformAtlas api.IAtlasObject // VectorAtlas
	vao          *VAO
}

// NewBufferObject creates a new vector object with an associated Mesh
func NewBufferObject() api.IBufferObject {
	vo := new(BufferObject)
	return vo
}

// Construct configures a vector object
func (bo *BufferObject) Construct(isStatic bool) {
	bo.uniformAtlas = NewUniformAtlas(isStatic)
	bo.vao = NewVAO(bo.uniformAtlas.Mesh())
}

// Use activates the VAO
func (bo *BufferObject) Use() {
	bo.vao.Use()
}

// UnUse deactivates the VAO
func (bo *BufferObject) UnUse() {
	bo.vao.UnUse()
}

// Bind binds the VAO
func (bo *BufferObject) Bind() {
	bo.vao.Bind()
}

// Render renders the given shape using the currently activated VAO
func (bo *BufferObject) Render(vs api.IAtlasShape) {
	bo.vao.Render(vs)
}

// UniformAtlas returns atlas
func (bo *BufferObject) UniformAtlas() api.IAtlasObject {
	return bo.uniformAtlas
}
