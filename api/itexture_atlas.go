package api

import (
	"image"
)

// ITextureAtlas represents a collection of textures
type ITextureAtlas interface {
	Name() string
	Build(relativePath string)
	AtlasImage() *image.NRGBA
	TextureXYCoords(name string) *[]int
	TextureSTCoords(name string) *[]float32
	TextureSTCoordsByIndex(idx int) *[]float32
}
