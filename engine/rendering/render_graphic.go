package rendering

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// RenderGraphic a graphic state for rendering against
type RenderGraphic struct {
	programID uint32

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

	o.programID = shader.Program()

	// ---------------------------------------
	// Query shader
	// ---------------------------------------
	o.modelLoc = gl.GetUniformLocation(o.programID, gl.Str("model\x00"))
	if o.modelLoc < 0 {
		panic("World: couldn't find 'model' uniform variable")
	}

	o.colorLoc = gl.GetUniformLocation(o.programID, gl.Str("fragColor\x00"))
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

// Program returns the internal shader program
func (r *RenderGraphic) Program() uint32 {
	return r.programID
}

// SetColor sets the shader's color
func (r *RenderGraphic) SetColor(color []float32) {
	gl.Uniform3fv(r.colorLoc, 1, &color[0])
}

// Render renders a shape
func (r *RenderGraphic) Render(shape api.IAtlasShape, model api.IMatrix4) {
	gl.UniformMatrix4fv(r.modelLoc, 1, false, &model.Matrix()[0])

	r.bufObj.Render(shape)
}
