package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
)

type cirPhysicsComponent struct {
	physicsComponent
}

func newCirPhysicsComponent() *cirPhysicsComponent {
	o := new(cirPhysicsComponent)
	return o
}

func (p *cirPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	px := p.phyNode.Position().X()
	py := p.phyNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	p.b2Body = phyWorld.CreateBody(&bDef)

	// Every Fixture has a shape
	circleShape := box2d.MakeB2CircleShape()
	circleShape.M_p.Set(0.0, 0.0) // Relative to body position
	tcc := p.phyNode.(*custom.StaticCircleNode)
	radius := tcc.Radius()
	circleShape.M_radius = float64(radius)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &circleShape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}
