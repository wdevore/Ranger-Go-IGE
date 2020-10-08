package api

// IDynamicRenderer renders a dynamic object
type IDynamicRenderer interface {
	Build(atlasName string)
	Draw(model IMatrix4)
	SetColor(color []float32)
	Update()

	Use()
	UnUse()
}
