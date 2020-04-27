package maths

import (
	"fmt"
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

const (
	// M00 XX: Typically the unrotated X component for scaling, also the cosine of the
	// angle when rotated on the Y and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source X component
	// and added to the target X component.
	M00 = 0
	// M01 XY: Typically the negative sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Y component and added to the target X component.
	M01 = 4
	// M02 XZ: Typically the sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Z component and added to the target X component.
	M02 = 8
	// M03 XW: Typically the translation of the X component.
	// On Vector3 multiplication this value is added to the target X component.
	M03 = 12

	// M10 YX: Typically the sine of the angle when rotated on the Z axis.
	// On Vector3 multiplication this value is multiplied with the
	// source X component and added to the target Y component.
	M10 = 1
	// M11 YY: Typically the unrotated Y component for scaling, also the cosine
	// of the angle when rotated on the X and/or Z axis. On
	// Vector3 multiplication this value is multiplied with the source Y
	// component and added to the target Y component.
	M11 = 5
	// M12 YZ: Typically the negative sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied
	// with the source Z component and added to the target Y component.
	M12 = 9
	// M13 YW: Typically the translation of the Y component.
	// On Vector3 multiplication this value is added to the target Y component.
	M13 = 13

	// M20 ZX: Typically the negative sine of the angle when rotated on the Y axis.
	// On Vector3 multiplication this value is multiplied
	// with the source X component and added to the target Z component.
	M20 = 2
	// M21 ZY: Typical the sine of the angle when rotated on the X axis.
	// On Vector3 multiplication this value is multiplied with the
	// source Y component and added to the target Z component.
	M21 = 6
	// M22 ZZ: Typically the unrotated Z component for scaling, also the cosine of the
	// angle when rotated on the X and/or Y axis.
	// On Vector3 multiplication this value is multiplied with the source Z component
	// and added to the target Z component.
	M22 = 10
	// M23 ZW: Typically the translation of the Z component.
	// On Vector3 multiplication this value is added to the target Z component.
	M23 = 14

	// M30 WX: Typically the value zero. On Vector3 multiplication this value is ignored.
	M30 = 3
	// M31 WY: Typically the value zero. On Vector3 multiplication this value is ignored.
	M31 = 7
	// M32 WZ: Typically the value zero. On Vector3 multiplication this value is ignored.
	M32 = 11
	// M33 WW: Typically the value one. On Vector3 multiplication this value is ignored.
	M33 = 15
)

// matrix4 represents a column major opengl array.
type matrix4 struct {
	e [16]float32

	// Rotation is in radians
	Rotation float32
	Scale    vector3
}

// A temporary matrix for multiplication
var temp = NewMatrix4()

// NewMatrix4 creates a Matrix4 initialized to an identity matrix
func NewMatrix4() api.IMatrix4 {
	m := new(matrix4)
	m.ToIdentity()
	return m
}

// E returns the internal 4x4 matrix
func (m *matrix4) E() [16]float32 {
	return m.e
}

// --------------------------------------------------------------------------
// Translation
// --------------------------------------------------------------------------

// TranslateBy adds a translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) TranslateBy(v api.IVector3) {
	m.e[M03] += v.X()
	m.e[M13] += v.Y()
	m.e[M23] += v.Z()
}

// TranslateBy3Comps adds a translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) TranslateBy3Comps(x, y, z float32) {
	m.e[M03] += x
	m.e[M13] += y
	m.e[M23] += z
}

// TranslateBy2Comps adds a translational component to the matrix in the 4th column.
// Z is unmodified. The other columns are unmodified.
func (m *matrix4) TranslateBy2Comps(x, y float32) {
	m.e[M03] += x
	m.e[M13] += y
}

// SetTranslateByVector sets the translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) SetTranslateByVector(v api.IVector3) {
	m.ToIdentity()
	m.e[M03] = v.X()
	m.e[M13] = v.Y()
	m.e[M23] = v.Z()
}

