package rendering

import (
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// TextureRenderGraphic a graphic state for texture rendering against
type TextureRenderGraphic struct {
	modelLoc int32
	colorLoc int32

	bufObj api.IBufferObject
	shader api.IShader

	bufferObjInUse bool
	atlas          api.IAtlas
}

// NewTextureRenderGraphic creates a new graphic
func NewTextureRenderGraphic(atlas api.IAtlas, shader api.IShader) api.IRenderGraphic {
	o := new(TextureRenderGraphic)
	o.atlas = atlas

	o.bufObj = NewTextureBufferObject()
	o.shader = shader

	programID := shader.Program()
	shader.Use()

	// ---------------------------------------
	// Query shader
	// ---------------------------------------
	o.modelLoc = gl.GetUniformLocation(programID, gl.Str("model\x00"))
	if o.modelLoc < 0 {
		panic("Couldn't find 'model' uniform variable")
	}

	return o
}

// Construct ...
func (r *TextureRenderGraphic) Construct(meshType int, atlas api.IAtlas) {
}

// ConstructWithImage ...
func (r *TextureRenderGraphic) ConstructWithImage(image *image.NRGBA, smooth bool, atlas api.IAtlas) {
	r.bufObj.ConstructWithImage(image, smooth, atlas)
}

// BufferObjInUse indicates if this graphic's buffer is activated
func (r *TextureRenderGraphic) BufferObjInUse() bool {
	return r.bufferObjInUse
}

// Use activates this graphic
func (r *TextureRenderGraphic) Use() {
	if !r.bufferObjInUse {
		r.shader.Use()
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// UnUse deactivates this graphic
func (r *TextureRenderGraphic) UnUse() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UnUseBufferObj deactivates buffer
func (r *TextureRenderGraphic) UnUseBufferObj() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UseBufferObj activates buffer
func (r *TextureRenderGraphic) UseBufferObj() {
	if !r.bufferObjInUse {
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// BufferObj returns internal buffer object
func (r *TextureRenderGraphic) BufferObj() api.IBufferObject {
	return r.bufObj
}

// SetColor sets the shader's color
func (r *TextureRenderGraphic) SetColor(color []float32) {
	// gl.Uniform3fv(r.colorLoc, 1, &color[0])
}

// SetColor4 sets the shader's vec4 color
func (r *TextureRenderGraphic) SetColor4(color []float32) {
	// gl.Uniform4fv(r.colorLoc, 1, &color[0])
}

// Render renders a shape
func (r *TextureRenderGraphic) Render(shape api.IAtlasShape, model api.IMatrix4) {
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(shape.ElementCount()), uint32(gl.UNSIGNED_INT), gl.PtrOffset(shape.Offset()))
}

// RenderElements renders the specificied # of elemens from the shape's vertices
func (r *TextureRenderGraphic) RenderElements(shape api.IAtlasShape, elementCount, elementOffset int, model api.IMatrix4) {
	// fmt.Println(elementCount, ", ", elementOffset)
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(elementCount), uint32(gl.UNSIGNED_INT), gl.PtrOffset(elementOffset))
}

// Update modifies the VBO buffer
func (r *TextureRenderGraphic) Update(shape api.IAtlasShape) {
	r.bufObj.Update(shape)
}
