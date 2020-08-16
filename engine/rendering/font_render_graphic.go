package rendering

import (
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// FontRenderGraphic a graphic state for font rendering
type FontRenderGraphic struct {
	modelLoc int32
	colorLoc int32

	bufObj api.IBufferObject
	shader api.IShader

	bufferObjInUse bool
	atlas          api.IAtlas
}

// NewFontRenderGraphic creates a new graphic
func NewFontRenderGraphic(atlas api.IAtlas, shader api.IShader) api.IRenderGraphic {
	o := new(FontRenderGraphic)
	o.atlas = atlas

	o.bufObj = NewTextureQuadObject()
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

	// o.colorLoc = gl.GetUniformLocation(programID, gl.Str("fragColor\x00"))
	// if o.colorLoc < 0 {
	// 	panic("Couldn't find 'fragColor' uniform variable")
	// }

	return o
}

// Construct ...
func (r *FontRenderGraphic) Construct(meshType int, atlas api.IAtlas) {
}

// ConstructWithImage ...
func (r *FontRenderGraphic) ConstructWithImage(image *image.NRGBA, smooth bool, atlas api.IAtlas) {
	r.bufObj.ConstructWithImage(image, 0, smooth, atlas)
}

// BufferObjInUse indicates if this graphic's buffer is activated
func (r *FontRenderGraphic) BufferObjInUse() bool {
	return r.bufferObjInUse
}

// Use activates this graphic
func (r *FontRenderGraphic) Use() {
	if !r.bufferObjInUse {
		r.shader.Use()
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// UnUse deactivates this graphic
func (r *FontRenderGraphic) UnUse() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UnUseBufferObj deactivates buffer
func (r *FontRenderGraphic) UnUseBufferObj() {
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UseBufferObj activates buffer
func (r *FontRenderGraphic) UseBufferObj() {
	if !r.bufferObjInUse {
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// BufferObj returns internal buffer object
func (r *FontRenderGraphic) BufferObj() api.IBufferObject {
	return r.bufObj
}

// SetColor currently textures are mixed with vertex colors
func (r *FontRenderGraphic) SetColor(color []float32) {
}

// SetColor4 currently textures are mixed with vertex colors
func (r *FontRenderGraphic) SetColor4(color []float32) {
	// gl.Uniform4fv(r.colorLoc, 1, &color[0])
}

// Render renders a shape
func (r *FontRenderGraphic) Render(shape api.IAtlasShape, model api.IMatrix4) {
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(shape.ElementCount()), uint32(gl.UNSIGNED_INT), gl.PtrOffset(shape.Offset()))
}

// RenderElements renders the specificied # of elemens from the shape's vertices
func (r *FontRenderGraphic) RenderElements(shape api.IAtlasShape, elementCount, elementOffset int, model api.IMatrix4) {
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])
	gl.DrawElements(shape.PrimitiveMode(), int32(elementCount), uint32(gl.UNSIGNED_INT), gl.PtrOffset(elementOffset))
}

// Update modifies the VBO buffer
func (r *FontRenderGraphic) Update(shape api.IAtlasShape) {
	r.bufObj.Update(shape)
}

// UpdateTexture modifies the VBO buffer
func (r *FontRenderGraphic) UpdateTexture(coords *[]float32) {
	bo := r.bufObj.(*textureQuadObject)
	bo.TextureUpdate(coords)
}
