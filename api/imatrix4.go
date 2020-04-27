package api

// IMatrix4 is a 4x4 matrix
type IMatrix4 interface {
	TranslateBy(v IVector3)
	TranslateBy3Comps(x, y, z float32)
	TranslateBy2Comps(x, y float32)
	SetTranslateByVector(v IVector3)
	SetTranslate3Comp(x, y, z float32)
	GetTranslation(out IVector3)
	SetRotation(angle float32)
	RotateBy(angle float32)
	ScaleBy(v IVector3)
	SetScale(v IVector3)
	SetScale3Comp(sx, sy, sz float32)
	SetScale2Comp(sx, sy float32)
	GetScale(out IVector3)
	PostScale(sx, sy, sz float32)
	Multiply(a, b IMatrix4)
	PreMultiply(b IMatrix4)
	PostMultiply(b IMatrix4)
	PostTranslate(tx, ty, tz float32)
	SetToOrtho(left, right, bottom, top, near, far float32)
	C(i int) float32
	E() [16]float32
	Clone() IMatrix4
	Set(src IMatrix4)
	ToIdentity()
}
