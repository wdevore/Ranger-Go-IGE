package particles

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Particle is the base object of a Particle system.
type Particle struct {
	elapsed  float32
	lifespan float32

	active bool
}

// NewParticle constructs a new Particle
func NewParticle() api.IParticle {
	o := new(Particle)

	o.active = false
	o.elapsed = 0.0
	o.lifespan = 0.0

	return o
}

// SetLifespan sets how long the Particle lives
func (p *Particle) SetLifespan(duration float32) {
	p.lifespan = duration
}

// Activate changes the Particle's state
func (p *Particle) Activate(active bool) {
	p.active = active
}

// IsActive indicates if the Particle is alive
func (p *Particle) IsActive() bool {
	return p.active
}

// Evaluate changes the Particle's state based on time
func (p *Particle) Evaluate(dt float32) {
	p.elapsed += dt

	p.active = p.elapsed < p.lifespan
}

// Reset resets the Particle
func (p *Particle) Reset() {
	p.active = false
	p.elapsed = 0.0
}
