package textures

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/assets/images"
)

// It requires a manifest doc that describes the sub-texture layouts
// For example:
//
// texture-atlas.png
// 256x256
// mine|0,64:64,64:64,128:0,128
// green ship|64,192:128,192:128,256:64,256
// orange ship|192,192:256,192:256,256:192,256
// ctype ship|0,192:64,192:64,256:0,256
// bomb|192,160:208,160:208,176:192,176
//

// SubTexture is a block within the image atlas.
//
//   TextureCoord  .------------.  TextureCoord
//                 |            |
//                 | SubTexture |           --> more sub textures...
//                 |            |
//                 |            |
//   TextureCoord  .------------.  TextureCoord
//
type textureAtlas struct {
	name          string
	manifest      string
	width, height int64

	atlas *image.NRGBA

	layer float32

	manifestJ images.TextureManifestJSON
}

// NewTextureAtlas creates a new atlas
func NewTextureAtlas(name, manifest string) api.ITextureAtlas {
	o := new(textureAtlas)
	o.name = name
	o.manifest = manifest
	return o
}

// Build setups the atlas based on manifest
func (t *textureAtlas) Load(relativePath string, flipped bool) {
	dataPath, err := filepath.Abs(relativePath)
	if err != nil {
		panic(err)
	}

	eConfFile, err := os.Open(dataPath + "/" + t.manifest)
	if err != nil {
		panic(err)
	}

	defer eConfFile.Close()

	bytes, err := ioutil.ReadAll(eConfFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, &t.manifestJ)
	if err != nil {
		panic(err)
	}

	file := dataPath + "/" + t.manifestJ.OutputPNG
	fmt.Println("TextureAtlas.Load loading: ", file)
	image, err := t.loadImage(file, flipped)

	if err != nil {
		panic(err)
	}

	t.atlas = image
	t.layer = t.manifestJ.Layer
}

// AtlasImage returns image atlas
func (t *textureAtlas) AtlasImage() *image.NRGBA {
	return t.atlas
}

// Name returns atlas name
func (t *textureAtlas) Name() string {
	return t.name
}

// SetLayer sets the 3D texture.z
func (t *textureAtlas) SetLayer(id float32) {
	t.layer = id
}

// Layer gets the 3D texture.z
func (t *textureAtlas) Layer() float32 {
	return t.layer
}

// GetIndex from name
func (t *textureAtlas) GetIndex(name string) int {
	for i, subTex := range t.manifestJ.Tiles {
		if name == subTex.Name {
			return i
		}
	}

	return -1
}

// TextureXYCoords returns the assigned coords of named sub texture to tile.
func (t *textureAtlas) TextureXYCoords(name string) *[]int {
	for _, subTex := range t.manifestJ.Tiles {
		if name == subTex.Name {
			return &subTex.XYCoords
		}
	}

	return nil
}

// TextureSTCoords returns the assigned coords of named sub texture to tile.
func (t *textureAtlas) TextureSTCoords(name string) *[]float32 {
	for _, subTex := range t.manifestJ.Tiles {
		if name == subTex.Name {
			return &subTex.STCoords
		}
	}

	return nil
}

func (t *textureAtlas) TextureSTCoordsByIndex(idx int) *[]float32 {
	return &t.manifestJ.Tiles[idx].STCoords
}

func (t *textureAtlas) loadImage(path string, flipped bool) (*image.NRGBA, error) {
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

	if flipped {
		r := image.Rect(0, 0, bounds.Dx(), bounds.Dy())
		flippedImg := image.NewNRGBA(r)

		// Flip horizontally or around Y-axis
		// for j := 0; j < nrgba.Bounds().Dy(); j++ {
		// 	for i := 0; i < nrgba.Bounds().Dx(); i++ {
		// 		flippedImg.Set(bounds.Dx()-i, j, nrgba.At(i, j))
		// 	}
		// }

		// Flip vertically or around the X-axis
		for j := 0; j < nrgba.Bounds().Dy(); j++ {
			for i := 0; i < nrgba.Bounds().Dx(); i++ {
				flippedImg.Set(i, bounds.Dy()-j, nrgba.At(i, j))
			}
		}

		return flippedImg, nil
	}

	return nrgba, nil
}
