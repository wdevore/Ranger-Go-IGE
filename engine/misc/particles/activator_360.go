package particles

import (
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Activator360 activates particles in random directions encompassing
// 360 degrees
type Activator360 struct {
	maxLife  float32
	maxSpeed float32
}

// NewActivator360 constructs an activator
func NewActivator360() api.IParticleActivator {
	o := new(Activator360)
	o.maxLife = api.MaxParticleLifetime
	o.maxSpeed = api.MaxParticleSpeed
	return o
}

// Activate configures a particle with a random direction and speed.
func (a *Activator360) Activate(particle api.IParticle, center api.IPoint) {
	direction := rand.Float64() * 360.0
	speed := float32(rand.Float64() * float64(a.maxSpeed))

	particle.SetVelocity(direction, speed)

	// The location of where the particle is emitted
	particle.SetPosition(center.X(), center.Y())

	// A random lifetime ranging from 0.0 to max_life
	lifespan := float32(rand.Float64()) * (a.maxLife * 1000.0)
	particle.SetLifespan(lifespan)

	color, isColorType := particle.Visual().(api.IColor)

	if isColorType {
		// Change the Red color component if the visual supports the IColor type
		// otherwise the color maintains its current color
		shade := float32(rand.Float64()*float64(speed)*50) / 255.0
		c := color.Color()
		c[0] = shade
		color.SetColor(c)
	}

	particle.Reset()

	particle.Activate(true)
}

// SetMaxLifetime sets maximum life a particle live
func (a *Activator360) SetMaxLifetime(duration float32) {
	a.maxLife = duration
}

// SetMaxSpeed set max magnitude
func (a *Activator360) SetMaxSpeed(speed float32) {
	a.maxSpeed = speed
}
