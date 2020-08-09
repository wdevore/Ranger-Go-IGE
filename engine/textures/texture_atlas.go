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
	atlas         *image.NRGBA

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
func (t *textureAtlas) Build(relativePath string) {
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
	fmt.Println("TextureAtlas.Build loading: ", file)
	image, err := t.loadImage(file, true)

	if err != nil {
		panic(err)
	}

	t.atlas = image
}

// AtlasImage returns image atlas
func (t *textureAtlas) AtlasImage() *image.NRGBA {
	return t.atlas
}

// Name returns atlas name
func (t *textureAtlas) Name() string {
	return t.name
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

// Build setups the atlas based on manifest
// func (t *textureAtlas) build(relativePath string) {
// 	dataPath, err := filepath.Abs(relativePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	manifestFile, err := os.Open(dataPath + "/" + t.manifest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer manifestFile.Close()

// 	lines := []string{}

// 	scanner := bufio.NewScanner(manifestFile)
// 	for scanner.Scan() {
// 		lines = append(lines, scanner.Text())
// 	}

// 	textureFile := lines[0]

// 	/// Use TextureManager to access image
// 	textureImgFile := dataPath + "/" + textureFile
// 	fmt.Println("TextureAtlas: loading " + textureImgFile)

// 	t.atlas, err = t.loadImage(textureImgFile, true)
// 	if err != nil {
// 		panic(err)
// 	}

// 	s := strings.Split(lines[1], "x")
// 	t.width, _ = strconv.ParseInt(s[0], 10, 64)
// 	t.height, _ = strconv.ParseInt(s[1], 10, 64)

// 	for i := 2; i < len(lines); i++ {
// 		s = strings.Split(lines[i], "|")

// 		ts := NewSubTexture(s[0])

// 		coords := strings.Split(s[1], ":")

// 		for _, coord := range coords {
// 			xy := strings.Split(coord, ",")
// 			x, _ := strconv.ParseInt(xy[0], 10, 64)
// 			y, _ := strconv.ParseInt(xy[1], 10, 64)

// 			// Scale x,y to fit texture space
// 			sc := float32(x) / float32(t.width)
// 			tc := float32(y) / float32(t.height)
// 			// fmt.Println(sc, ",", tc)
// 			ts.textureCoords = append(ts.textureCoords, &textureCoord{s: sc, t: tc})
// 		}

// 		t.subTextures = append(t.subTextures, ts)
// 	}
// }
