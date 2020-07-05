// Package rendering defines TBO features of shaders.
package rendering

import (
	"fmt"
	"image"
	"runtime"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v4.5-core/gl"
)

// TBO represents a shader's Texture features.
type TBO struct {
	tboID uint32 // GLuint
}

// NewTBO creates a empty TBO
func NewTBO() *TBO {
	o := new(TBO)

	gl.GenTextures(1, &o.tboID)

	runtime.SetFinalizer(o, (*TBO).delete)

	return o
}

// TboID returns internal TBO buffer id.
func (t *TBO) TboID() uint32 {
	return t.tboID
}

// Use bind vertex array to Id
func (t *TBO) Use() {
	// activate the texture unit first before binding
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.tboID)
}

// UnUse ...
func (t *TBO) UnUse() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Bind binds the pixels to buffer
// *image.NRGBA has a .Pix() that returns the []uint8 array.
func (t *TBO) Bind(width, height int32, smooth bool, pixels []uint8) {
	t.Use()

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	if smooth {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}
}

// BindUsingImage binds the image to buffer
func (t *TBO) BindUsingImage(image *image.NRGBA, smooth bool) {
	t.Use()

	width := int32(image.Bounds().Dx())
	height := int32(image.Bounds().Dy())

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(image.Pix))

	if smooth {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}
}

func (t *TBO) delete() {
	fmt.Println("TBO delete")
	mainthread.CallNonBlock(func() {
		gl.DeleteTextures(1, &t.tboID)
	})
}
