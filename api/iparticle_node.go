package api

// IParticleNode represents a particle
type IParticleNode interface {
	SetPosition(x, y float32)
	GetPosition() IPoint
	SetVelocity(angle float64, speed float32)
	Visual() INode
}
