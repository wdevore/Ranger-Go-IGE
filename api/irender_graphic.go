package api

// IRenderGraphic represents a graphic render state
type IRenderGraphic interface {
	BufferObjInUse() bool

	Use()
	UnUse()

	UnUseBufferObj()
	UseBufferObj()

	BufferObj() IBufferObject

	Render(shape IAtlasShape, model IMatrix4)
	SetColor(color []float32)

	Program() uint32
}
