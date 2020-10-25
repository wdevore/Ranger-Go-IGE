package main

import (
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// ActivatorCone activates particles in random directions encompassing
// fixed range of angles
type ActivatorCone struct {
	minLife            float64
	maxLife            float64
	startAngle         float64
	endAngle           float64
	minAngularVelocity float64
	maxAngularVelocity float64
	minForce           float64
	maxForce           float64
}

// NewActivatorCone constructs an activator
func NewActivatorCone() api.IParticleActivator {
	o := new(ActivatorCone)
	o.minLife = 2000.0
	o.maxLife = 15000.0
	o.minAngularVelocity = 5.0
	o.maxAngularVelocity = 15.0
	o.minForce = 15.0
	o.maxForce = 30.0
	o.startAngle = maths.DegreeToRadians * 85.0
	o.endAngle = maths.DegreeToRadians * 95.0
	return o
}

// Activate configures a particle with a random direction and speed.
func (a *ActivatorCone) Activate(particle api.IParticle, center api.IPoint) {
	direction := maths.Lerp(a.startAngle, a.endAngle, rand.Float64())

	angularVel := maths.Lerp(a.minAngularVelocity, a.maxAngularVelocity, rand.Float64())
	if rand.Float64() < 0.5 {
		angularVel = -angularVel
	}

	upForce := maths.Lerp(a.minForce, a.maxForce, rand.Float64())

	triPhy := particle.(*triPhysicsComponent)
	triPhy.ParticleConfigure(direction, angularVel, upForce)

	// A random lifetime ranging from 0.0 to max_life
	lifespan := maths.Lerp(a.minLife, a.maxLife, rand.Float64())
	particle.SetLifespan(float32(lifespan))

	// fmt.Println("dir: ", direction, ", ang: ", angularVel, ", force: ", upForce, ", life: ", lifespan)

	particle.Reset()

	particle.Activate(true)
}

// SetMaxLifetime sets maximum life a particle live
func (a *ActivatorCone) SetMaxLifetime(duration float64) {
	a.maxLife = duration
}

// SetStartAngle set cone spread angle
func (a *ActivatorCone) SetStartAngle(angle float64) {
	a.startAngle = angle
}

// SetEndAngle set cone spread angle
func (a *ActivatorCone) SetEndAngle(angle float64) {
	a.endAngle = angle
}
