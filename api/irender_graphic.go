package api

import "image"

// IRenderGraphic represents a graphic render state
type IRenderGraphic interface {
	BufferObjInUse() bool
	Construct(meshType int, atlas IAtlas)
	ConstructWithImage(image *image.NRGBA, smooth bool, atlas IAtlas)

	Use()
	UnUse()

	UnUseBufferObj()
	UseBufferObj()

	BufferObj() IBufferObject

	SetColor(color []float32)
	SetColor4(color []float32)

	Render(shape IAtlasShape, model IMatrix4)
	RenderElements(shape IAtlasShape, elementCount, elementOffset int, model IMatrix4)

	Update(shape IAtlasShape)
	UpdateTexture(coords *[]float32)
}
