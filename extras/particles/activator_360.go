package particles

import (
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

// Activator360 activates particles in random directions encompassing
// 360 degrees
type Activator360 struct {
	maxLife  float64
	maxSpeed float64
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

	pNode := particle.(api.IParticleNode)
	pNode.SetVelocity(direction, speed)

	// The location of where the particle is emitted
	pNode.SetPosition(center.X(), center.Y())

	color, isColorType := pNode.Visual().(api.IColor)

	if isColorType {
		// Change the Red color component if the visual supports the IColor type
		// otherwise the color maintains its current color
		shade := float32(rand.Float64()*float64(speed)*50) / 255.0
		c := color.Color()
		c[0] = shade
		color.SetColor(c)
	}

	// A random lifetime ranging from 0.0 to max_life
	lifespan := rand.Float64() * (a.maxLife * 1000.0)
	particle.SetLifespan(float32(lifespan))

	particle.Reset()

	particle.Activate(true)
}

// SetMaxLifetime sets maximum life a particle live
func (a *Activator360) SetMaxLifetime(duration float64) {
	a.maxLife = duration
}
