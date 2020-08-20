package particles

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
)

// ParticleSystem is a collection of particles
type ParticleSystem struct {
	epiCenter api.IPoint

	activator api.IParticleActivator
	particles []api.IParticle

	active      bool
	autoTrigger bool
}

// NewParticleSystem constructs a new system
func NewParticleSystem(activator api.IParticleActivator) api.IParticleSystem {
	o := new(ParticleSystem)
	o.activator = activator
	o.epiCenter = geometry.NewPoint()
	o.particles = []api.IParticle{}
	o.active = false
	o.autoTrigger = false
	return o
}

// AddParticle adds a particle to system
func (ps *ParticleSystem) AddParticle(particle api.IParticle) {
	ps.particles = append(ps.particles, particle)
}

// SetPosition set the epi center of the system, where the particles
// emit from.
func (ps *ParticleSystem) SetPosition(x, y float32) {
	ps.epiCenter.SetByComp(x, y)
}

// SetAutoTrigger enables/disable autotrigger
func (ps *ParticleSystem) SetAutoTrigger(enable bool) {
	ps.autoTrigger = enable
}

// Activate enables/disable system
func (ps *ParticleSystem) Activate(enable bool) {
	ps.active = enable
}

// Update updates all active particles
func (ps *ParticleSystem) Update(dt float32) {
	if ps.active {
		ps.active = false
		for _, p := range ps.particles {
			if !p.IsActive() {
				if ps.autoTrigger {
					ps.TriggerOneshot()
				} else {
					p.Activate(false) // deactivate particle
				}
			} else {
				p.Evaluate(dt)
				ps.active = true // Keep ps alive until all died
			}
		}
	}
}

// TriggerOneshot activates a single particle
func (ps *ParticleSystem) TriggerOneshot() {
	ps.active = true
	// Look for a dead particle to resurrect.
	for _, p := range ps.particles {
		if !p.IsActive() {
			ps.activator.Activate(p, ps.epiCenter)
			break
		}
	}
}

// TriggerAt activates a single particle at a specific position
func (ps *ParticleSystem) TriggerAt(pos api.IPoint) {
	// Look for a dead particle to resurrect.
	for _, p := range ps.particles {
		if !p.IsActive() {
			ps.activator.Activate(p, pos)
			break
		}
	}
}

// TriggerExplosion activates the entire system at once
func (ps *ParticleSystem) TriggerExplosion() {
	// Look for a dead particle to resurrect.
	for _, p := range ps.particles {
		ps.activator.Activate(p, ps.epiCenter)
	}
}
