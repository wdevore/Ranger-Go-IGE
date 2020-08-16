package api

import (
	"image"
)

// ITextureAtlas represents a collection of textures
type ITextureAtlas interface {
	Name() string
	Load(relativePath string, flipped bool)
	AtlasImage() *image.NRGBA
	TextureXYCoords(name string) *[]int
	TextureSTCoords(name string) *[]float32
	TextureSTCoordsByIndex(idx int) *[]float32
	SetLayer(float32)
	Layer() float32
	GetIndex(name string) int
}
