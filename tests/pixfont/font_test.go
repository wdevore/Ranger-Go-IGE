package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/pbnjay/pixfont"
)

func TestRunner(t *testing.T) {
	drawFont(t)
}

// https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go

func drawFont(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 128, 128))

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

	f, _ := os.OpenFile("font9x9.png", os.O_CREATE|os.O_RDWR, 0644)
	png.Encode(f, img)
}
