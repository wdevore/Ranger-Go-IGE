package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type boxPhysicsComponent struct {
	physicsComponent
}

func newBoxPhysicsComponent() *boxPhysicsComponent {
	o := new(boxPhysicsComponent)
	return o
}

// EnableGravity enables/disables gravity for this component
func (p *boxPhysicsComponent) EnableGravity(enable bool) {
	if enable {
		p.b2Body.SetGravityScale(1.0)
	} else {
		p.b2Body.SetGravityScale(0.0)
	}
}

// ApplyForce applies linear force to box center
func (p *boxPhysicsComponent) ApplyForce(dirX, dirY float64) {
	p.b2Body.ApplyForce(box2d.B2Vec2{X: dirX, Y: dirY}, p.b2Body.GetWorldCenter(), true)
}

// ApplyImpulse applies linear impulse to box center
func (p *boxPhysicsComponent) ApplyImpulse(dirX, dirY float64) {
	p.b2Body.ApplyLinearImpulse(box2d.B2Vec2{X: dirX, Y: dirY}, p.b2Body.GetWorldCenter(), true)
}

// ApplyImpulseToCorner applies linear impulse to 1,1 box corner
// As the box rotates the 1,1 corner rotates which means impulses
// could change the rotation to either CW or CCW.
func (p *boxPhysicsComponent) ApplyImpulseToCorner(dirX, dirY float64) {
	p.b2Body.ApplyLinearImpulse(box2d.B2Vec2{X: dirX, Y: dirY}, p.b2Body.GetWorldPoint(box2d.B2Vec2{X: 1.0, Y: 1.0}), true)
}

// ApplyTorque applies torgue to box center
func (p *boxPhysicsComponent) ApplyTorque(torgue float64) {
	p.b2Body.ApplyTorque(torgue, true)
}

// ApplyAngularImpulse applies angular impulse to box center
func (p *boxPhysicsComponent) ApplyAngularImpulse(impulse float64) {
	p.b2Body.ApplyAngularImpulse(impulse, true)
}

func (p *boxPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Set the position of the Body
	px := p.phyNode.Position().X()
	py := p.phyNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	p.b2Body = phyWorld.CreateBody(&bDef)

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()

	// Box2D assumes the same is defined in unit-space which
	// means if the object is defined otherwise we need the object
	// to return the correct value
	tcc := p.phyNode.(*extras.StaticSquareNode)
	b2Shape.SetAsBoxFromCenterAndAngle(
		float64(tcc.HalfSide()), float64(tcc.HalfSide()),
		box2d.B2Vec2{X: 0.0, Y: 0.0}, 0.0)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}