// SetTranslate3Comp sets the translational component to the matrix in the 4th column.
// The other columns are unmodified.
func (m *matrix4) SetTranslate3Comp(x, y, z float32) {
	m.ToIdentity()

	m.e[M03] = x
	m.e[M13] = y
	m.e[M23] = z
}

// GetTranslation returns the translational components in 'out' Vector3 field.
func (m *matrix4) GetTranslation(out api.IVector3) {
	out.Set3Components(m.e[M03], m.e[M13], m.e[M23])
}

// --------------------------------------------------------------------------
// Rotation
// --------------------------------------------------------------------------

// SetRotation set a rotation matrix about Z axis. 'angle' is specified
// in radians.
//
//      [  M00  M01   _    _   ]
//      [  M10  M11   _    _   ]
//      [   _    _    _    _   ]
//      [   _    _    _    _   ]
func (m *matrix4) SetRotation(angle float32) {
	m.ToIdentity()

	if angle == 0 {
		return
	}

	m.Rotation = angle

	// Column major
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	m.e[M00] = c
	m.e[M01] = -s
	m.e[M10] = s
	m.e[M11] = c
}

// RotateBy postmultiplies this matrix with a (counter-clockwise) rotation matrix whose
// angle is specified in radians.
func (m *matrix4) RotateBy(angle float32) {
	if angle == 0.0 {
		return
	}

	m.Rotation += angle

	// Column major
	c := float32(math.Cos(float64(angle)))
	s := float32(math.Sin(float64(angle)))

	e := temp.E()
	e[M00] = c
	e[M01] = -s
	e[M02] = 0.0
	e[M03] = 0.0
	e[M10] = s
	e[M11] = c
	e[M12] = 0.0
	e[M13] = 0.0
	e[M20] = 0.0
	e[M21] = 0.0
	e[M22] = 1.0
	e[M23] = 0.0
	e[M30] = 0.0
	e[M31] = 0.0
	e[M32] = 0.0
	e[M33] = 1.0

	m.PostMultiply(temp)
}

// --------------------------------------------------------------------------
// Scale
// --------------------------------------------------------------------------

// ScaleBy scales the scale components.
func (m *matrix4) ScaleBy(v api.IVector3) {
	m.Scale.Set(v)

	m.e[M00] *= v.X()
	m.e[M11] *= v.Y()
	m.e[M22] *= v.Z()
}

// SetScale sets the scale components of an identity matrix and captures
// scale values into Scale property.
func (m *matrix4) SetScale(v api.IVector3) {
	m.Scale.Set(v)

	m.ToIdentity()

	m.e[M00] = v.X()
	m.e[M11] = v.Y()
	m.e[M22] = v.Z()
}

// SetScale3Comp sets the scale components of an identity matrix and captures
// scale values into Scale property.
func (m *matrix4) SetScale3Comp(sx, sy, sz float32) {
	m.Scale.Set3Components(sx, sy, sz)

	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = sz
}

// SetScale2Comp sets the scale components of an identity matrix and captures
// scale values into Scale property where Z component = 1.0.
func (m *matrix4) SetScale2Comp(sx, sy float32) {
	m.Scale.Set3Components(sx, sy, 1.0)

	m.ToIdentity()

	m.e[M00] = sx
	m.e[M11] = sy
	m.e[M22] = 1.0
}

// GetScale returns the scale in 'out' field.
func (m *matrix4) GetScale(out api.IVector3) {
	out.Set(&m.Scale)
}

// PostScale postmultiplies this matrix with a scale matrix.
func (m *matrix4) PostScale(sx, sy, sz float32) {
	e := temp.E()

	e[M00] = sx
	e[M01] = 0
	e[M02] = 0
	e[M03] = 0
	e[M10] = 0
	e[M11] = sy
	e[M12] = 0
	e[M13] = 0
	e[M20] = 0
	e[M21] = 0
	e[M22] = sz
	e[M23] = 0
	e[M30] = 0
	e[M31] = 0
	e[M32] = 0
	e[M33] = 1

	m.PostMultiply(temp)
}

