package api

// ITextureManager represents a collection of textures
type ITextureManager interface {
	AddAtlas(name, relativePath, textureManifest string)
	GetSTCoords(atlas, index int) *[]float32
	GetAtlasByIndex(index int) ITextureAtlas
	GetAtlasByName(name string) ITextureAtlas
	// LoadTexture(relativePath, image string, flipped bool) (int, error)
	// AccessTexture(index int) (*image.NRGBA, error)
}
