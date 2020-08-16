package api

import "image"

// IBufferObject represents a vector object
type IBufferObject interface {
	Construct(meshType int, atlas IAtlas)
	ConstructWithImage(image *image.NRGBA, textureIndex uint32, smooth bool, atlas IAtlas)

	Use()
	UnUse()

	Update(shape IAtlasShape)
}