// --------------------------------------------------------------------------
// Transforms
// --------------------------------------------------------------------------

// --------------------------------------------------------------------------
// Matrix methods
// --------------------------------------------------------------------------

// Multiply multiplies a * b and places result into 'out', (i.e. out = a * b)
func Multiply(a, b, out api.IMatrix4) {
	oe := out.E()
	ae := a.E()
	be := b.E()

	oe[M00] = ae[M00]*be[M00] + ae[M01]*be[M10] + ae[M02]*be[M20] + ae[M03]*be[M30]
	oe[M01] = ae[M00]*be[M01] + ae[M01]*be[M11] + ae[M02]*be[M21] + ae[M03]*be[M31]
	oe[M02] = ae[M00]*be[M02] + ae[M01]*be[M12] + ae[M02]*be[M22] + ae[M03]*be[M32]
	oe[M03] = ae[M00]*be[M03] + ae[M01]*be[M13] + ae[M02]*be[M23] + ae[M03]*be[M33]
	oe[M10] = ae[M10]*be[M00] + ae[M11]*be[M10] + ae[M12]*be[M20] + ae[M13]*be[M30]
	oe[M11] = ae[M10]*be[M01] + ae[M11]*be[M11] + ae[M12]*be[M21] + ae[M13]*be[M31]
	oe[M12] = ae[M10]*be[M02] + ae[M11]*be[M12] + ae[M12]*be[M22] + ae[M13]*be[M32]
	oe[M13] = ae[M10]*be[M03] + ae[M11]*be[M13] + ae[M12]*be[M23] + ae[M13]*be[M33]
	oe[M20] = ae[M20]*be[M00] + ae[M21]*be[M10] + ae[M22]*be[M20] + ae[M23]*be[M30]
	oe[M21] = ae[M20]*be[M01] + ae[M21]*be[M11] + ae[M22]*be[M21] + ae[M23]*be[M31]
	oe[M22] = ae[M20]*be[M02] + ae[M21]*be[M12] + ae[M22]*be[M22] + ae[M23]*be[M32]
	oe[M23] = ae[M20]*be[M03] + ae[M21]*be[M13] + ae[M22]*be[M23] + ae[M23]*be[M33]
	oe[M30] = ae[M30]*be[M00] + ae[M31]*be[M10] + ae[M32]*be[M20] + ae[M33]*be[M30]
	oe[M31] = ae[M30]*be[M01] + ae[M31]*be[M11] + ae[M32]*be[M21] + ae[M33]*be[M31]
	oe[M32] = ae[M30]*be[M02] + ae[M31]*be[M12] + ae[M32]*be[M22] + ae[M33]*be[M32]
	oe[M33] = ae[M30]*be[M03] + ae[M31]*be[M13] + ae[M32]*be[M23] + ae[M33]*be[M33]
}

