package api

// ITransformStack a transformation stack
type ITransformStack interface {
	Initialize()
	Apply(aft IAffineTransform)
	Save()
	Restore()
}
