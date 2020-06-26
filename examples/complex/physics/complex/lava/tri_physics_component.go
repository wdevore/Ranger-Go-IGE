package main

import (
	"fmt"
	"math"

	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/misc/particles"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
)

// Lava particle
type triPhysicsComponent struct {
	particles.Particle
	physicsComponent

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...

	initialPosition api.IPoint
	direction       float64
	angularVel      float64
	force           float64

	debug int
}

func newTriPhysicsComponent() *triPhysicsComponent {
	o := new(triPhysicsComponent)
	return o
}

func (p *triPhysicsComponent) Configure(
	phyWorld *box2d.B2World, parent api.INode,
	id int, scale float32, color api.IPalette,
	position api.IPoint) {
	var err error

	p.initialPosition = position

	p.phyNode, err = custom.NewStaticTriangleNode(fmt.Sprintf("::Lava%d", id), true, true, parent.World(), parent)
	if err != nil {
		panic(err)
	}
	p.phyNode.SetVisible(false)

	p.phyNode.SetScale(scale)
	p.phyNode.SetPosition(position.X(), position.Y())
	gp := p.phyNode.(*custom.StaticTriangleNode)
	gp.SetColor(color)

	p.Build(phyWorld, position)
}

func (p *triPhysicsComponent) Build(phyWorld *box2d.B2World, position api.IPoint) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Set the initial position of the Body
	px := p.phyNode.Position().X()
	py := p.phyNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	p.b2Body = phyWorld.CreateBody(&bDef)

	tcc := p.phyNode.(*custom.StaticTriangleNode)
	// Box2D expects polygon edges to be defined at full length, not
	// half-side
	scale := tcc.SideLength()
	verts := tcc.Vertices()

	vertices := []box2d.B2Vec2{}
	for i := 0; i < len(verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64(verts[i] * scale), Y: float64(verts[i+1] * scale)})
	}

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0

	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	fd.Friction = 0.1
	fd.Restitution = 0.4

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *triPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

// EnableGravity enables/disables gravity for this component
func (p *triPhysicsComponent) EnableGravity(enable bool) {
	if enable {
		p.b2Body.SetGravityScale(1.0)
	} else {
		p.b2Body.SetGravityScale(0.0)
	}
}

func (p *triPhysicsComponent) Update() {
	if p.b2Body.IsActive() {
		pos := p.b2Body.GetPosition()
		p.phyNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot := p.b2Body.GetAngle()
		p.phyNode.SetRotation(rot)
	}
}

//-----------------------------------------------------------------
// Particles
//-----------------------------------------------------------------

// SetLifespan sets how long the Particle lives
func (p *triPhysicsComponent) SetLifespan(duration float32) {
	p.Particle.SetLifespan(duration)
}

// Activate changes the Particle's state
func (p *triPhysicsComponent) Activate(active bool) {
	p.Particle.Activate(active)

	// Note: All physic properties must be reset otherwise they
	// carry over into the next activation causing compounding effects.
	// So even though only an Angular velocity was applied, during the
	// simulation linear velocities manifest during collisions, we need
	// to clear any linear effects too.
	p.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))

	p.b2Body.SetTransform(
		box2d.MakeB2Vec2(float64(p.initialPosition.X()), float64(p.initialPosition.Y())),
		0.0)

	p.phyNode.SetVisible(active)

	if !active {
		return
	}

	p.b2Body.SetAngularVelocity(p.angularVel)

	// Apply force upwards
	p.b2Body.ApplyLinearImpulse(
		box2d.B2Vec2{X: math.Cos(p.direction) * p.force, Y: math.Sin(p.direction) * p.force},
		p.b2Body.GetWorldCenter(), true)
}

// IsActive indicates if the Particle is alive
func (p *triPhysicsComponent) IsActive() bool {
	return p.Particle.IsActive()
}

// Update changes the Particle's state based on time
func (p *triPhysicsComponent) Evaluate(dt float32) {
	p.Particle.Evaluate(dt)

	// Update Particle's position as long as the Particle is active.
	if !p.Particle.IsActive() {
		p.phyNode.SetVisible(false)
	}
}

// Reset resets the Particle
func (p *triPhysicsComponent) Reset() {
	p.Particle.Reset()
	p.phyNode.SetVisible(false)
}

// ParticleConfigure configures a particle prior to activation.
func (p *triPhysicsComponent) ParticleConfigure(direction, angularVel, force float64) {
	p.direction = direction
	p.angularVel = angularVel
	p.force = force
}
