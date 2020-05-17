package rendering

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// BufferObject associates an Atlas with a VAO
type BufferObject struct {
	atlasObject api.IAtlasObject // VectorAtlas
	vao         *VAO
	atlas       api.IAtlas
}

// NewBufferObject creates a new vector object with an associated Mesh
func NewBufferObject() api.IBufferObject {
	vo := new(BufferObject)
	return vo
}

// Construct configures a buffer object
// The second arg is a functor that given an atlas-object
// populate the atlas
func (b *BufferObject) Construct(isStatic bool, populator api.FunctorAtlasPopulator) {
	b.vao = NewVAO()
	b.vao.BindStart()

	// The Atlas contains a Mesh and the Mesh contains
	// VBOs and EBOs
	b.atlasObject = NewUniformAtlas(isStatic)

	b.atlas = NewAtlas()
	// Populate atlas with default objects
	populator(b.atlas, b.atlasObject)

	mesh := b.atlasObject.Mesh()

	mesh.BindVBO()
	mesh.BindEBO()

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

// Atlas returns the internal atlas
func (b *BufferObject) Atlas() api.IAtlas {
	return b.atlas
}
