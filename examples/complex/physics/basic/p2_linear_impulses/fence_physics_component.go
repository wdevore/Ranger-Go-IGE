package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type fencePhysicsComponent struct {
	physicsComponent

	bottomLineNode api.INode
}

func newFencePhysicsComponent() *fencePhysicsComponent {
	o := new(fencePhysicsComponent)
	return o
}

func (p *fencePhysicsComponent) Build(world api.IWorld, parent api.INode, phyWorld *box2d.B2World, position api.IPoint) {
	p.buildNodes(world, parent)
	p.buildPhysics(phyWorld, position)
}

func (p *fencePhysicsComponent) buildPhysics(phyWorld *box2d.B2World, position api.IPoint) {
	p.position = position

	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_staticBody

	// Set the position of the Body
	px := position.X()
	py := position.Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures.
	// This body represents the entire fence (i.e. all four sides)
	p.b2Body = phyWorld.CreateBody(&bDef)

	fd := box2d.MakeB2FixtureDef()

	// ------------------------------------------------------------
	// Bottom fixture
	px = p.bottomLineNode.Position().X()
	py = p.bottomLineNode.Position().Y()
	tln := p.bottomLineNode.(*shapes.MonoHLineNode)
	halfLength := float64(tln.HalfLength())

	b2Shape := box2d.MakeB2EdgeShape()
	// Positioning and size relative to the body
	b2Shape.Set(
		box2d.MakeB2Vec2(-halfLength, float64(py)),
		box2d.MakeB2Vec2(halfLength, float64(py)))
	fd.Shape = &b2Shape
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *fencePhysicsComponent) buildNodes(world api.IWorld, parent api.INode) error {
	var err error

	p.bottomLineNode, err = shapes.NewMonoHLineNode("Ground", world, parent)
	if err != nil {
		return err
	}
	p.bottomLineNode.SetScale(100.0)
	p.bottomLineNode.SetPosition(0.0, -12.5)
	ghl := p.bottomLineNode.(*shapes.MonoHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}
