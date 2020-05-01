package main

import (
	"fmt"
	"math"
	"testing"

	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

func TestRunner(t *testing.T) {
	matrix3(t)
}

func matrix1(t *testing.T) {
	am := maths.NewTransform()
	am.MakeTranslate(2.0, 2.0)
	am.Scale(2.0, 2.0)
	p := geometry.NewPointUsing(2.0, 2.0)

	am.TransformPoint(p)

	fmt.Println(p)

	m4 := maths.NewMatrix4()
	m4.SetTranslate3Comp(2.0, 2.0, 0.0)
	m4.ScaleByComp(2.0, 2.0, 1.0)
	fmt.Println(m4)
	v := maths.NewVector3With3Components(2.0, 2.0, 0.0)
	v.Mul(m4)
	fmt.Println(v)
}

func matrix2(t *testing.T) {
	om := maths.NewTransform()

	am := maths.NewTransform()
	am.MakeTranslate(1.0, 1.0)
	am.Rotate(math.Pi / 4.0)
	am.Scale(2.0, 2.0)

	bm := maths.NewTransform()
	bm.MakeTranslate(3.0, 3.0)
	bm.Scale(4.0, 4.0)

	maths.Multiply(am, bm, om)
	fmt.Println(om)

	p := geometry.NewPointUsing(2.0, 2.0)
	om.TransformPoint(p)

	fmt.Println(p)

	mo := maths.NewMatrix4()

	m4 := maths.NewMatrix4()
	m4.SetTranslate3Comp(1.0, 1.0, 0.0)
	m4.Rotate(math.Pi / 4.0)
	m4.ScaleByComp(2.0, 2.0, 1.0)

	m5 := maths.NewMatrix4()
	m5.SetTranslate3Comp(3.0, 3.0, 0.0)
	m5.ScaleByComp(4.0, 4.0, 1.0)

	maths.Multiply4(m4, m5, mo)
	fmt.Println(mo)

	v := maths.NewVector3With3Components(2.0, 2.0, 0.0)
	v.Mul(mo)
	fmt.Println(v)
}

func matrix3(t *testing.T) {
	// --------------------------------------------------
	// This section demonstrates using seperate matrices
	// to perform matrix concatination rather than using the
	// "appending" methods, for Rotate or Scale.
	// --------------------------------------------------
	// First we create each descrete matrix
	xt := maths.NewTransform()
	xt.MakeTranslate(1.0, 1.0)
	xs := maths.NewTransform()
	xs.MakeScale(2.0, 2.0)
	xr := maths.NewTransform()
	xr.MakeRotate(math.Pi / 4.0)

	// The output
	xo := maths.NewTransform()

	// A proper composite transform matrix sequence would be:
	// translate, concatinate rotation, then concatinate scale.
	// Using "post" multiply we would perform the sequence in the
	// opposite order:
	// scale x rotate x translate ==> xs * xr * xt
	// Thus translate is "post" multiplied last.
	maths.Multiply(xs, xr, xo)
	maths.Multiply(xo, xt, xo)
	// Which yields:
	// |  1.414, -1.414,  1.000|
	// |  1.414,  1.414,  1.000|
	fmt.Println("xo:")
	fmt.Println(xo)

	// ----------------------------------------------------------
	// This section uses the AffineTransform approach which has
	// concatinating methods, for example Rotate, verses
	// MakeRotate which overwrites the matrix.
	// ----------------------------------------------------------
	om := maths.NewTransform()
	am := maths.NewTransform()
	// Because the affinetransform's methods are "pre" multiply
	// types, we need do the opposite sequence from above:
	// translate x rotate x scale
	am.MakeTranslate(1.0, 1.0)
	am.Rotate(math.Pi / 4.0) // "pre" multiplies
	am.Scale(2.0, 2.0)       // "pre" multiplies
	fmt.Println("am:")
	fmt.Println(am)
	fmt.Println("am 4x4:")
	fmt.Println(am.String4x4())

	bm := maths.NewTransform()
	bm.MakeTranslate(3.0, 3.0)
	bm.Scale(4.0, 4.0)
	fmt.Println("bm:")
	fmt.Println(bm)

	maths.Multiply(am, bm, om)
	fmt.Println("om:")
	fmt.Println(om)

	// ----------------------------------------------------------
	// This section is an example of moving the affine transform
	// into a 4x4 matrix that can be used with OpenGL.
	// ----------------------------------------------------------
	m44 := maths.NewMatrix4()
	m44.SetFromAffine(om)
	fmt.Println("m44:")
	fmt.Println(m44)
	fmt.Println("--------------------------------------")
	fmt.Println("######################################")

	// ----------------------------------------------------------
	// This section does the same thing as the affine approach
	// above but using 4x4 matrices.
	// It does it 4 different ways using post and pre, and method
	// and functions.
	// ----------------------------------------------------------
	mt := maths.NewMatrix4()
	mt.SetTranslate3Comp(1.0, 1.0, 0.0)
	mr := maths.NewMatrix4()
	mr.SetRotation(math.Pi / 4.0)
	ms := maths.NewMatrix4()
	ms.SetScale3Comp(2.0, 2.0, 1.0)

	// First we use Multiply4 "pre" method which means we need
	// perform the sequence: translate x rotate x scale, which
	// is the reverse of the affine matrix above.
	mo := maths.NewMatrix4()
	// mt * mr * ms
	maths.Multiply4(mt, mr, mo)
	maths.Multiply4(mo, ms, mo)
	fmt.Println("(Mul4) mo:")
	fmt.Println(mo)

	// The Multiply "method" does the same thing requiring the
	// same sequence order.
	mo.ToIdentity()
	mo.Multiply(mt, mr)
	mo.Multiply(mo, ms)
	fmt.Println("(Mul) mo:")
	fmt.Println(mo)

	// The PostMultiply approach speaks for itself, so the
	// sequence is: scale x rotate x translate which is
	// the same as the affine approach above.
	mo.ToIdentity()
	mo.Set(ms)
	mo.PostMultiply(mr)
	mo.PostMultiply(mt)
	fmt.Println("(Pre) mo:")
	fmt.Println(mo)

	// And finally the PreMultiply approach
	// This is the reverse as the affine approach above
	mo.ToIdentity()
	mo.Set(mt)
	mo.PreMultiply(mr)
	mo.PreMultiply(ms)
	fmt.Println("(Post) mo:")
	fmt.Println(mo)

	// ----------------------------------------------------------
	// Now we use the "pre" concatinating methods that are functionally
	// the same as the affine concatinating methods but in 4x4 format.
	m4 := maths.NewMatrix4()
	m4.SetTranslate3Comp(1.0, 1.0, 0.0)
	m4.Rotate(math.Pi / 4.0)
	m4.ScaleByComp(2.0, 2.0, 1.0)
	fmt.Println("m4:")
	fmt.Println(m4)

	m5 := maths.NewMatrix4()
	m5.SetTranslate3Comp(3.0, 3.0, 0.0)
	m5.ScaleByComp(4.0, 4.0, 1.0)
	fmt.Println("m5:")
	fmt.Println(m5)

	mo.ToIdentity()
	maths.Multiply4(m5, m4, mo)
	fmt.Println("mo:")
	fmt.Println(mo)

	if !mo.Eq(m44) {
		t.Fatal("Expecting mo == m44")
	}
}