// Multiply multiplies a * b and places result into this matrix, (i.e. this = a * b)
func (m *matrix4) Multiply(a, b api.IMatrix4) {
	ae := a.E()
	be := b.E()

	m.e[M00] = ae[M00]*be[M00] + ae[M01]*be[M10] + ae[M02]*be[M20] + ae[M03]*be[M30]
	m.e[M01] = ae[M00]*be[M01] + ae[M01]*be[M11] + ae[M02]*be[M21] + ae[M03]*be[M31]
	m.e[M02] = ae[M00]*be[M02] + ae[M01]*be[M12] + ae[M02]*be[M22] + ae[M03]*be[M32]
	m.e[M03] = ae[M00]*be[M03] + ae[M01]*be[M13] + ae[M02]*be[M23] + ae[M03]*be[M33]
	m.e[M10] = ae[M10]*be[M00] + ae[M11]*be[M10] + ae[M12]*be[M20] + ae[M13]*be[M30]
	m.e[M11] = ae[M10]*be[M01] + ae[M11]*be[M11] + ae[M12]*be[M21] + ae[M13]*be[M31]
	m.e[M12] = ae[M10]*be[M02] + ae[M11]*be[M12] + ae[M12]*be[M22] + ae[M13]*be[M32]
	m.e[M13] = ae[M10]*be[M03] + ae[M11]*be[M13] + ae[M12]*be[M23] + ae[M13]*be[M33]
	m.e[M20] = ae[M20]*be[M00] + ae[M21]*be[M10] + ae[M22]*be[M20] + ae[M23]*be[M30]
	m.e[M21] = ae[M20]*be[M01] + ae[M21]*be[M11] + ae[M22]*be[M21] + ae[M23]*be[M31]
	m.e[M22] = ae[M20]*be[M02] + ae[M21]*be[M12] + ae[M22]*be[M22] + ae[M23]*be[M32]
	m.e[M23] = ae[M20]*be[M03] + ae[M21]*be[M13] + ae[M22]*be[M23] + ae[M23]*be[M33]
	m.e[M30] = ae[M30]*be[M00] + ae[M31]*be[M10] + ae[M32]*be[M20] + ae[M33]*be[M30]
	m.e[M31] = ae[M30]*be[M01] + ae[M31]*be[M11] + ae[M32]*be[M21] + ae[M33]*be[M31]
	m.e[M32] = ae[M30]*be[M02] + ae[M31]*be[M12] + ae[M32]*be[M22] + ae[M33]*be[M32]
	m.e[M33] = ae[M30]*be[M03] + ae[M31]*be[M13] + ae[M32]*be[M23] + ae[M33]*be[M33]
}

// PreMultiply premultiplies 'b' matrix with 'this' and places the result into 'this' matrix.
// (i.e. this = b * this)
func (m *matrix4) PreMultiply(b api.IMatrix4) {
	te := temp.E()
	be := b.E()

	te[M00] = be[M00]*m.e[M00] + be[M01]*m.e[M10] + be[M02]*m.e[M20] + be[M03]*m.e[M30]
	te[M01] = be[M00]*m.e[M01] + be[M01]*m.e[M11] + be[M02]*m.e[M21] + be[M03]*m.e[M31]
	te[M02] = be[M00]*m.e[M02] + be[M01]*m.e[M12] + be[M02]*m.e[M22] + be[M03]*m.e[M32]
	te[M03] = be[M00]*m.e[M03] + be[M01]*m.e[M13] + be[M02]*m.e[M23] + be[M03]*m.e[M33]
	te[M10] = be[M10]*m.e[M00] + be[M11]*m.e[M10] + be[M12]*m.e[M20] + be[M13]*m.e[M30]
	te[M11] = be[M10]*m.e[M01] + be[M11]*m.e[M11] + be[M12]*m.e[M21] + be[M13]*m.e[M31]
	te[M12] = be[M10]*m.e[M02] + be[M11]*m.e[M12] + be[M12]*m.e[M22] + be[M13]*m.e[M32]
	te[M13] = be[M10]*m.e[M03] + be[M11]*m.e[M13] + be[M12]*m.e[M23] + be[M13]*m.e[M33]
	te[M20] = be[M20]*m.e[M00] + be[M21]*m.e[M10] + be[M22]*m.e[M20] + be[M23]*m.e[M30]
	te[M21] = be[M20]*m.e[M01] + be[M21]*m.e[M11] + be[M22]*m.e[M21] + be[M23]*m.e[M31]
	te[M22] = be[M20]*m.e[M02] + be[M21]*m.e[M12] + be[M22]*m.e[M22] + be[M23]*m.e[M32]
	te[M23] = be[M20]*m.e[M03] + be[M21]*m.e[M13] + be[M22]*m.e[M23] + be[M23]*m.e[M33]
	te[M30] = be[M30]*m.e[M00] + be[M31]*m.e[M10] + be[M32]*m.e[M20] + be[M33]*m.e[M30]
	te[M31] = be[M30]*m.e[M01] + be[M31]*m.e[M11] + be[M32]*m.e[M21] + be[M33]*m.e[M31]
	te[M32] = be[M30]*m.e[M02] + be[M31]*m.e[M12] + be[M32]*m.e[M22] + be[M33]*m.e[M32]
	te[M33] = be[M30]*m.e[M03] + be[M31]*m.e[M13] + be[M32]*m.e[M23] + be[M33]*m.e[M33]

	// Place results in "this"
	m.e[M00] = te[M00]
	m.e[M01] = te[M01]
	m.e[M02] = te[M02]
	m.e[M03] = te[M03]
	m.e[M10] = te[M10]
	m.e[M11] = te[M11]
	m.e[M12] = te[M12]
	m.e[M13] = te[M13]
	m.e[M20] = te[M20]
	m.e[M21] = te[M21]
	m.e[M22] = te[M22]
	m.e[M23] = te[M23]
	m.e[M30] = te[M30]
	m.e[M31] = te[M31]
	m.e[M32] = te[M32]
	m.e[M33] = te[M33]
}

