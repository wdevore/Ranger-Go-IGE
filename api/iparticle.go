package api

// IParticle represents a particle
type IParticle interface {
	Update(dt float32)

	IsActive() bool
	Activate(bool)
	Reset()

	SetPosition(x, y float32)
	GetPosition() IPoint
	SetLifespan(duration float32)

	SetVelocity(angle float64, speed float32)

	Visual() INode
}
