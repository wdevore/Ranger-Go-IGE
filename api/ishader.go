package api

// IShader represents a shader program
type IShader interface {
	Load() error
	Compile() error
	Use()
	Program() uint32
}
