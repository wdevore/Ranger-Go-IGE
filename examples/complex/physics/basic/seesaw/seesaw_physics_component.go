package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type seesawPhysicsComponent struct {
	physicsComponent

	circle api.INode
	square api.INode

	// Var for variant #3
	squareShape box2d.B2PolygonShape
	circleShape box2d.B2CircleShape
}

func newSeesawPhysicsComponent() *seesawPhysicsComponent {
	o := new(seesawPhysicsComponent)
	return o
}

func (p *seesawPhysicsComponent) Update(msPerUpdate, secPerUpdate float64) {
	// for fix := p.b2Body.GetFixtureList(); fix != nil; fix = fix.GetNext() {
	// 	shape := fix.GetShape()
	// 	if circle, ok := shape.(box2d.B2CircleShape); ok {
	// 		pos := circle.M_p;
	// 	}
	// }

	if p.b2Body.IsActive() {
		pos := p.b2Body.GetPosition()
		p.phyNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot := p.b2Body.GetAngle()
		p.phyNode.SetRotation(rot)
	}
}

func (p *seesawPhysicsComponent) Reset() {
	x := p.position.X()
	y := p.position.Y()

	p.phyNode.SetPosition(float32(x), float32(y))
	p.b2Body.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
	p.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	p.b2Body.SetAngularVelocity(0.0)
	p.b2Body.SetAwake(true)
}

func (p *seesawPhysicsComponent) Build(world api.IWorld, parent api.INode, phyWorld *box2d.B2World, position api.IPoint) {
	p.position = position

	p.buildPolygon(world, parent)
	p.buildCircle(world, p.phyNode)
	p.buildSquare(world, p.phyNode)

	p.buildPhysics(phyWorld, position)
}

func (p *seesawPhysicsComponent) buildPhysics(phyWorld *box2d.B2World, position api.IPoint) {

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
	tcc := p.phyNode.(*extras.StaticPolygonNode)

	vertices := []box2d.B2Vec2{}
	verts := tcc.Vertices()
	s := p.phyNode.Scale()

	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[0] * s), Y: float64((*verts)[1] * s)})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[3] * s), Y: float64((*verts)[4] * s)})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[6] * s), Y: float64((*verts)[7] * s)})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[9] * s), Y: float64((*verts)[10] * s)})
	vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[12] * s), Y: float64((*verts)[13] * s)})

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
	b2CircleShape.M_p.Set(float64(p.circle.Position().X()*s), float64(p.circle.Position().Y()*s)) // Relative to body position
	gcir := p.circle.(*extras.StaticCircleNode)
	b2CircleShape.SetRadius(float64(gcir.Radius() * s))

	fd = box2d.MakeB2FixtureDef()
	fd.Shape = &b2CircleShape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	// ---------------------------------------------------------------
	// Square
	// ---------------------------------------------------------------
	// Every Fixture has a shape
	b2SquareShape := box2d.MakeB2PolygonShape()
	gss := p.square.(*extras.StaticSquareNode)

	b2SquareShape.SetAsBoxFromCenterAndAngle(
		float64(gss.HalfSide()*s), float64(gss.HalfSide()*s),
		box2d.B2Vec2{X: float64(p.square.Position().X() * s), Y: float64(p.square.Position().Y() * s)}, 0.0)

	fd = box2d.MakeB2FixtureDef()
	fd.Shape = &b2SquareShape
	fd.Density = 1.0
	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *seesawPhysicsComponent) buildPolygon(world api.IWorld, parent api.INode) error {
	var err error

	// --------------------------------------------------------------
	p.phyNode, err = extras.NewStaticPolygonNode("Polygon", false, world, parent)
	if err != nil {
		return err
	}
	p.phyNode.SetScale(3.0)
	p.phyNode.SetPosition(p.position.X(), p.position.Y())
	gpol := p.phyNode.(*extras.StaticPolygonNode)
	gpol.SetColor(color.NewPaletteInt64(color.LightOrange))

	vertices := []float32{
		-1.0, 2.0, 0.0,
		-1.0, 0.0, 0.0,
		0.0, -3.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 1.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2, 3, 4,
	}

	gpol.Populate(vertices, indices)

	return nil
}

func (p *seesawPhysicsComponent) buildCircle(world api.IWorld, parent api.INode) error {
	var err error

	p.circle, err = extras.NewStaticCircleNode("Circle", false, world, parent)
	if err != nil {
		return err
	}
	p.circle.SetScale(5.0 / 3)
	p.circle.SetPosition(-5.0, 0.0)
	gol2 := p.circle.(*extras.StaticCircleNode)
	gol2.SetColor(color.NewPaletteInt64(color.Green))

	return nil
}

func (p *seesawPhysicsComponent) buildSquare(world api.IWorld, parent api.INode) error {
	var err error

	p.square, err = extras.NewStaticSquareNode("Square", true, false, world, parent)
	if err != nil {
		return err
	}
	p.square.SetScale(5.0 / 3)
	p.square.SetPosition(5.0, 0.0)
	gol2 := p.square.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.Aqua))

	return nil
}
