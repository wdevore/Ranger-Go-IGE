package api

// IDynamicText represents dynamic text type objects
type IDynamicText interface {
	Text() string
	SetText(string)
	SetPixelSize(size float32)
}
