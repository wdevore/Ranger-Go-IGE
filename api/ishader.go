package api

// IShader represents a shader program
type IShader interface {
	Load(relativePath string) error
	Compile() error
	Use()
	Program() uint32
}
