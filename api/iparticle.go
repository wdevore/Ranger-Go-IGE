package api

// IParticle represents a baseline particle
type IParticle interface {
	Evaluate(dt float32)

	IsActive() bool
	Activate(bool)
	Reset()

	SetLifespan(duration float32)
}
