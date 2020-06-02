package api

// IParticleSystem represents a particle system
type IParticleSystem interface {
	AddParticle(particle IParticle)
	Update(dt float32)

	SetPosition(x, y float32)
	SetAutoTrigger(bool)
	Activate(bool)

	TriggerOneshot()
	TriggerAt(pos IPoint)
	TriggerExplosion()
}