// PostMultiply postmultiplies 'b' matrix with 'this' and places the result into 'this' matrix.
// (i.e. this = this * b)
func (m *matrix4) PostMultiply(b api.IMatrix4) {
	te := temp.E()
	be := b.E()

	te[M00] = be[M00]*m.e[M00] + be[M01]*m.e[M10] + be[M02]*m.e[M20] + be[M03]*m.e[M30]
	te[M01] = be[M00]*m.e[M01] + be[M01]*m.e[M11] + be[M02]*m.e[M21] + be[M03]*m.e[M31]
	te[M02] = be[M00]*m.e[M02] + be[M01]*m.e[M12] + be[M02]*m.e[M22] + be[M03]*m.e[M32]
	te[M03] = be[M00]*m.e[M03] + be[M01]*m.e[M13] + be[M02]*m.e[M23] + be[M03]*m.e[M33]
	te[M10] = be[M10]*m.e[M00] + be[M11]*m.e[M10] + be[M12]*m.e[M20] + be[M13]*m.e[M30]
	te[M11] = be[M10]*m.e[M01] + be[M11]*m.e[M11] + be[M12]*m.e[M21] + be[M13]*m.e[M31]
	te[M12] = be[M10]*m.e[M02] + be[M11]*m.e[M12] + be[M12]*m.e[M22] + be[M13]*m.e[M32]
	te[M13] = be[M10]*m.e[M03] + be[M11]*m.e[M13] + be[M12]*m.e[M23] + be[M13]*m.e[M33]
	te[M20] = be[M20]*m.e[M00] + be[M21]*m.e[M10] + be[M22]*m.e[M20] + be[M23]*m.e[M30]
	te[M21] = be[M20]*m.e[M01] + be[M21]*m.e[M11] + be[M22]*m.e[M21] + be[M23]*m.e[M31]
	te[M22] = be[M20]*m.e[M02] + be[M21]*m.e[M12] + be[M22]*m.e[M22] + be[M23]*m.e[M32]
	te[M23] = be[M20]*m.e[M03] + be[M21]*m.e[M13] + be[M22]*m.e[M23] + be[M23]*m.e[M33]
	te[M30] = be[M30]*m.e[M00] + be[M31]*m.e[M10] + be[M32]*m.e[M20] + be[M33]*m.e[M30]
	te[M31] = be[M30]*m.e[M01] + be[M31]*m.e[M11] + be[M32]*m.e[M21] + be[M33]*m.e[M31]
	te[M32] = be[M30]*m.e[M02] + be[M31]*m.e[M12] + be[M32]*m.e[M22] + be[M33]*m.e[M32]
	te[M33] = be[M30]*m.e[M03] + be[M31]*m.e[M13] + be[M32]*m.e[M23] + be[M33]*m.e[M33]

	// Place results in "this"
	m.e[M00] = te[M00]
	m.e[M01] = te[M01]
	m.e[M02] = te[M02]
	m.e[M03] = te[M03]
	m.e[M10] = te[M10]
	m.e[M11] = te[M11]
	m.e[M12] = te[M12]
	m.e[M13] = te[M13]
	m.e[M20] = te[M20]
	m.e[M21] = te[M21]
	m.e[M22] = te[M22]
	m.e[M23] = te[M23]
	m.e[M30] = te[M30]
	m.e[M31] = te[M31]
	m.e[M32] = te[M32]
	m.e[M33] = te[M33]
}

