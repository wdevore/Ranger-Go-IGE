package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type boxPhysicsComponent struct {
	physicsComponent
}

func newBoxPhysicsComponent() *boxPhysicsComponent {
	o := new(boxPhysicsComponent)
	return o
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
	tcc := p.phyNode.(*shapes.MonoSquareNode)
	b2Shape.SetAsBoxFromCenterAndAngle(
		float64(tcc.HalfSide()), float64(tcc.HalfSide()),
		box2d.B2Vec2{X: 0.0, Y: 0.0}, 0.0)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}
