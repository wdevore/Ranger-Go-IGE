package api

// ITextureRenderer renders a single texture
type ITextureRenderer interface {
	Build(atlasName string)
	Draw(model IMatrix4)
	SetColor(color []float32)
	SelectCoordsByIndex(index int)
	Use()
	UnUse()
}
