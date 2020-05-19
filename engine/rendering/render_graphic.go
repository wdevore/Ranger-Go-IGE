package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// RenderGraphic a graphic state for rendering against
type RenderGraphic struct {
	modelLoc int32
	colorLoc int32

	bufObj api.IBufferObject

	bufferObjInUse bool
}

// NewRenderGraphic creates a new graphic
func NewRenderGraphic(isStatic bool, atlas api.IAtlas, shader api.IShader) api.IRenderGraphic {
	o := new(RenderGraphic)

	o.bufObj = NewBufferObject()
	o.bufObj.Construct(isStatic, atlas)

	programID := shader.Program()

	// ---------------------------------------
	// Query shader
	// ---------------------------------------
	o.modelLoc = gl.GetUniformLocation(programID, gl.Str("model\x00"))
	if o.modelLoc < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	o.colorLoc = gl.GetUniformLocation(programID, gl.Str("fragColor\x00"))
	if o.colorLoc < 0 {
		panic("World: couldn't find 'fragColor' uniform variable")
	}

	return o
}

// BufferObjInUse indicates if this graphic's buffer is activated
func (r *RenderGraphic) BufferObjInUse() bool {
	return r.bufferObjInUse
}

// Use activates this graphic
func (r *RenderGraphic) Use() {
	if !r.bufferObjInUse {
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// UnUse deactivates this graphic
func (r *RenderGraphic) UnUse() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UnUseBufferObj deactivates buffer
func (r *RenderGraphic) UnUseBufferObj() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UseBufferObj activates buffer
func (r *RenderGraphic) UseBufferObj() {
	if !r.bufferObjInUse {
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// BufferObj returns internal buffer object
func (r *RenderGraphic) BufferObj() api.IBufferObject {
	return r.bufObj
}

// SetColor sets the shader's color
func (r *RenderGraphic) SetColor(color []float32) {
	gl.Uniform3fv(r.colorLoc, 1, &color[0])
}

// Render renders a shape
func (r *RenderGraphic) Render(shape api.IAtlasShape, model api.IMatrix4) {
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])

	// b.vao.Render(vs)
	// The signature of glDrawElements was defined back before there were buffer objects;
	// originally you'd pass an actual pointer to data in a client-side vertex array.
	// When device-side buffers were introduced, this function was extended to support them
	// as well, by shoehorning a buffer offset into the address argument.
	// Because we are using VBOs we need to awkwardly cast the offset value into a
	// pointer to void.
	// If we weren't using VBOs then we would use client-side addresses: &_mesh->indices[offset]
	// indices := b.atlasObject.Mesh().Indices()
	gl.DrawElements(shape.PrimitiveMode(), int32(shape.Count()), uint32(gl.UNSIGNED_INT), gl.PtrOffset(shape.Offset()))

	// gl.DrawElementsBaseVertex(shape.PrimitiveMode(), int32(shape.Count()), uint32(gl.UNSIGNED_INT), gl.Ptr(&indices[0]), int32(shape.Offset()))

}

// Update modifies the VBO buffer
func (r *RenderGraphic) Update(offset, vertexCount int) {
	r.bufObj.Update(offset, vertexCount*api.XYZComponentCount)
}
