package rendering

import "unsafe"

// VectorShape defines shape element attributes
type VectorShape struct {
	Name          string
	PrimitiveMode uint32
	// Offset is multiplied by the size of an Unsigned Int in preparation for
	// drawing.
	offset int
	Count  int32
}

// NewVectorShape creates a new vector shape
func NewVectorShape() *VectorShape {
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
