package api

// IVector represents 2D vectors
// that have direction and magnitude
type IVector interface {
	Components() (float32, float32)
	// X sets the x component
	X() float32
	// Y sets the y component
	Y() float32
	// SetByComp sets by component
	SetByComp(x, y float32)
	// SetByPoint sets point using another point
	SetByPoint(IPoint)
	// SetByVector sets point using another vector
	SetByVector(IVector)

	SetByAngle(radians float64)

	// Length returns the square root length
	Length() float32
	// LengthSqr return the squared length
	LengthSqr() float32

	// Scale vector by s
	Scale(s float32)

	// Add offsets a this vector
	Add(x, y float32)
	// Sub offsets a this vector
	Sub(x, y float32)

	AddV(IVector)
	SubV(IVector)

	// Div vector by d
	Div(d float32)

	// AngleX computes the angle (radians) between this vector and v.
	AngleX(v IVector) float32

	// Normalize normalizes this vector, if the vector is zero
	// then nothing happens
	Normalize()

	// SetDirection set the direction of the vector, however,
	// it will erase the magnitude
	SetDirection(radians float64)

	// CrossCW computes the cross-product faster in the CW direction
	CrossCW()
	// CrossCCW computes the cross-product faster in the CCW direction
	CrossCCW()
}
