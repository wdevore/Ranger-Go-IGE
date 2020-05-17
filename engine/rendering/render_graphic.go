package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// RenderGraphic a graphic state for rendering against
type RenderGraphic struct {
	shader    api.IShader
	programID uint32

	modelLoc int32
	colorLoc int32

	bufObj api.IBufferObject

	shaderInUse    bool
	bufferObjInUse bool
}

// NewRenderGraphic creates a new graphic
func NewRenderGraphic(vertexShaderCode, fragmentShaderCode string, isStatic bool, populator api.FunctorAtlasPopulator) api.IRenderGraphic {
	o := new(RenderGraphic)

	o.bufObj = NewBufferObject()
	// pass functor for populating
	o.bufObj.Construct(isStatic, populator)

	// ---------------------------------------
	// Compile shader
	// ---------------------------------------
	o.shader = NewShaderFromCode(vertexShaderCode, fragmentShaderCode)

	err := o.shader.Compile()

	if err != nil {
		fmt.Println("RenderGraphic error: ")
		panic(err)
	}

	// Activate shader so we can query it.
	o.shader.Use()

	o.programID = o.shader.Program()

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

// ShaderInUse indicates if the graphic's shader is activated
func (r *RenderGraphic) ShaderInUse() bool {
	return r.shaderInUse
}

// BufferObjInUse indicates if this graphic's buffer is activated
func (r *RenderGraphic) BufferObjInUse() bool {
	return r.bufferObjInUse
}

// Use activates this graphic
func (r *RenderGraphic) Use() {
	if !r.shaderInUse {
		r.shader.Use()
		r.shaderInUse = true
	}
	if !r.bufferObjInUse {
		r.bufObj.Use()
		r.bufferObjInUse = true
	}
}

// UnUse deactivates this graphic
func (r *RenderGraphic) UnUse() {
	r.shaderInUse = false
	r.bufObj.UnUse()
	r.bufferObjInUse = false
}

// UseShader activates this graphic's shader
func (r *RenderGraphic) UseShader() {
	if !r.shaderInUse {
		r.shader.Use()
		r.shaderInUse = true
	}
}

// UnUseShader deactivates current shader
func (r *RenderGraphic) UnUseShader() {
	r.shaderInUse = false
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

// Atlas returns internal uniform atlas
func (r *RenderGraphic) Atlas() api.IAtlas {
	return r.bufObj.Atlas()
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
