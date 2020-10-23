package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type triPhysicsComponent struct {
	physicsComponent
}

func newTriPhysicsComponent() *triPhysicsComponent {
	o := new(triPhysicsComponent)
	return o
}

func (p *triPhysicsComponent) Build(phyWorld *box2d.B2World, node api.INode, position api.IPoint) {
	p.phyNode = node
	p.position = position

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

	tcc := p.phyNode.(*shapes.MonoTriangleNode)
	// Box2D expects polygon edges to be defined at full length, not
	// half-side
	scale := tcc.SideLength()
	verts := tcc.Vertices()

	vertices := []box2d.B2Vec2{}
	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scale), Y: float64((*verts)[i+1] * scale)})
	}
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[0] * scale), Y: float64(verts[1] * scale)})
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[3] * scale), Y: float64(verts[4] * scale)})
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[6] * scale), Y: float64(verts[7] * scale)})

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	b2Shape.Set(vertices, len(*verts)/api.XYZComponentCount)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}
