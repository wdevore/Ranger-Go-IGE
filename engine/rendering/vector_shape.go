package rendering

import (
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// VectorShape defines shape element attributes
type VectorShape struct {
	name          string
	primitiveMode uint32
	// Offset is multiplied by the size of an Unsigned Int in preparation for
	// drawing.
	offset int
	count  int32
}

// NewVectorShape creates a new vector shape
func NewVectorShape() api.IVectorShape {
	vs := new(VectorShape)
	return vs
}

// SetOffset scales offset by size of an uint32
func (vs *VectorShape) SetOffset(offset int) {
	vs.offset = offset * int(unsafe.Sizeof(uint32(0)))
}

// Offset returns calculated offset
func (vs *VectorShape) Offset() int {
	return vs.offset
}

// Name returns name
func (vs *VectorShape) Name() string {
	return vs.name
}

// SetName sets name
func (vs *VectorShape) SetName(n string) {
	vs.name = n
}

// PrimitiveMode returns OpenGL primative mode
func (vs *VectorShape) PrimitiveMode() uint32 {
	return vs.primitiveMode
}

// SetPrimitiveMode sets name
func (vs *VectorShape) SetPrimitiveMode(m uint32) {
	vs.primitiveMode = m
}

// Count returns count
func (vs *VectorShape) Count() int32 {
	return vs.count
}

// SetCount sets index count
func (vs *VectorShape) SetCount(c int32) {
	vs.count = c
}
