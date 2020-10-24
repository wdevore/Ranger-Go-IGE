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
	tln := p.bottomLineNode.(*shapes.MonoHLineNode)
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

	scaleX := float32(130.0)
	scaleY := float32(75.0)

	p.bottomLineNode, err = shapes.NewMonoHLineNode("Bottom", world, parent)
	if err != nil {
		return err
	}
	p.bottomLineNode.SetScale(scaleX)
	p.bottomLineNode.SetPosition(0.0, -scaleY/2)
	ghl := p.bottomLineNode.(*shapes.MonoHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.Yellow))

	p.rightLineNode, err = shapes.NewMonoVLineNode("Right", world, parent)
	if err != nil {
		return err
	}
	p.rightLineNode.SetScale(scaleY)
	p.rightLineNode.SetPosition(scaleX/2, 0.0)
	ghv := p.rightLineNode.(*shapes.MonoVLineNode)
	ghv.SetColor(color.NewPaletteInt64(color.Yellow))

	p.topLineNode, err = shapes.NewMonoHLineNode("Top", world, parent)
	if err != nil {
		return err
	}
	p.topLineNode.SetScale(scaleX)
	p.topLineNode.SetPosition(0.0, scaleY/2)
	ghl = p.topLineNode.(*shapes.MonoHLineNode)
	ghl.SetColor(color.NewPaletteInt64(color.Yellow))

	p.leftLineNode, err = shapes.NewMonoVLineNode("Left", world, parent)
	if err != nil {
		return err
	}
	p.leftLineNode.SetScale(scaleY)
	p.leftLineNode.SetPosition(-scaleX/2, 0.0)
	ghv = p.leftLineNode.(*shapes.MonoVLineNode)
	ghv.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}

func (p *fencePhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}
