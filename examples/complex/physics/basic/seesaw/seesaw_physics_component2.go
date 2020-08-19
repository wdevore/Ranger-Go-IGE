package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

func (p *seesawPhysicsComponent) Build2(world api.IWorld, parent api.INode, phyWorld *box2d.B2World, position api.IPoint) {
	p.position = position

	p.buildPolygon2(world, parent)
	p.buildCircle2(world, p.phyNode)
	p.buildSquare2(world, p.phyNode)

	p.buildPhysics2(phyWorld, position)
}

func (p *seesawPhysicsComponent) buildPhysics2(phyWorld *box2d.B2World, position api.IPoint) {

	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Set the position of the Body
	px := position.X()
	py := position.Y()
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
	tcc := p.phyNode.(*custom.StaticPolygonNode)

	vertices := []box2d.B2Vec2{}
	verts := tcc.Vertices()
	// s := p.phyNode.Scale()

	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[0]), Y: float64((*verts)[1])})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[3]), Y: float64((*verts)[4])})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[6]), Y: float64((*verts)[7])})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[9]), Y: float64((*verts)[10])})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[12]), Y: float64((*verts)[13])})

	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ---------------------------------------------------------------
	// Circle
	// ---------------------------------------------------------------
	// Every Fixture has a shape
	b2CircleShape := box2d.MakeB2CircleShape()
	b2CircleShape.M_p.Set(float64(p.circle.Position().X()), float64(p.circle.Position().Y())) // Relative to body position
	gcir := p.circle.(*custom.StaticCircleNode)
	b2CircleShape.SetRadius(float64(gcir.Radius()))

	fd = box2d.MakeB2FixtureDef()
	fd.Shape = &b2CircleShape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ---------------------------------------------------------------
	// Square
	// ---------------------------------------------------------------
	// Every Fixture has a shape
	b2SquareShape := box2d.MakeB2PolygonShape()
	gss := p.square.(*custom.StaticSquareNode)

	b2SquareShape.SetAsBoxFromCenterAndAngle(
		float64(gss.HalfSide()), float64(gss.HalfSide()),
		box2d.B2Vec2{X: float64(p.square.Position().X()), Y: float64(p.square.Position().Y())}, 0.0)

	fd = box2d.MakeB2FixtureDef()
	fd.Shape = &b2SquareShape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *seesawPhysicsComponent) buildPolygon2(world api.IWorld, parent api.INode) error {
	var err error

	// --------------------------------------------------------------
	p.phyNode, err = custom.NewStaticPolygonNode("Polygon", false, world, parent)
	if err != nil {
		return err
	}

	scale := float32(3.0)
	// This build version embeds the scale directly inside
	// the vertices below rather than scaling the node
	// p.phyNode.SetScale(scale)

	p.phyNode.SetPosition(p.position.X(), p.position.Y())
	gpol := p.phyNode.(*custom.StaticPolygonNode)
	gpol.SetColor(color.NewPaletteInt64(color.LightOrange))

	vertices := []float32{
		-1.0, 2.0, 0.0,
		-1.0, 0.0, 0.0,
		0.0, -3.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 1.0, 0.0,
	}

	for i, v := range vertices {
		vertices[i] = v * scale
	}

	indices := []uint32{
		0, 1, 2, 3, 4,
	}

	gpol.Populate(vertices, indices)

	return nil
}

func (p *seesawPhysicsComponent) buildCircle2(world api.IWorld, parent api.INode) error {
	var err error

	p.circle, err = custom.NewStaticCircleNode("Circle", false, world, parent)
	if err != nil {
		return err
	}
	p.circle.SetScale(5.0)
	p.circle.SetPosition(-15.0, 0.0)
	gol2 := p.circle.(*custom.StaticCircleNode)
	gol2.SetColor(color.NewPaletteInt64(color.Green))

	return nil
}

func (p *seesawPhysicsComponent) buildSquare2(world api.IWorld, parent api.INode) error {
	var err error

	p.square, err = custom.NewStaticSquareNode("Square", true, false, world, parent)
	if err != nil {
		return err
	}
	p.square.SetScale(5.0)
	p.square.SetPosition(15.0, 0.0)
	gol2 := p.square.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.Aqua))

	return nil
}
