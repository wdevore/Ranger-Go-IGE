package api

// IColor represents color types that support a []float32 array
type IColor interface {
	Color() []float32
	SetColor([]float32)
}
