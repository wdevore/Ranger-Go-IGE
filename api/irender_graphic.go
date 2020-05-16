package api

// IRenderGraphic represents a graphic render state
type IRenderGraphic interface {
	ShaderInUse() bool
	BufferObjInUse() bool

	Use()
	UnUse()

	UseShader()
	UnUseShader()

	UnUseBufferObj()
	UseBufferObj()

	BufferObj() IBufferObject

	Atlas() IAtlas

	Render(shape IAtlasShape, model IMatrix4)
	SetColor(color []float32)

	Program() uint32
}
