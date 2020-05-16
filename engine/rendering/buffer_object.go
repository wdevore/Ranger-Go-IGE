package rendering

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// BufferObject associates an Atlas with a VAO
type BufferObject struct {
	uniformAtlas api.IAtlasObject // VectorAtlas
	vao          *VAO
	atlas        api.IAtlas
}

// NewBufferObject creates a new vector object with an associated Mesh
func NewBufferObject() api.IBufferObject {
	vo := new(BufferObject)
	return vo
}

// Construct configures a buffer object
func (b *BufferObject) Construct(isStatic bool) {
	b.vao = NewVAO()
	b.vao.BindStart()

	// The Atlas contains a Mesh and the Mesh contains
	// VBOs and EBOs
	b.uniformAtlas = NewUniformAtlas(isStatic)

	b.atlas = NewAtlas()
	// Populate atlas with default objects
	b.atlas.Build(b)

	b.uniformAtlas.BindAndBufferVBO()

	// Now EBO
	b.uniformAtlas.BindAndBufferEBO()

	b.vao.BindComplete()
}

// Use activates the VAO
func (b *BufferObject) Use() {
	b.vao.Use()
}

// UnUse deactivates the VAO
func (b *BufferObject) UnUse() {
	b.vao.UnUse()
}

// Render renders the given shape using the currently activated VAO
func (b *BufferObject) Render(vs api.IAtlasShape) {
	b.vao.Render(vs)
}

// UniformAtlas returns a uniform atlasobject
func (b *BufferObject) UniformAtlas() api.IAtlasObject {
	return b.uniformAtlas
}

// Atlas returns the internal atlas
func (b *BufferObject) Atlas() api.IAtlas {
	return b.atlas
}
