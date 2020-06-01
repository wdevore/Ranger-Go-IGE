package api

// IZoomTransform represents 2D zoom transform
type IZoomTransform interface {
	// GetTransform updates and returns the internal transform.
	GetTransform() IAffineTransform

	// Update modifies the internal transform state based on current values.
	Update()

	// SetPosition is an absolute position. Typically you would use TranslateBy.
	SetPosition(x, y float32)

	// Scale returns the current scale factor
	Scale() float32

	PsuedoScale() float32

	// SetScale sets the scale based on the current scale value making
	// this a relative scale.
	SetScale(scale float32)

	// SetAt sets the center zoom point.
	SetAt(x, y float32)

	// ZoomBy performs a relative zoom based on the current scale/zoom.
	ZoomBy(dx, dy float32)

	// TranslateBy is a relative positional translation.
	TranslateBy(dx, dy float32)
}
