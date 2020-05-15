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
	count  int32
}

// NewAtlasShape creates a new vector shape
func NewAtlasShape() api.IAtlasShape {
	vs := new(AtlasShape)
	return vs
}

// SetOffset scales offset by size of an uint32
func (vs *AtlasShape) SetOffset(offset int) {
	vs.offset = offset * int(unsafe.Sizeof(uint32(0)))
}

// Offset returns calculated offset
func (vs *AtlasShape) Offset() int {
	return vs.offset
}

// Name returns name
func (vs *AtlasShape) Name() string {
	return vs.name
}

// SetName sets name
func (vs *AtlasShape) SetName(n string) {
	vs.name = n
}

// PrimitiveMode returns OpenGL primative mode
func (vs *AtlasShape) PrimitiveMode() uint32 {
	return vs.primitiveMode
}

// SetPrimitiveMode sets name
func (vs *AtlasShape) SetPrimitiveMode(m uint32) {
	vs.primitiveMode = m
}

// Count returns count
func (vs *AtlasShape) Count() int32 {
	return vs.count
}

// SetCount sets index count
func (vs *AtlasShape) SetCount(c int32) {
	vs.count = c
}
