package api

// IRenderGraphic represents a graphic render state
type IRenderGraphic interface {
	BufferObjInUse() bool

	Use()
	UnUse()

	UnUseBufferObj()
	UseBufferObj()

	BufferObj() IBufferObject

	SetColor(color []float32)
	Render(shape IAtlasShape, model IMatrix4)

	Update(offset, vertexCount int)
}
