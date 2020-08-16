package api

// ITextureManager represents a collection of textures
type ITextureManager interface {
	AddAtlas(name, relativePath, textureManifest string, flipped bool)
	GetSTCoords(atlas, index int) *[]float32
	GetAtlasByIndex(index int) ITextureAtlas
	GetAtlasByName(name string) ITextureAtlas
}
