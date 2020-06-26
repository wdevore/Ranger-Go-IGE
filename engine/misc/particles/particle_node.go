package particles

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// NodeParticle is the base object of a NodeParticle system.
type NodeParticle struct {
	elapsed  float32
	lifespan float32
	position api.IPoint
	velocity api.IVelocity

	active bool

	// Visual representation
	node api.INode
}

// NewNodeParticle constructs a new ParticleNode
func NewNodeParticle(visual api.INode) api.IParticle {
	o := new(NodeParticle)

	o.active = false
	o.elapsed = 0.0
	o.lifespan = 0.0
	o.velocity = maths.NewVelocity()
	o.position = geometry.NewPoint()

	o.node = visual

	return o
}

// SetPosition sets ParticleNodes initial position
func (p *NodeParticle) SetPosition(x, y float32) {
	p.position.SetByComp(x, y)
	p.node.SetPosition(x, y)
}

// GetPosition gets the ParticleNode's current position
func (p *NodeParticle) GetPosition() api.IPoint {
	return p.position
}

// SetLifespan sets how long the ParticleNode lives
func (p *NodeParticle) SetLifespan(duration float32) {
	p.lifespan = duration
}

// Visual gets the current INode assigned to this ParticleNode
func (p *NodeParticle) Visual() api.INode {
	return p.node
}

// Activate changes the ParticleNode's state
func (p *NodeParticle) Activate(active bool) {
	p.active = active
	p.node.SetVisible(active)
}

// IsActive indicates if the ParticleNode is alive
func (p *NodeParticle) IsActive() bool {
	return p.active
}

// SetVelocity changes the velocity
func (p *NodeParticle) SetVelocity(angle float64, speed float32) {
	p.velocity.SetDirectionByAngle(angle)
	p.velocity.SetMagnitude(speed)
}

// Evaluate changes the ParticleNode's state based on time
func (p *NodeParticle) Evaluate(dt float32) {
	p.elapsed += dt

	p.active = p.elapsed < p.lifespan

	// Update ParticleNode's position as long as the ParticleNode is active.
	if p.active {
		p.velocity.ApplyToPoint(p.position)
		p.node.SetPosition(p.position.X(), p.position.Y())
	}
}

// Reset resets the ParticleNode
func (p *NodeParticle) Reset() {
	p.active = false
	p.elapsed = 0.0
	p.node.SetVisible(p.active)
}
