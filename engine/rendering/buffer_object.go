package rendering

import (
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// BufferObject associates an Atlas with a VAO
type BufferObject struct {
	atlasObject api.IAtlasObject // VectorAtlas
	vao         *VAO
}

// NewBufferObject creates a new vector object with an associated Mesh
func NewBufferObject() api.IBufferObject {
	vo := new(BufferObject)
	return vo
}

// Construct configures a buffer object
// The second arg is a functor that given an atlas-object
// populates the atlas
// func (b *BufferObject) Construct(isStatic bool, populator api.FunctorAtlasPopulator) {
func (b *BufferObject) Construct(isStatic bool, atlas api.IAtlas) {
	b.vao = NewVAO()
	b.vao.BindStart()

	// The Atlas contains a Mesh and the Mesh contains
	// VBOs and EBOs
	b.atlasObject = NewUniformAtlas(isStatic)

	// Populate atlas with objects
	atlas.Populate(b.atlasObject)

	mesh := b.atlasObject.Mesh()

	mesh.BindVBO()
	mesh.BindEBO()

	b.vao.BindComplete()
}

func (b *BufferObject) Vertices() []float32 {
	return b.atlasObject.Mesh().Vertices()
}

// Use activates the VAO
func (b *BufferObject) Use() {
	b.vao.Use()
}

// UnUse deactivates the VAO
func (b *BufferObject) UnUse() {
	b.vao.UnUse()
}

// Update modifies the VBO buffer
func (b *BufferObject) Update(offset, size int) {
	b.atlasObject.Mesh().Update(offset, size)
}

func (b *BufferObject) UpdatePreScaled(offset, size int) {
	b.atlasObject.Mesh().UpdatePreScaled(offset, size)
}

func (b *BufferObject) UpdatePreScaledUsing(offset, size int, vertices []float32) {
	// if v.bufferUsage == gl.STATIC_DRAW {
	// 	panic("VBO is not configured as DYNAMIC_DRAW.")
	// }

	gl.BindBuffer(gl.ARRAY_BUFFER, b.atlasObject.Mesh().VboID())

	// The last parameter should be a separate buffer
	// The 'offset' and 'size' are parameters for the destination
	// buffer and in is bytes
	// In other words the source buffer is captured as a whole,
	// the destination buffer is piece-meal based on offset/size
	gl.BufferSubData(gl.ARRAY_BUFFER, offset, size, gl.Ptr(vertices))

	if errNum := gl.GetError(); errNum != gl.NO_ERROR {
		log.Fatal("(BufObj)GL Error: ", errNum)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
