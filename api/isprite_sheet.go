package api

import (
	"image"
)

// ISpriteSheet represents a collection of sprites (aka sub textures)
type ISpriteSheet interface {
	Name() string
	Load(relativePath string, flipped bool)
	SheetImage() *image.NRGBA
	TextureXYCoords(name string) *[]int
	TextureSTCoords(name string) *[]float32
	TextureSTCoordsByIndex(idx int) *[]float32
	GetIndex(name string) int
}
