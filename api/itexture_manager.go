package api

import "image"

// ITextureManager represents a collection of textures
type ITextureManager interface {
	LoadTexture(image string) (int, error)
	AccessTexture(index int) (*image.NRGBA, error)
}
