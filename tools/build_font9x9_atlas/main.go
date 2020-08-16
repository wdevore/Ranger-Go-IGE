package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/pbnjay/pixfont"
)

const (
	// Outputs: two files
	textureAtlasFile = "font9x9_texture_atlas.png"
	textureManifest  = "font9x9_texture_manifest.json"
)

func main() {

	txf, _ := os.Create(textureAtlasFile)
	defer txf.Close()

	// Now build and write JSON
	packSize := 128

	// Build font texture atlas
	atlas := buildFont2(packSize)

	png.Encode(txf, atlas)

	writeAsJSONFile(textureManifest, textureAtlasFile, packSize, packSize)
}

func buildFont(packSize int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, packSize, packSize))

	clr := color.White
	x := 1
	y := 1
	c := 97
	for rc := 0; rc < 14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	c = 65
	for rc := 0; rc < 14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	c = 33
	for rc := 0; rc < 14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 14; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 4; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	c = 91
	for rc := 0; rc < 6; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	c = 123
	for rc := 0; rc < 4; rc++ {
		pixfont.DrawString(img, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	return img
}

func buildFont2(packSize int) *image.RGBA {
	runeSize := 9
	imgFlipRune := image.NewRGBA(image.Rect(0, 0, runeSize, runeSize))

	atlas := image.NewRGBA(image.Rect(0, 0, packSize, packSize))

	clr := color.White
	x := 1
	y := 1
	c := 97
	for rc := 0; rc < 14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		// pixfont.DrawString(atlas, x, y, string(rune(c)), clr)
		x += 9
		c++
	}

	x = 1
	y += 9
	c = 65
	for rc := 0; rc < 14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	x = 1
	y += 9
	c = 33
	for rc := 0; rc < 14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 14; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	x = 1
	y += 9
	for rc := 0; rc < 4; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	c = 91
	for rc := 0; rc < 6; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	c = 123
	for rc := 0; rc < 4; rc++ {
		drawRune(runeSize, c, x, y, clr, imgFlipRune, atlas)
		x += 9
		c++
	}

	return atlas
}

func drawRune(runeSize, c, x, y int, colr color.Color, flippedRune, atlas *image.RGBA) {
	imgRune := image.NewRGBA(image.Rect(0, 0, runeSize-1, runeSize-1))

	pixfont.DrawString(imgRune, 0, 0, string(rune(c)), colr)

	for j := 0; j < imgRune.Bounds().Dy(); j++ {
		for i := 0; i < imgRune.Bounds().Dx(); i++ {
			flippedRune.Set(i, imgRune.Bounds().Dy()-j, imgRune.At(i, j))
		}
	}

	target := image.Rect(x, y-1, x+runeSize, y+runeSize)

	draw.Draw(atlas, target, flippedRune, flippedRune.Bounds().Min, draw.Src)

}
