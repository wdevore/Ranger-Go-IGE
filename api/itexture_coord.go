package api

// ITextureCoord is an s,t coordinate
type ITextureCoord interface {
	Set(s, t float32)
	ST() (float32, float32)
	S() float32
	T() float32
}
