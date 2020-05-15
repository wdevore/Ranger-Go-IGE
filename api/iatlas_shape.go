package api

// IAtlasShape represents a vector shape
type IAtlasShape interface {
	SetOffset(offset int)
	Offset() int

	Name() string
	SetName(string)

	PrimitiveMode() uint32
	SetPrimitiveMode(uint32)

	Count() int32
	SetCount(int32)
}
