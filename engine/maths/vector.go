package maths

import (
	"fmt"
	"math"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type vector struct {
	x, y float32
}

// NewVector constructs a new IVector
func NewVector() api.IVector {
	o := new(vector)
	return o
}

// NewVectorUsing constructs a new IVector using components
func NewVectorUsing(x, y float32) api.IVector {
	o := new(vector)
	o.x = x
	o.y = y
	return o
}

func (v *vector) Components() (x, y float32) {
	return v.x, v.y
}

func (v *vector) X() float32 {
	return v.x
}

func (v *vector) Y() float32 {
	return v.y
}

func (v *vector) SetByComp(x, y float32) {
	v.x = x
	v.y = y
}

func (v *vector) SetByPoint(ip api.IPoint) {
	v.x = ip.X()
	v.y = ip.Y()
}

func (v *vector) SetByAngle(radians float64) {
	v.x = float32(math.Cos(radians))
	v.y = float32(math.Sin(radians))
}

func (v *vector) SetByVector(ip api.IVector) {
	v.x = ip.X()
	v.y = ip.Y()
}

func (v *vector) Length() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}

func (v *vector) LengthSqr() float32 {
	return v.x*v.x + v.y*v.y
}

func (v *vector) Add(x, y float32) {
	v.x += x
	v.y += y
}

func (v *vector) Sub(x, y float32) {
	v.x -= x
	v.y -= y
}

func (v *vector) AddV(iv api.IVector) {
	v.x += iv.X()
	v.y += iv.Y()
}

func (v *vector) SubV(iv api.IVector) {
	v.x -= iv.X()
	v.y -= iv.Y()
}

// Add performs: out = v1 + v2
func Add(v1, v2, out api.IVector) {
	out.SetByComp(v1.X()+v2.X(), v1.Y()+v2.Y())
}

// Sub performs: out = v1 - v2
func Sub(v1, v2, out api.IVector) {
	out.SetByComp(v1.X()-v2.X(), v1.Y()-v2.Y())
}

func (v *vector) Scale(s float32) {
	v.x = v.x * s
	v.y = v.y * s
}

// ScaleBy performs: out = v * s
func ScaleBy(v api.IVector, s float32, out api.IVector) {
	out.SetByComp(v.X()*s, v.Y()*s)
}

func (v *vector) Div(d float32) {
	v.x = v.x / d
	v.y = v.y / d
}

var tmp = NewVector()

// Distance between two vectors
func VectorDistance(v1, v2 api.IVector) float32 {
	Sub(v1, v2, tmp)
	return tmp.Length()
}

func (v *vector) AngleX(vo api.IVector) float32 {
	return float32(math.Atan2(float64(vo.Y()), float64(vo.X())))
}

func (v *vector) Normalize() {
	len := v.Length()
	if len != 0.0 {
		v.Div(len)
	}
}

func (v *vector) SetDirection(radians float64) {
	v.SetByComp(float32(math.Cos(radians)), float32(math.Sin(radians)))
}

// Dot computes the dot-product between the vectors
func VectorDot(v1, v2 api.IVector) float32 {
	return v1.X()*v2.X() + v1.Y()*v2.Y()
}

// Cross computes the cross-product of two vectors
func Cross(v1, v2 api.IVector) float32 {
	return v1.X()*v2.Y() - v1.Y()*v2.X()
}

func (v *vector) CrossCW() {
	v.SetByComp(v.Y(), -v.X())
}

func (v *vector) CrossCCW() {
	v.SetByComp(-v.Y(), v.X())
}

var tmp2 = NewVector()

// Angle computes the angle in radians between two vector directions
func Angle(v1, v2 api.IVector) float32 {
	tmp.SetByVector(v1)
	tmp2.SetByVector(v2)

	tmp.Normalize()  // a2
	tmp2.Normalize() // b2

	angle := math.Atan2(float64(Cross(tmp, tmp2)), float64(VectorDot(tmp, tmp2)))

	if math.Abs(angle) < Epsilon {
		return 0.0
	}

	return float32(angle)
}

func (v vector) String() string {
	return fmt.Sprintf("<%0.3f,%0.3f>", v.x, v.y)
}
