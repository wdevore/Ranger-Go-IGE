package api

// IDelay scene timing and transitions
type IDelay interface {
	Reset()

	SetPauseTime(milliseconds float64)
	Inc(dt float64)
	UpdateTransition(dt float64)

	ReadyToTransition() bool
}
