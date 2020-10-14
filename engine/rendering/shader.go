// Package rendering defines features of shaders.
package rendering

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Shader represents a shader program
type Shader struct {
	vertexCode   string
	fragmentCode string

	vertexSrc   string
	fragmentSrc string

	program uint32 // GLuint
}

// NewShaderFromCode creates a blank shader -
// you must call Compile before that shader is valid.
func NewShaderFromCode(vertexCode, fragmentCode string) api.IShader {
	s := new(Shader)
	s.vertexCode = vertexCode
	s.fragmentCode = fragmentCode
	return s
}

// NewShader creates a blank shader. You must call Load before that shader is valid.
func NewShader(vertexSrc, fragmentSrc string) api.IShader {
	s := new(Shader)
	s.vertexSrc = vertexSrc
	s.fragmentSrc = fragmentSrc
	return s
}

// Load reads and compiles shader programs
func (s *Shader) Load(relativePath string) error {

	var err error
	s.vertexCode, s.fragmentCode, err = fetch(relativePath, s.vertexSrc, s.fragmentSrc)
	if err != nil {
		return err
	}

	return s.Compile()
}

// Compile compiles shader programs
func (s *Shader) Compile() error {
	// fmt.Println("Shader: compiling...")
	var err error

	s.program, err = newProgram(s.vertexCode, s.fragmentCode)
	if err != nil {
		return err
	}

	// fmt.Println("Shader compiling complete for ID:", s.program)
	return nil
}

// Use activates program
func (s *Shader) Use() {
	gl.UseProgram(s.program)
}

// Program returns the shader's program id
func (s *Shader) Program() uint32 {
	return s.program
}

func fetch(relativePath string, vertexSrc, fragmentSrc string) (vCode, fCode string, err error) {

	// Vertex source -----------------------------------------------
	filePath := fmt.Sprintf(relativePath+api.RelativeShaderPath+"%s", vertexSrc)

	var bytes []byte
	bytes, err = ioutil.ReadFile(filePath)

	if err != nil {
		return "", "", err
	}
	bytes = append(bytes, 0)

	vCode = string(bytes)

	// Fragment source -----------------------------------------------
	filePath = fmt.Sprintf(relativePath+api.RelativeShaderPath+"%s", fragmentSrc)

	bytes, err = ioutil.ReadFile(filePath)

	if err != nil {
		return "", "", err
	}
	bytes = append(bytes, 0)

	fCode = string(bytes)

	return
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compile(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compile(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil

}

func compile(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
