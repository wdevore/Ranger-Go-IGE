package textures

import (
	"errors"
	"image"
	"image/draw"
	_ "image/png" // For 'png' images
	"os"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type textureManager struct {
	// Images and sprite sheets (aka texture atlases)
	images []*image.NRGBA

	index int
}

// NewTextureManager creates a texture manager
func NewTextureManager() api.ITextureManager {
	o := new(textureManager)

	return o
}

func (t *textureManager) LoadTexture(image string) (int, error) {
	rgb, err := loadImage(image)
	if err != nil {
		return 0, err
	}

	t.images = append(t.images, rgb)

	idx := t.index
	t.index++

	return idx, nil
}

func (t *textureManager) AccessTexture(index int) (image *image.NRGBA, err error) {
	if index > len(t.images) {
		return nil, errors.New("TextureManager: Index out of range")
	}

	image = t.images[index]

	if image == nil {
		return nil, errors.New("TextureManager: Image not in collection")
	}

	return image, nil
}

func loadImage(path string) (*image.NRGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	nrgba := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	// Transfer data to image
	draw.Draw(nrgba, nrgba.Bounds(), img, bounds.Min, draw.Src)

	return nrgba, nil
}
