package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type boxPhysicsComponent struct {
	physicsComponent

	beginContactColor api.IPalette
	endContactColor   api.IPalette

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

func newBoxPhysicsComponent() *boxPhysicsComponent {
	o := new(boxPhysicsComponent)
	return o
}

func (p *boxPhysicsComponent) MoveLeft() {
	velocity := p.b2Body.GetLinearVelocity()
	// velocity.X = math.Max(velocity.X-0.1, -5.0)
	velocity.X = velocity.X - 0.1
	p.b2Body.SetLinearVelocity(velocity)
}

func (p *boxPhysicsComponent) MoveRight() {
	velocity := p.b2Body.GetLinearVelocity()
	velocity.X = velocity.X + 0.1
	p.b2Body.SetLinearVelocity(velocity)
}

func (p *boxPhysicsComponent) MoveUp() {
	velocity := p.b2Body.GetLinearVelocity()
	velocity.Y = velocity.Y + 0.1
	p.b2Body.SetLinearVelocity(velocity)
}

func (p *boxPhysicsComponent) MoveDown() {
	velocity := p.b2Body.GetLinearVelocity()
	velocity.Y = velocity.Y - 0.1
	p.b2Body.SetLinearVelocity(velocity)
}

// EnableGravity enables/disables gravity for this component
func (p *boxPhysicsComponent) EnableGravity(enable bool) {
	if enable {
		p.b2Body.SetGravityScale(1.0)
	} else {
		p.b2Body.SetGravityScale(0.0)
	}
}

func (p *boxPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

	p.beginContactColor = color.NewPaletteInt64(color.LightPurple)
	p.endContactColor = color.NewPaletteInt64(color.Orange)

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

	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *boxPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

// ------------------------------------------------------
// Physics feedback
// ------------------------------------------------------

// HandleBeginContact processes BeginContact events
func (p *boxPhysicsComponent) HandleBeginContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*shapes.MonoSquareNode)

	if !ok {
		n, ok = nodeB.(*shapes.MonoSquareNode)
	}

	if ok {
		n.SetFilledColor(p.beginContactColor)
	}

	return false
}

// HandleEndContact processes EndContact events
func (p *boxPhysicsComponent) HandleEndContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*shapes.MonoSquareNode)

	if !ok {
		n, ok = nodeB.(*shapes.MonoSquareNode)
	}

	if ok {
		n.SetFilledColor(p.endContactColor)
	}

	return false
}
