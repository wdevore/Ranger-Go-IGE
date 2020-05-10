package api

// IVectorShape represents a vector shape
type IVectorShape interface {
	SetOffset(offset int)
	Offset() int

	Name() string
	SetName(string)

	PrimitiveMode() uint32
	SetPrimitiveMode(uint32)

	Count() int32
	SetCount(int32)
}
