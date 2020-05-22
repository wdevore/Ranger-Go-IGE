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
	RenderElements(shape IAtlasShape, elementCount, elementOffset int, model IMatrix4)

	Vertices() []float32

	Update(offset, size int)
	UpdatePreScaled(offset, size int)
	UpdatePreScaledUsing(offset, size int, vertices []float32)
}
