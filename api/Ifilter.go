package api

// IFilter represents Transform Filter nodes
type IFilter interface {
	Build(IWorld)
	Visit(transStack ITransformStack, interpolation float64)

	InheritOnlyRotation()
	InheritOnlyScale()
	InheritOnlyTranslation()
	InheritRotationAndTranslation()
	InheritAll()
}