// MultiplyIntoA multiplies a * b and places result into 'a', (i.e. a = a * b)
func MultiplyIntoA(a, b api.IMatrix4) {
	te := temp.E()
	ae := a.E()
	be := b.E()

	te[M00] = ae[M00]*be[M00] + ae[M01]*be[M10] + ae[M02]*be[M20] + ae[M03]*be[M30]
	te[M01] = ae[M00]*be[M01] + ae[M01]*be[M11] + ae[M02]*be[M21] + ae[M03]*be[M31]
	te[M02] = ae[M00]*be[M02] + ae[M01]*be[M12] + ae[M02]*be[M22] + ae[M03]*be[M32]
	te[M03] = ae[M00]*be[M03] + ae[M01]*be[M13] + ae[M02]*be[M23] + ae[M03]*be[M33]
	te[M10] = ae[M10]*be[M00] + ae[M11]*be[M10] + ae[M12]*be[M20] + ae[M13]*be[M30]
	te[M11] = ae[M10]*be[M01] + ae[M11]*be[M11] + ae[M12]*be[M21] + ae[M13]*be[M31]
	te[M12] = ae[M10]*be[M02] + ae[M11]*be[M12] + ae[M12]*be[M22] + ae[M13]*be[M32]
	te[M13] = ae[M10]*be[M03] + ae[M11]*be[M13] + ae[M12]*be[M23] + ae[M13]*be[M33]
	te[M20] = ae[M20]*be[M00] + ae[M21]*be[M10] + ae[M22]*be[M20] + ae[M23]*be[M30]
	te[M21] = ae[M20]*be[M01] + ae[M21]*be[M11] + ae[M22]*be[M21] + ae[M23]*be[M31]
	te[M22] = ae[M20]*be[M02] + ae[M21]*be[M12] + ae[M22]*be[M22] + ae[M23]*be[M32]
	te[M23] = ae[M20]*be[M03] + ae[M21]*be[M13] + ae[M22]*be[M23] + ae[M23]*be[M33]
	te[M30] = ae[M30]*be[M00] + ae[M31]*be[M10] + ae[M32]*be[M20] + ae[M33]*be[M30]
	te[M31] = ae[M30]*be[M01] + ae[M31]*be[M11] + ae[M32]*be[M21] + ae[M33]*be[M31]
	te[M32] = ae[M30]*be[M02] + ae[M31]*be[M12] + ae[M32]*be[M22] + ae[M33]*be[M32]
	te[M33] = ae[M30]*be[M03] + ae[M31]*be[M13] + ae[M32]*be[M23] + ae[M33]*be[M33]

	ae[M00] = te[M00]
	ae[M01] = te[M01]
	ae[M02] = te[M02]
	ae[M03] = te[M03]
	ae[M10] = te[M10]
	ae[M11] = te[M11]
	ae[M12] = te[M12]
	ae[M13] = te[M13]
	ae[M20] = te[M20]
	ae[M21] = te[M21]
	ae[M22] = te[M22]
	ae[M23] = te[M23]
	ae[M30] = te[M30]
	ae[M31] = te[M31]
	ae[M32] = te[M32]
	ae[M33] = te[M33]
}

