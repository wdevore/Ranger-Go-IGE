package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type cirPhysicsComponent struct {
	physicsComponent

	beginContactColor api.IPalette
	endContactColor   api.IPalette

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

func newCirPhysicsComponent() *cirPhysicsComponent {
	o := new(cirPhysicsComponent)
	return o
}

func (p *cirPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

	p.beginContactColor = color.NewPaletteInt64(color.LightPurple)
	p.endContactColor = color.NewPaletteInt64(color.Orange)

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
	tcc := p.phyNode.(*extras.StaticCircleNode)
	radius := tcc.Radius()
	circleShape.M_radius = float64(radius)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &circleShape
	fd.Density = 1.0

	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *cirPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

// ------------------------------------------------------
// Physics feedback
// ------------------------------------------------------

// HandleBeginContact processes BeginContact events
func (p *cirPhysicsComponent) HandleBeginContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*extras.StaticCircleNode)

	if !ok {
		n, ok = nodeB.(*extras.StaticCircleNode)
	}

	if ok {
		n.SetColor(p.beginContactColor)
	}

	return false
}

// HandleEndContact processes EndContact events
func (p *cirPhysicsComponent) HandleEndContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*extras.StaticCircleNode)

	if !ok {
		n, ok = nodeB.(*extras.StaticCircleNode)
	}

	if ok {
		n.SetColor(p.endContactColor)
	}

	return false
}
