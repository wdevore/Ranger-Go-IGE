package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type slopePhysicsComponent struct {
	physicsComponent

	slope api.INode
}

func newFencePhysicsComponent() *slopePhysicsComponent {
	o := new(slopePhysicsComponent)
	return o
}

func (p *slopePhysicsComponent) Build(world api.IWorld, parent api.INode, phyWorld *box2d.B2World, position api.IPoint, rotation float64) {
	p.buildNodes(world, parent, position, rotation)
	p.buildPhysics(phyWorld, position, rotation)
}

func (p *slopePhysicsComponent) buildPhysics(phyWorld *box2d.B2World, position api.IPoint, rotation float64) {
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
	bDef.Angle = p.slope.Rotation()

	// An instance of a body to contain Fixtures.
	// This body represents the entire fence (i.e. all four sides)
	p.b2Body = phyWorld.CreateBody(&bDef)

	fd := box2d.MakeB2FixtureDef()

	tln := p.slope.(*extras.StaticHLineNode)
	halfLength := float64(tln.HalfLength())

	// ------------------------------------------------------------
	// Top fixture
	// px := p.topLineNode.Position().X()
	// py := p.topLineNode.Position().Y()
	b2Shape := box2d.MakeB2EdgeShape()
	b2Shape.Set(
		box2d.MakeB2Vec2(-halfLength, 0.0),
		box2d.MakeB2Vec2(halfLength, 0.0))
	fd.Shape = &b2Shape
	// fmt.Println(p.topLineNode.Rotation())

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *slopePhysicsComponent) buildNodes(world api.IWorld, parent api.INode, position api.IPoint, rotation float64) error {
	var err error

	p.slope, err = extras.NewStaticHLineNode("Bottom", world, parent)
	if err != nil {
		return err
	}
	p.slope.SetScale(25.0)
	p.slope.SetPosition(position.X(), position.Y())
	p.slope.SetRotation(rotation)
	glh := p.slope.(*extras.StaticHLineNode)
	glh.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}
