package rendering

import (
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// AtlasShape defines shape element attributes
type AtlasShape struct {
	name          string
	primitiveMode uint32
	// Offset is multiplied by the size of an Unsigned Int in preparation for
	// drawing.
	offset int
	count  int

	atlasObj api.IAtlasObject
}

// NewAtlasShape creates a new vector shape
func NewAtlasShape(atlasObj api.IAtlasObject) api.IAtlasShape {
	o := new(AtlasShape)
	o.atlasObj = atlasObj
	return o
}

// SetOffset scales offset by size of an uint32
func (a *AtlasShape) SetOffset(offset int) {
	a.offset = offset * int(unsafe.Sizeof(uint32(0)))
}

// Offset returns calculated offset
func (a *AtlasShape) Offset() int {
	return a.offset
}

// Name returns name
func (a *AtlasShape) Name() string {
	return a.name
}

// SetName sets name
func (a *AtlasShape) SetName(n string) {
	a.name = n
}

// PrimitiveMode returns OpenGL primative mode
func (a *AtlasShape) PrimitiveMode() uint32 {
	return a.primitiveMode
}

// SetPrimitiveMode sets name
func (a *AtlasShape) SetPrimitiveMode(m uint32) {
	a.primitiveMode = m
}

// Count returns count
func (a *AtlasShape) Count() int {
	return a.count
}

// SetCount sets index count
func (a *AtlasShape) SetCount(c int) {
	a.count = c
}

// SetVertex3D sets the atlas object's buffer data
func (a *AtlasShape) SetVertex3D(x, y, z float32, index int) {
	a.atlasObj.SetVertex(x, y, z, index)
}

// SetVertex2D sets the atlas object's buffer data using a Z = 0
func (a *AtlasShape) SetVertex2D(x, y float32, index int) {
	a.atlasObj.SetVertex(x, y, 0.0, index)
}

// Update modifies the VBO buffer
func (a *AtlasShape) Update(offset, vertexCount int) {
	a.atlasObj.Update(offset, vertexCount)
}
