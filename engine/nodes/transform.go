package nodes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// Transform holds the transform properties and methods.
type Transform struct {
	position api.IPoint
	rotation float64
	scale    api.IPoint

	aft     api.IAffineTransform
	inverse api.IAffineTransform
}

func (t *Transform) initializeTransform() {
	t.position = geometry.NewPoint()
	t.scale = geometry.NewPointUsing(1.0, 1.0)

	t.aft = maths.NewTransform()
	t.inverse = maths.NewTransform()
}

// AffineTransform returns this node's transform matrix
func (t *Transform) AffineTransform() api.IAffineTransform {
	return t.aft
}

// InverseTransform returns an inverted cached version of "transform"
func (t *Transform) InverseTransform() api.IAffineTransform {
	return t.inverse
}

// SetPosition set the translation components of the matrix
func (t *Transform) SetPosition(x, y float32) {
	t.position.SetByComp(x, y)
}

// Position returns the position independent of the matrix
func (t *Transform) Position() api.IPoint {
	return t.position
}

// SetRotation set the rotation given as radians
func (t *Transform) SetRotation(radians float64) {
	t.rotation = radians
}

// Rotation is the current rotation in radians
func (t *Transform) Rotation() float64 {
	return t.rotation
}

// SetScale sets the scale uniformly for x and y
func (t *Transform) SetScale(scale float32) {
	t.scale.SetByComp(scale, scale)
}

// SetScaleComps sets the scale Non-uniformly for x and y
func (t *Transform) SetScaleComps(sx, sy float32) {
	t.scale.SetByComp(sx, sy)
}

// Scale returns the X scale component for uniform scales.
func (t *Transform) Scale() float32 {
	return t.scale.X()
}

// ScaleComps returns X and Y scale components for non-uniform scales.
func (t *Transform) ScaleComps() (float32, float32) {
	return t.scale.X(), t.scale.Y()
}

// CalcFilteredTransform performs a filter transform calculation.
func (t *Transform) CalcFilteredTransform(excludeTranslation bool,
	excludeRotation bool,
	excludeScale bool,
	aft api.IAffineTransform) {
	aft.ToIdentity()

	if !excludeTranslation {
		aft.MakeTranslate(t.position.X(), t.position.Y())
	}

	if !excludeRotation && t.rotation != 0.0 {
		aft.Rotate(t.rotation)
	}

	if !excludeScale && (t.scale.X() != 0.0 || t.scale.Y() != 0.0) {
		aft.Scale(t.scale.X(), t.scale.Y())
	}
}
