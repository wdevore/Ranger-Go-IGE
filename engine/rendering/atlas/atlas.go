package atlas

import "unsafe"

const (
	floatSize = int(unsafe.Sizeof(float32(0)))
	uintSize  = int(unsafe.Sizeof(uint32(0)))
)
