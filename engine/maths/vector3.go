package maths

import (
	"fmt"
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// vector3 contains base components
type vector3 struct {
	x, y, z float32
}

// NewVector3 creates a Vector3 initialized to 0.0, 0.0, 0.0
func NewVector3() api.IVector3 {
	v := new(vector3)
	v.x = 0.0
	v.y = 0.0
	v.z = 0.0
	return v
}

// NewVector3With3Components creates a Vector3 initialized with x,y,z
func NewVector3With3Components(x, y, z float32) api.IVector3 {
	v := new(vector3)
	v.x = x
	v.y = y
	v.z = z
	return v
}

// NewVector3With2Components creates a Vector3 initialized with x,y and z = 0.0
func NewVector3With2Components(x, y float32) api.IVector3 {
	v := new(vector3)
	v.x = x
	v.y = y
	v.z = 0.0
	return v
}

// Clone returns a new copy this vector
func (v *vector3) Clone() api.IVector3 {
	c := new(vector3)
	c.x = v.x
	c.y = v.y
	c.z = v.z
	return c
}

// Set3Components modifies x,y,z
func (v *vector3) Set3Components(x, y, z float32) {
	v.x = x
	v.y = y
	v.z = z
}

// Set2Components modifies x,y only
func (v *vector3) Set2Components(x, y float32) {
	v.x = x
	v.y = y
}

// X return x component
func (v *vector3) X() float32 {
	return v.x
}

// Y return x component
func (v *vector3) Y() float32 {
	return v.y
}

// Z return z component
func (v *vector3) Z() float32 {
	return v.x
}

// Components2D returns x,y
func (v *vector3) Components2D() (x, y float32) {
	return v.x, v.y
}

// Components3D returns x,y
func (v *vector3) Components3D() (x, y, z float32) {
	return v.x, v.y, v.z
}

// Set modifies x,y,z from source
func (v *vector3) Set(source api.IVector3) {
	v.x = source.X()
	v.y = source.Y()
	v.z = source.Z()
}

// Add a Vector3 to this vector
func (v *vector3) Add(src api.IVector3) {
	v.x += src.X()
	v.y += src.Y()
	v.z += src.Z()
}

// Add2Components adds x and y to this vector
func (v *vector3) Add2Components(x, y float32) {
	v.x += x
	v.y += y
}

// Sub subtracts a Vector3 to this vector
func (v *vector3) Sub(src api.IVector3) {
	v.x -= src.X()
	v.y -= src.Y()
	v.z -= src.Z()
}

// Sub2Components subtracts x and y to this vector
func (v *vector3) Sub2Components(x, y float32) {
	v.x -= x
	v.y -= y
}

// ScaleBy scales this vector by s
func (v *vector3) ScaleBy(s float32) {
	v.x *= s
	v.y *= s
	v.z *= s
}

// ScaleBy2Components scales this vector by sx and sy
func (v *vector3) ScaleBy2Components(sx, sy float32) {
	v.x *= sx
	v.y *= sy
}

// MulAdd scales and adds src to this vector
func (v *vector3) MulAdd(src api.IVector3, scalar float32) {
	v.x += src.X() * scalar
	v.y += src.Y() * scalar
	v.z += src.Z() * scalar
}

// Length returns the euclidean length
func Length(x, y, z float32) float32 {
	return float32(math.Sqrt(float64(x*x + y*y + z*z)))
}

// Length returns the euclidean length
func (v *vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
}

// LengthSquared returns the euclidean length squared
func LengthSquared(x, y, z float32) float32 {
	return x*x + y*y + z*z
}

// LengthSquared returns the euclidean length squared
func (v *vector3) LengthSquared() float32 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Equal makes an exact equality check. Use EqEpsilon, it is more realistic.
func (v *vector3) Equal(other api.IVector3) bool {
	return v.x == other.X() && v.y == other.Y() && v.z == other.Z()
}

// EqEpsilon makes an approximate equality check. Preferred
func (v *vector3) EqEpsilon(other api.IVector3) bool {
	return (v.x-other.X()) < Epsilon && (v.y-other.Y()) < Epsilon && (v.z-other.Z()) < Epsilon
}

// Distance finds the euclidean distance between the two specified vectors
func Distance(x1, y1, z1, x2, y2, z2 float32) float32 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1

	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

// Distance finds the euclidean distance between the two specified vectors
func (v *vector3) Distance(src api.IVector3) float32 {
	a := src.X() - v.x
	b := src.Y() - v.y
	c := src.Z() - v.z

	return float32(math.Sqrt(float64(a*a + b*b + c*c)))
}

// DistanceSquared finds the euclidean distance between the two specified vectors squared
func DistanceSquared(x1, y1, z1, x2, y2, z2 float32) float32 {
	a := x2 - x1
	b := y2 - y1
	c := z2 - z1

	return a*a + b*b + c*c
}

// DistanceSquared finds the euclidean distance between the two specified vectors squared
func (v *vector3) DistanceSquared(src api.IVector3) float32 {
	a := src.X() - v.x
	b := src.Y() - v.y
	c := src.Z() - v.z

	return a*a + b*b + c*c
}

// Dot returns the product between the two vectors
func Dot(x1, y1, z1, x2, y2, z2 float32) float32 {
	return x1*x2 + y1*y2 + z1*z2
}

// DotByComponent returns the product between the two vectors
func (v *vector3) DotByComponent(x, y, z float32) float32 {
	return v.x*x + v.y*y + v.z*z
}

// Dot returns the product between the two vectors
func (v *vector3) Dot(o api.IVector3) float32 {
	return v.x*o.X() + v.y*o.Y() + v.z*o.Z()
}

// Cross sets this vector to the cross product between it and the other vector.
func (v *vector3) Cross(o api.IVector3) {
	v.Set3Components(
		v.y*o.Z()-v.z*o.Y(),
		v.z*o.X()-v.x*o.Z(),
		v.x*o.Y()-v.y*o.X())
}

// --------------------------------------------------------------------------
// Transforms
// --------------------------------------------------------------------------

// |M00 M01 M02 M03|   |x|
// |M10 M11 M12 M13| x |y|
// |M20 M21 M22 M23|   |z|
// |M30 M31 M32 M33|   |1|

// Mul left-multiplies the vector by the given matrix, assuming the fourth (w) component of the vector is 1.
func (v *vector3) Mul(m api.IMatrix4) {
	me := m.Matrix()
	v.Set3Components(
		v.x*me[M00]+v.y*me[M01]+v.z*me[M02]+me[M03],
		v.x*me[M10]+v.y*me[M11]+v.z*me[M12]+me[M13],
		v.x*me[M20]+v.y*me[M21]+v.z*me[M22]+me[M23])
}

func (v vector3) String() string {
	return fmt.Sprintf("<%7.3f, %7.3f, %7.3f>", v.x, v.y, v.z)
}
