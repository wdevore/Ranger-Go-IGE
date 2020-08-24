package textures

import (
	"image"
	_ "image/png" // For 'png' images

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type textureManager struct {
	// Images and sprite sheets (aka texture atlases)
	images []*image.NRGBA

	atlases []api.ITextureAtlas

	index int
}

// NewTextureManager creates a texture manager that contains
// TextureAtlas(s)
func NewTextureManager() api.ITextureManager {
	o := new(textureManager)
	o.atlases = []api.ITextureAtlas{}
	return o
}

func (t *textureManager) AddAtlas(name, relativePath, textureManifest string, flipped bool) {
	// Open manifest to get texture file name
	ta := NewTextureAtlas(name, textureManifest)
	ta.Load(relativePath, flipped)

	t.atlases = append(t.atlases, ta)
}

func (t *textureManager) GetSTCoords(atlas, index int) *[]float32 {
	return t.atlases[atlas].TextureSTCoordsByIndex(index)
}

func (t *textureManager) GetAtlasByIndex(index int) api.ITextureAtlas {
	return t.atlases[index]
}

func (t *textureManager) GetAtlasByName(name string) api.ITextureAtlas {
	for _, ta := range t.atlases {
		if ta.Name() == name {
			return ta
		}
	}

	return nil
}

func (t *textureManager) GetAtlasPairByName(name string) (index int, atlas api.ITextureAtlas) {
	for i, ta := range t.atlases {
		if ta.Name() == name {
			return i, ta
		}
	}

	return -1, nil
}
