package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type fencePhysicsComponent struct {
	physicsComponent

	bottomLineNode api.INode
	rightLineNode  api.INode
	topLineNode    api.INode
	leftLineNode   api.INode

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
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
	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	// ------------------------------------------------------------
	// Bottom fixture
	px = p.bottomLineNode.Position().X()
	py = p.bottomLineNode.Position().Y()
	tln := p.bottomLineNode.(*extras.StaticHLineNode)
	halfLength := float64(tln.HalfLength())

	b2Shape := box2d.MakeB2EdgeShape()
	// Positioning and size relative to the body
	b2Shape.Set(
		box2d.MakeB2Vec2(-halfLength, float64(py)),
		box2d.MakeB2Vec2(halfLength, float64(py)))
	fd.Shape = &b2Shape
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ------------------------------------------------------------
	// Right fixture
	px = p.rightLineNode.Position().X()
	py = p.rightLineNode.Position().Y()
	b2Shape = box2d.MakeB2EdgeShape()
	// Positioning and size relative to the body
	b2Shape.Set(
		box2d.MakeB2Vec2(float64(px), -halfLength),
		box2d.MakeB2Vec2(float64(px), halfLength))
	fd.Shape = &b2Shape
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ------------------------------------------------------------
	// Top fixture
	px = p.topLineNode.Position().X()
	py = p.topLineNode.Position().Y()
	b2Shape = box2d.MakeB2EdgeShape()
	b2Shape.Set(
		box2d.MakeB2Vec2(-halfLength, float64(py)),
		box2d.MakeB2Vec2(halfLength, float64(py)))
	fd.Shape = &b2Shape
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ------------------------------------------------------------
	// Left fixture
	px = p.leftLineNode.Position().X()
	py = p.leftLineNode.Position().Y()
	b2Shape = box2d.MakeB2EdgeShape()
	b2Shape.Set(
		box2d.MakeB2Vec2(float64(px), -halfLength),
		box2d.MakeB2Vec2(float64(px), halfLength))
	fd.Shape = &b2Shape
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *fencePhysicsComponent) buildNodes(world api.IWorld, parent api.INode) error {
	var err error

	scale := float32(75.0)

	p.bottomLineNode, err = extras.NewStaticHLineNode("Bottom", world, parent)
	if err != nil {
		return err
	}
	p.bottomLineNode.SetScale(scale * 2)
	p.bottomLineNode.SetPosition(0.0, -scale/2)
	glh := p.bottomLineNode.(*extras.StaticHLineNode)
	glh.SetColor(color.NewPaletteInt64(color.Yellow))

	p.rightLineNode, err = extras.NewStaticVLineNode("Right", world, parent)
	if err != nil {
		return err
	}
	p.rightLineNode.SetScale(scale)
	p.rightLineNode.SetPosition(scale, 0.0)
	glv := p.rightLineNode.(*extras.StaticVLineNode)
	glv.SetColor(color.NewPaletteInt64(color.Yellow))

	p.topLineNode, err = extras.NewStaticHLineNode("Top", world, parent)
	if err != nil {
		return err
	}
	p.topLineNode.SetScale(scale * 2)
	p.topLineNode.SetPosition(0.0, scale/2)
	glh = p.topLineNode.(*extras.StaticHLineNode)
	glh.SetColor(color.NewPaletteInt64(color.Yellow))

	p.leftLineNode, err = extras.NewStaticVLineNode("Left", world, parent)
	if err != nil {
		return err
	}
	p.leftLineNode.SetScale(scale)
	p.leftLineNode.SetPosition(-scale, 0.0)
	glv = p.leftLineNode.(*extras.StaticVLineNode)
	glv.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}

func (p *fencePhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}