// PostTranslate postmultiplies this matrix by a translation matrix.
// Postmultiplication is also used by OpenGL ES.
func (m *matrix4) PostTranslate(tx, ty, tz float32) {
	te := temp.E()

	te[M00] = 1.0
	te[M01] = 0.0
	te[M02] = 0.0
	te[M03] = tx
	te[M10] = 0.0
	te[M11] = 1.0
	te[M12] = 0.0
	te[M13] = ty
	te[M20] = 0.0
	te[M21] = 0.0
	te[M22] = 1.0
	te[M23] = tz
	te[M30] = 0.0
	te[M31] = 0.0
	te[M32] = 0.0
	te[M33] = 1.0

	// this = this * temp;
	m.PostMultiply(temp)
}

// --------------------------------------------------------------------------
// Projections
// --------------------------------------------------------------------------

// SetToOrtho sets the matrix for a 2d ortho graphic projection
func (m *matrix4) SetToOrtho(left, right, bottom, top, near, far float32) {
	m.ToIdentity()

	xorth := 2.0 / (right - left)
	yorth := 2.0 / (top - bottom)
	zorth := -2.0 / (far - near)

	tx := -(right + left) / (right - left)
	ty := -(top + bottom) / (top - bottom)
	tz := -(far + near) / (far - near)

	m.e[M00] = xorth
	m.e[M10] = 0.0
	m.e[M20] = 0.0
	m.e[M30] = 0.0
	m.e[M01] = 0.0
	m.e[M11] = yorth
	m.e[M21] = 0.0
	m.e[M31] = 0.0
	m.e[M02] = 0.0
	m.e[M12] = 0.0
	m.e[M22] = zorth
	m.e[M32] = 0.0
	m.e[M03] = tx
	m.e[M13] = ty
	m.e[M23] = tz
	m.e[M33] = 1.0
}

// --------------------------------------------------------------------------
// Misc
// --------------------------------------------------------------------------

// C returns a cell value based on Mxx index
func (m *matrix4) C(i int) float32 {
	return m.e[i]
}

// Clone returns a clone of this matrix
func (m *matrix4) Clone() api.IMatrix4 {
	c := new(matrix4)
	c.Set(m)
	return c
}

// Set copies src into this matrix
func (m *matrix4) Set(src api.IMatrix4) {
	se := src.E()

	m.e[M00] = se[M00]
	m.e[M01] = se[M01]
	m.e[M02] = se[M02]
	m.e[M03] = se[M03]

	m.e[M10] = se[M10]
	m.e[M11] = se[M11]
	m.e[M12] = se[M12]
	m.e[M13] = se[M13]

	m.e[M20] = se[M20]
	m.e[M21] = se[M21]
	m.e[M22] = se[M22]
	m.e[M23] = se[M23]

	m.e[M30] = se[M30]
	m.e[M31] = se[M31]
	m.e[M32] = se[M32]
	m.e[M33] = se[M33]
}

// ToIdentity set this matrix to the identity matrix
func (m *matrix4) ToIdentity() {
	m.e[M00] = 1.0
	m.e[M01] = 0.0
	m.e[M02] = 0.0
	m.e[M03] = 0.0

	m.e[M10] = 0.0
	m.e[M11] = 1.0
	m.e[M12] = 0.0
	m.e[M13] = 0.0

	m.e[M20] = 0.0
	m.e[M21] = 0.0
	m.e[M22] = 1.0
	m.e[M23] = 0.0

	m.e[M30] = 0.0
	m.e[M31] = 0.0
	m.e[M32] = 0.0
	m.e[M33] = 1.0
}

func (m matrix4) String() string {
	s := fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M00], m.e[M01], m.e[M02], m.e[M03])
	s += fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M10], m.e[M11], m.e[M12], m.e[M13])
	s += fmt.Sprintf("[%f, %f, %f, %f]\n", m.e[M20], m.e[M21], m.e[M22], m.e[M23])
	s += fmt.Sprintf("[%f, %f, %f, %f]", m.e[M30], m.e[M31], m.e[M32], m.e[M33])
	return s
}
