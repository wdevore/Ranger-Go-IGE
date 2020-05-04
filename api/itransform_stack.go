package api

// ITransformStack a transformation stack
type ITransformStack interface {
	Initialize(IMatrix4)
	Apply(IMatrix4) IMatrix4
	ApplyAffine(IAffineTransform) IMatrix4
	Save()
	Restore()
}
