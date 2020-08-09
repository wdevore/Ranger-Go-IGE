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
func (t *TBO) Use(textureIndex int) {
	// activate the texture unit first before binding. There can be anywhere from 8 to 32
	gl.ActiveTexture(indexToTextureID(textureIndex))
	gl.BindTexture(gl.TEXTURE_2D, t.tboID)
}

// UnUse ...
func (t *TBO) UnUse() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// Bind binds the pixels to buffer
// *image.NRGBA has a .Pix() that returns the []uint8 array.
func (t *TBO) Bind(width, height int32, smooth bool, pixels []uint8) {
	t.Use(0)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	if smooth {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
}

// BindUsingImage binds the image to buffer
func (t *TBO) BindUsingImage(image *image.NRGBA, smooth bool) {
	t.Use(0)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	width := int32(image.Bounds().Dx())
	height := int32(image.Bounds().Dy())

	if smooth {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	// Give the image to OpenGL
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(image.Pix))
	// gl.GenerateMipmap(gl.TEXTURE_2D)
}

// BindTextureVbo binds the vertex attributes for xyzst
func (t *TBO) BindTextureVbo(points []float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.DYNAMIC_DRAW)

	sizeOfFloat := int32(4)

	// If the data per-vertex is (x,y,z,s,t = 5) then
	// Stride = 5 * size of float
	// OR
	// If the data per-vertex is (x,y,z,r,g,b,s,t = 8) then
	// Stride = 8 * size of float

	// Our data layout is x,y,z,s,t
	stride := 5 * sizeOfFloat

	// position attribute
	size := int32(3)   // x,y,z
	offset := int32(0) // position is first thus this attrib is offset by 0
	attribIndex := uint32(0)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset)))
	gl.EnableVertexAttribArray(0)

	// texture coord attribute is offset by 3 (i.e. x,y,z)
	size = int32(2)   // s,t
	offset = int32(3) // the preceeding component size = 3, thus this attrib is offset by 3
	attribIndex = uint32(1)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset*sizeOfFloat)))
	gl.EnableVertexAttribArray(1)
}

// BindTextureVbo2 binds the vertex attributes for xyst
func (t *TBO) BindTextureVbo2(points []float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.DYNAMIC_DRAW)

	sizeOfFloat := int32(4)

	// If the data per-vertex is (x,y,s,t = 4) then
	stride := 4 * sizeOfFloat

	// position attribute
	size := int32(2)   // x,y
	offset := int32(0) // position is first thus this attrib is offset by 0
	attribIndex := uint32(0)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset)))
	gl.EnableVertexAttribArray(0)

	// texture coord attribute is offset by 2 (i.e. x,y)
	size = int32(2)   // s,t
	offset = int32(2) // the preceeding component size = 2, thus this attrib is offset by 2
	attribIndex = uint32(1)
	gl.VertexAttribPointer(attribIndex, size, gl.FLOAT, false, stride, gl.PtrOffset(int(offset*sizeOfFloat)))
	gl.EnableVertexAttribArray(1)
}

// UpdateTextureVbo moves any modified data to the buffer.
func (t *TBO) UpdateTextureVbo(data []float32, vbo uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(data)*4, gl.Ptr(data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (t *TBO) delete() {
	fmt.Println("TBO delete")
	mainthread.CallNonBlock(func() {
		gl.DeleteTextures(1, &t.tboID)
	})
}

func indexToTextureID(index int) (textID uint32) {
	textID = gl.TEXTURE0

	switch index {
	case 0:
		textID = gl.TEXTURE0
	case 1:
		textID = gl.TEXTURE1
	case 2:
		textID = gl.TEXTURE2
	case 3:
		textID = gl.TEXTURE3
	case 4:
		textID = gl.TEXTURE4
	case 5:
		textID = gl.TEXTURE5
	case 6:
		textID = gl.TEXTURE6
	case 7:
		textID = gl.TEXTURE7
	case 8:
		textID = gl.TEXTURE8
	case 9:
		textID = gl.TEXTURE9
	case 10:
		textID = gl.TEXTURE10
	case 11:
		textID = gl.TEXTURE11
	case 12:
		textID = gl.TEXTURE12
	case 13:
		textID = gl.TEXTURE13
	case 14:
		textID = gl.TEXTURE14
	case 15:
		textID = gl.TEXTURE15
	case 16:
		textID = gl.TEXTURE16
	case 17:
		textID = gl.TEXTURE17
	case 18:
		textID = gl.TEXTURE18
	case 19:
		textID = gl.TEXTURE19
	case 20:
		textID = gl.TEXTURE20
	case 21:
		textID = gl.TEXTURE21
	case 22:
		textID = gl.TEXTURE22
	case 23:
		textID = gl.TEXTURE23
	case 24:
		textID = gl.TEXTURE24
	case 25:
		textID = gl.TEXTURE25
	case 26:
		textID = gl.TEXTURE26
	case 27:
		textID = gl.TEXTURE27
	case 28:
		textID = gl.TEXTURE28
	case 29:
		textID = gl.TEXTURE29
	case 30:
		textID = gl.TEXTURE30
	case 31:
		textID = gl.TEXTURE31
	}

	return textID
}
