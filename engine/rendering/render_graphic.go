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
func NewRenderGraphic(meshType int, atlas api.IAtlas, shader api.IShader) api.IRenderGraphic {
	o := new(RenderGraphic)

	o.bufObj = NewBufferObject()
	o.bufObj.Construct(meshType, atlas)

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

// Vertices returns buffer object's mesh's backing array
func (r *RenderGraphic) Vertices() []float32 {
	return r.bufObj.Vertices()
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

// SetColor4 sets the shader's vec4 color
func (r *RenderGraphic) SetColor4(color []float32) {
	gl.Uniform4fv(r.colorLoc, 1, &color[0])
}

// Render renders a shape
//    The signature of glDrawElements was defined back before there were buffer objects;
//    originally you'd pass an actual pointer to data in a client-side vertex array.
//    When device-side buffers were introduced, this function was extended to support them
//    as well, by shoehorning a buffer offset into the address argument.
//    Because we are using VBOs we need to awkwardly cast the offset value into a
//    pointer to void.
//    If we weren't using VBOs then we would use client-side addresses: &_mesh->indices[offset]
//    indices := b.atlasObject.Mesh().Indices()
//    2nd parm = how many elements to draw
//    fmt.Println("count: ", int32(shape.Count()), ", ", shape.Offset())
//    Count = # of vertices to render
func (r *RenderGraphic) Render(shape api.IAtlasShape, model api.IMatrix4) {

	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(shape.Count()), uint32(gl.UNSIGNED_INT), gl.PtrOffset(shape.Offset()))

	// TODO evaluate for use
	// gl.DrawElementsBaseVertex(shape.PrimitiveMode(), int32(shape.Count()), uint32(gl.UNSIGNED_INT), gl.Ptr(&indices[0]), int32(shape.Offset()))
}

// RenderElements renders the specificied # of elemens from the shape's vertices
func (r *RenderGraphic) RenderElements(shape api.IAtlasShape, elementCount, elementOffset int, model api.IMatrix4) {
	// fmt.Println(elementCount, ", ", elementOffset)
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(elementCount), uint32(gl.UNSIGNED_INT), gl.PtrOffset(elementOffset))
}

// Update modifies the VBO buffer
func (r *RenderGraphic) Update(offset, count int) {
	r.bufObj.Update(offset, count*api.XYZComponentCount)
}

// UpdatePreScaled expects parameters already pre-scaled by data-type
func (r *RenderGraphic) UpdatePreScaled(offset, count int) {
	r.bufObj.UpdatePreScaled(offset, count)
}

// UpdatePreScaledUsing expects parameters already pre-scaled by data-type
// and a source buffer
func (r *RenderGraphic) UpdatePreScaledUsing(offset, size int, vertices []float32) {
	r.bufObj.UpdatePreScaledUsing(offset, size, vertices)
}
