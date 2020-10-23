package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type arcPhysicsComponent struct {
	physicsComponent

	beginContactColor api.IPalette
	endContactColor   api.IPalette

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

func newArcPhysicsComponent() *arcPhysicsComponent {
	o := new(arcPhysicsComponent)
	return o
}

func (p *arcPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

	p.beginContactColor = color.NewPaletteInt64WithAlpha(color.DarkGray, 0.75)
	p.endContactColor = color.NewPaletteInt64WithAlpha(color.LightPurple, 0.75)

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

	p.b2Body.SetAngularVelocity(45.0 * maths.DegreeToRadians)

	tcc := p.phyNode.(*shapes.MonoArcNode)
	scale := p.phyNode.Scale()
	verts := tcc.Vertices()

	vertices := []box2d.B2Vec2{}
	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scale), Y: float64((*verts)[i+1] * scale)})
	}

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	// fd.Density = 1.0

	// Remember to set the UserData so the ContactListner can call the node's
	// Begin/End contact methods.
	fd.UserData = p.phyNode
	fd.IsSensor = true

	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *arcPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

// ------------------------------------------------------
// Physics feedback
// ------------------------------------------------------

// HandleBeginContact processes BeginContact events
func (p *arcPhysicsComponent) HandleBeginContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*extras.StaticArcNode)

	if !ok {
		n, ok = nodeB.(*extras.StaticArcNode)
	}

	if ok {
		n.SetColor(p.beginContactColor)
	}

	return false
}

// HandleEndContact processes EndContact events
func (p *arcPhysicsComponent) HandleEndContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*extras.StaticArcNode)

	if !ok {
		n, ok = nodeB.(*extras.StaticArcNode)
	}

	if ok {
		n.SetColor(p.endContactColor)
	}

	return false
}
