package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func TestRunner(t *testing.T) {
	drawFont(t)
}

// https://stackoverflow.com/questions/38299930/how-to-add-a-simple-text-label-to-an-image-in-go

func drawFont(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 128, 128))
	y := 13
	i := y - 3
	o := 1
	addLabel(img, o, i, "abcdefghijklmnopqr")
	i += y
	addLabel(img, o, i, "stuvwxyzABCDEFGHIJ")
	i += y
	addLabel(img, o, i, "KLMNOPQRSTUVWXYZ01")
	i += y
	addLabel(img, o, i, "23456789~`!@#$%^&*")
	i += y
	addLabel(img, o, i, "()_+-={}|[]\\:\";'<>")
	i += y
	addLabel(img, o, i, "?,./")

	f, err := os.Create("hello-go.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
