package maths

import (
	"fmt"
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// A minified affine transform. This is generally passed to
// OpenGL via uniformMatrix3fv()
//
// Column major (passed to OpenGL)
// [ a, b, i, c, d, j, e, f, k ]
//
//     x'   |a c e|   |x|
//     y' = |b d f| x |y|
//     1    |0 0 1|   |1|
//           i j k
//
//  tx = e, ty = f
//
// Column vectors:
// 1. Used by OpenGL
// 2. Matrix applied to left of vector, ex: M * v
// 3. Concatenated transformations are "interpreted" as right to left
//    MT * MR * MS would interpret as
//    Scale(MS) then Rotate(MR) then Translate(MT)
// which is very different in meaning than:
//    Translate then Rotate then Scale
// Each produces a different visual effect.
// The approach that Ranger uses is the "right to left" post multiply.

// https://developer.mozilla.org/en-US/docs/Web/CSS/transform-function/matrix
// | a c e |    | a c e |      | a c 0 e |     | ma mc m8  me  |
// | b d f |    | b d f |  =>  | b d 0 f | ==> | mb md m9  mf  |
// 	            | 0 0 1 |      | 0 0 1 0 |     | m2 m6 m10 m14 |
// 	                           | 0 0 0 1 |     | m3 m7 m11 m15 |

// For a 3x3
// const (
// 	ma  = 0
// 	mb  = 1
// 	mi  = 2 // 0
// 	mc  = 3
// 	md  = 4
// 	mj  = 5 // 0
// 	mtx = 6 // e
// 	mty = 7 // f
// 	mk  = 8 // 1
// )

// For a 4x4
const (
	ma = 0 // a
	mb = 1 // b
	m2 = 2 // 0
	m3 = 3 // 0
	mc = 4 // c
	md = 5 // d
	m6 = 6 // 0
	m7 = 7 // 0

	m8  = 8  // 0
	m9  = 9  // 0
	m10 = 10 // 1
	m11 = 11 // 0
	me  = 12 // e
	mf  = 13 // f
	m14 = 14 // 0
	m15 = 15 // 1
)

type affineTransform struct {
	// [a, b, 0, c, d, 0, e, f, 1]
	// m [9]float32

	// [a, b, 0, 0, c, d, 0, 0, 0, 0, 1, 0, 0, e, f, 0, 1]
	m [16]float32
}

// NewTransform constructs an Identity Affine Transform matrix
func NewTransform() api.IAffineTransform {
	o := new(affineTransform)
	o.ToIdentity()
	return o
}

func (at *affineTransform) Matrix() *([16]float32) {
	return &at.m
}

// ----------------------------------------------------------
// Methods
// ----------------------------------------------------------

// | 1 0 0 0 |     | a c 0 e |  | ma mc m8  me  |
// | 0 1 0 0 | ==> | b d 0 f |  | mb md m9  mf  |
// | 0 0 1 0 |     | 0 0 1 0 |  | m2 m6 m10 m14 |
// | 0 0 0 1 |     | 0 0 0 1 |  | m3 m7 m11 m15 |
func (at *affineTransform) ToIdentity() {
	at.m[ma] = 1.0
	at.m[mb] = 0.0
	at.m[m2] = 0.0
	at.m[m3] = 0.0

	at.m[mc] = 0.0
	at.m[md] = 1.0
	at.m[m6] = 0.0
	at.m[m7] = 0.0

	at.m[m8] = 0.0
	at.m[m9] = 0.0
	at.m[m10] = 1.0
	at.m[m11] = 0.0

	at.m[me] = 0.0
	at.m[mf] = 0.0
	at.m[m14] = 0.0
	at.m[m15] = 1.0
}

func (at *affineTransform) Components() (a, b, c, d, tx, ty float32) {
	return at.m[ma], at.m[mb], at.m[mc], at.m[md], at.m[me], at.m[mf]
}

func (at *affineTransform) SetByComp(a, b, c, d, tx, ty float32) {
	at.m[ma] = a
	at.m[mb] = b
	at.m[mc] = c
	at.m[md] = d
	at.m[me] = tx
	at.m[mf] = ty
}

func (at *affineTransform) SetByTransform(t api.IAffineTransform) {
	at.m[ma], at.m[mb], at.m[mc], at.m[md], at.m[me], at.m[mf] = t.Components()
}

func (at *affineTransform) TransformPoint(p api.IPoint) {
	p.SetByComp(
		(at.m[ma]*p.X())+(at.m[mc]*p.Y())+at.m[me],
		(at.m[mb]*p.X())+(at.m[md]*p.Y())+at.m[mf])
}

func (at *affineTransform) TransformToPoint(in api.IPoint, out api.IPoint) {
	out.SetByComp(
		(at.m[ma]*in.X())+(at.m[mc]*in.Y())+at.m[me],
		(at.m[mb]*in.X())+(at.m[md]*in.Y())+at.m[mf])
}

func (at *affineTransform) TransformCompToPoint(x, y float32, out api.IPoint) {
	out.SetByComp(
		(at.m[ma]*x)+(at.m[mc]*y)+at.m[me],
		(at.m[mb]*x)+(at.m[md]*y)+at.m[mf])
}

func (at *affineTransform) TransformToComps(in api.IPoint) (x, y float32) {
	return (at.m[ma] * in.X()) + (at.m[mc] * in.Y()) + at.m[me], (at.m[mb] * in.X()) + (at.m[md] * in.Y()) + at.m[mf]
}

func (at *affineTransform) Translate(x, y float32) {
	at.m[me] += (at.m[ma] * x) + (at.m[mc] * y)
	at.m[mf] += (at.m[mb] * x) + (at.m[md] * y)
}

func (at *affineTransform) MakeTranslate(x, y float32) {
	at.SetByComp(1.0, 0.0, 0.0, 1.0, x, y)
}

func (at *affineTransform) MakeTranslateUsingPoint(p api.IPoint) {
	at.SetByComp(1.0, 0.0, 0.0, 1.0, p.X(), p.Y())
}

func (at *affineTransform) Scale(sx, sy float32) {
	at.m[ma] *= sx
	at.m[mb] *= sx
	at.m[mc] *= sy
	at.m[md] *= sy
}

func (at *affineTransform) MakeScale(sx, sy float32) {
	at.SetByComp(sx, 0.0, 0.0, sy, 0.0, 0.0)
}

func (at *affineTransform) GetPsuedoScale() float32 {
	return at.m[ma]
}

// Concatinate a rotation (radians) onto this transform.
//
// Rotation is just a matter of perspective. A CW rotation can be seen as
// CCW depending on what you are talking about rotating. For example,
// if the coordinate system is thought as rotating CCW then objects are
// seen as rotating CW, and that is what the 2x2 matrix below represents.
//
// It is also the frame of reference we use. In this library +Y axis is downward
//     |cos  -sin|   object appears to rotate CW.
//     |sin   cos|
//
// In the matrix below the object appears to rotate CCW.
//     |cos  sin|
//     |-sin cos|
//
//     |a  c|    |cos  -sin|
//     |b  d|  x |sin   cos|
//
// If Y axis is downward (default for SDL and Image) then:
// +angle yields a CW rotation
// -angle yeilds a CCW rotation.
//
// else
// -angle yields a CW rotation
// +angle yeilds a CCW rotation.
func (at *affineTransform) Rotate(radians float64) {
	s := float32(math.Sin(radians))
	cr := float32(math.Cos(radians))
	a := at.m[ma]
	b := at.m[mb]
	c := at.m[mc]
	d := at.m[md]

	at.m[ma] = a*cr + c*s
	at.m[mb] = b*cr + d*s
	at.m[mc] = c*cr - a*s
	at.m[md] = d*cr - b*s
}

func (at *affineTransform) MakeRotate(radians float64) {
	s := float32(math.Sin(radians))
	c := float32(math.Cos(radians))
	at.m[ma] = c
	at.m[mb] = s
	at.m[mc] = -s
	at.m[md] = c
	at.m[me] = 0
	at.m[mf] = 0
}

// MultiplyPre performs: n = n * m
func MultiplyPre(m api.IAffineTransform, n api.IAffineTransform) {
	na, nb, nc, nd, _, _ := n.Components()
	ma, mb, mc, md, mtx, mty := m.Components()

	n.SetByComp(
		na*ma+nb*mc,
		na*mb+nb*md,
		nc*ma+nd*mc,
		nc*mb+nd*md,
		(na*mtx)+(nc*mty)+mtx,
		(nb*mtx)+(nd*mty)+mty)
}

// MultiplyPost performs: n = m * n
func MultiplyPost(m api.IAffineTransform, n api.IAffineTransform) {
	na, nb, nc, nd, ntx, nty := n.Components()
	ma, mb, mc, md, _, _ := m.Components()

	n.SetByComp(
		ma*na+mb*nc,
		ma*nb+mb*nd,
		mc*na+md*nc,
		mc*nb+md*nd,
		(ma*ntx)+(mc*nty)+ntx,
		(mb*ntx)+(md*nty)+nty)
}

// Multiply performs: out = m * n
func Multiply(m api.IAffineTransform, n api.IAffineTransform, out api.IAffineTransform) {
	na, nb, nc, nd, ntx, nty := n.Components()
	ma, mb, mc, md, mtx, mty := m.Components()

	out.SetByComp(
		ma*na+mb*nc,
		ma*nb+mb*nd,
		mc*na+md*nc,
		mc*nb+md*nd,
		(mtx*na)+(mty*nc)+ntx,
		(mtx*nb)+(mty*nd)+nty)
}

func (at *affineTransform) Invert() {
	a := at.m[ma]
	b := at.m[mb]
	c := at.m[mc]
	d := at.m[md]
	tx := at.m[me]
	ty := at.m[mf]

	determinant := 1.0 / (a*d - b*c)

	at.m[ma] = determinant * d
	at.m[mb] = -determinant * b
	at.m[mc] = -determinant * c
	at.m[md] = determinant * a
	at.m[me] = determinant * (c*ty - d*tx)
	at.m[mf] = determinant * (b*tx - a*ty)
}

func (at *affineTransform) InvertTo(out api.IAffineTransform) {
	determinant := 1.0 / (at.m[ma]*at.m[md] - at.m[mb]*at.m[mc])
	out.SetByComp(
		determinant*at.m[md],
		-determinant*at.m[mb],
		-determinant*at.m[mc],
		determinant*at.m[ma],
		determinant*(at.m[mc]*at.m[mf]-at.m[md]*at.m[me]),
		determinant*(at.m[mb]*at.m[me]-at.m[ma]*at.m[mf]))
}

func (at *affineTransform) Transpose() {
	c := at.m[mc]
	at.m[mc] = at.m[mb]
	at.m[mb] = c
}

func (at *affineTransform) Populate(destination api.IMatrix4) {
	d := destination.Matrix()
	m := at.m
	d[0] = m[0]
	d[1] = m[1]
	d[2] = m[2]
	d[3] = m[3]
	d[4] = m[4]
	d[5] = m[5]
	d[6] = m[6]
	d[7] = m[7]
	d[8] = m[8]
	d[9] = m[9]
	d[10] = m[10]
	d[11] = m[11]
	d[12] = m[12]
	d[13] = m[13]
	d[14] = m[14]
	d[15] = m[15]
}

func (at affineTransform) String() string {
	return fmt.Sprintf("|%7.3f,%7.3f,%7.3f|\n|%7.3f,%7.3f,%7.3f|",
		at.m[ma], at.m[mc], at.m[me],
		at.m[mb], at.m[md], at.m[mf])
}

// | ma mc m8  me  |
// | mb md m9  mf  |
// | m2 m6 m10 m14 |
// | m3 m7 m11 m15 |
func (at affineTransform) String4x4() string {
	return fmt.Sprintf("|%7.3f,%7.3f,%7.3f,%7.3f|\n|%7.3f,%7.3f,%7.3f,%7.3f|\n|%7.3f,%7.3f,%7.3f,%7.3f|\n|%7.3f,%7.3f,%7.3f,%7.3f|",
		at.m[ma], at.m[mc], at.m[m8], at.m[me],
		at.m[mb], at.m[md], at.m[m9], at.m[mf],
		at.m[m2], at.m[m6], at.m[m10], at.m[m14],
		at.m[m3], at.m[m7], at.m[m11], at.m[m15])
}
