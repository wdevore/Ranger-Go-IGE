package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	fallingCirPos    api.IPoint
	fallingCirNode   api.INode
	b2FallingCirBody *box2d.B2Body

	fallingSqrPos    api.IPoint
	fallingSqrNode   api.INode
	b2FallingSqrBody *box2d.B2Body

	fallingTriPos    api.IPoint
	fallingTriNode   api.INode
	b2FallingTriBody *box2d.B2Body

	groundLineNode api.INode

	// Box 2D system
	b2Gravity box2d.B2Vec2
	b2World   box2d.B2World

	b2GroundBody *box2d.B2Body
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	setupPhysicsWorld(g)

	if err := g.addCircle(); err != nil {
		return err
	}

	if err := g.addSquare(); err != nil {
		return err
	}

	if err := g.addTri(); err != nil {
		return err
	}

	// ---------------------------------------------------------
	if err := g.addGround(); err != nil {
		return err
	}

	buildGroundPhysics(g)

	return nil
}

func (g *gameLayer) addGround() error {
	var err error

	g.groundLineNode, err = custom.NewStaticHLineNode("Ground", g.World(), g)
	if err != nil {
		return err
	}
	g.groundLineNode.SetScale(25.0)
	g.groundLineNode.SetRotation(maths.DegreeToRadians * 10.0)
	g.groundLineNode.SetPosition(0.0, 0.0)
	gl2 := g.groundLineNode.(*custom.StaticHLineNode)
	gl2.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}

func (g *gameLayer) addCircle() error {
	var err error

	g.fallingCirPos = geometry.NewPointUsing(0.0, 15.0)
	g.fallingCirNode, err = custom.NewStaticCircleNode("Circle", true, g.World(), g)
	if err != nil {
		return err
	}
	g.fallingCirNode.SetScale(3.0)
	// g.circleNode.SetPosition(100.0, 100.0)
	g.fallingCirNode.SetPosition(g.fallingCirPos.X(), g.fallingCirPos.Y())
	gol2 := g.fallingCirNode.(*custom.StaticCircleNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	buildCirclePhysics(g)

	return nil
}

func (g *gameLayer) addSquare() error {
	var err error

	g.fallingSqrPos = geometry.NewPointUsing(0.0, 10.0)
	g.fallingSqrNode, err = custom.NewStaticSquareNode("Square", true, true, g.World(), g)
	if err != nil {
		return err
	}
	g.fallingSqrNode.SetScale(3.0)
	g.fallingSqrNode.SetPosition(g.fallingSqrPos.X(), g.fallingSqrPos.Y())
	gol2 := g.fallingSqrNode.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.Aqua))

	buildSquarePhysics(g)

	return nil
}

func (g *gameLayer) addTri() error {
	var err error

	g.fallingTriPos = geometry.NewPointUsing(0.0, 5.0)
	g.fallingTriNode, err = custom.NewStaticTriangleNode("Triangle", true, true, g.World(), g)
	if err != nil {
		return err
	}
	g.fallingTriNode.SetScale(3.0)
	g.fallingTriNode.SetPosition(g.fallingTriPos.X(), g.fallingTriPos.Y())
	gol2 := g.fallingTriNode.(*custom.StaticTriangleNode)
	gol2.SetColor(color.NewPaletteInt64(color.Pink))

	buildTrianglePhysics(g)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	// Box2D expects a fractional number of dt not ms/frame which is
	// why I use secPerUpdate.

	// Instruct the world to perform a single step of simulation.
	// It is generally best to keep the time step and iterations fixed.
	g.b2World.Step(secPerUpdate, api.VelocityIterations, api.PositionIterations)

	g.UpdateNode(msPerUpdate, secPerUpdate)
}

func (g *gameLayer) UpdateNode(msPerUpdate, secPerUpdate float64) {
	if g.b2FallingCirBody.IsActive() {
		pos := g.b2FallingCirBody.GetPosition()
		g.fallingCirNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot := g.b2FallingCirBody.GetAngle()
		g.fallingCirNode.SetRotation(rot)

		// -----------------------------------------------------------
		pos = g.b2FallingSqrBody.GetPosition()
		g.fallingSqrNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot = g.b2FallingSqrBody.GetAngle()
		g.fallingSqrNode.SetRotation(rot)

		// -----------------------------------------------------------
		pos = g.b2FallingTriBody.GetPosition()
		g.fallingTriNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot = g.b2FallingTriBody.GetAngle()
		g.fallingTriNode.SetRotation(rot)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
	man.UnRegisterEventTarget(g)
	g.b2World.Destroy()
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (g *gameLayer) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		switch event.GetKeyScan() {
		case 65: // A
			if event.GetState() == 1 {
			}
		case 82: // R
			if event.GetState() == 1 {
				x := g.fallingCirPos.X()
				y := g.fallingCirPos.Y()

				g.fallingCirNode.SetPosition(float32(x), float32(y))
				g.b2FallingCirBody.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
				g.b2FallingCirBody.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
				g.b2FallingCirBody.SetAngularVelocity(0.0)
				g.b2FallingCirBody.SetAwake(true)

				x = g.fallingSqrPos.X()
				y = g.fallingSqrPos.Y()

				g.fallingSqrNode.SetPosition(float32(x), float32(y))
				g.b2FallingSqrBody.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
				g.b2FallingSqrBody.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
				g.b2FallingSqrBody.SetAngularVelocity(0.0)
				g.b2FallingSqrBody.SetAwake(true)

				x = g.fallingTriPos.X()
				y = g.fallingTriPos.Y()

				g.fallingTriNode.SetPosition(float32(x), float32(y))
				g.b2FallingTriBody.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
				g.b2FallingTriBody.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
				g.b2FallingTriBody.SetAngularVelocity(0.0)
				g.b2FallingTriBody.SetAwake(true)
			}
		}
	}

	return false
}

func setupPhysicsWorld(g *gameLayer) {
	// --------------------------------------------
	// Box 2d configuration
	// --------------------------------------------

	// Define the gravity vector.
	// Ranger's device coordinate space is oriented the same as OpenGL
	// ^ +Y
	// |
	// |
	// |
	// .--------> +X
	// Thus gravity is specified as negative for downward motion.
	g.b2Gravity = box2d.MakeB2Vec2(0.0, -9.8)

	// Construct a world object, which will hold and simulate the rigid bodies.
	g.b2World = box2d.MakeB2World(g.b2Gravity)
}

func buildSquarePhysics(g *gameLayer) {
	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Set the position of the Body
	px := g.fallingSqrNode.Position().X()
	py := g.fallingSqrNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	g.b2FallingSqrBody = g.b2World.CreateBody(&bDef)

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()

	// Box2D assumes the same is defined in unit-space which
	// means if the object is defined otherwise we need the object
	// to return the correct value
	tcc := g.fallingSqrNode.(*custom.StaticSquareNode)
	b2Shape.SetAsBoxFromCenterAndAngle(
		float64(tcc.HalfSide()), float64(tcc.HalfSide()),
		box2d.B2Vec2{X: 0.0, Y: 0.0}, 0.0)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	g.b2FallingSqrBody.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildTrianglePhysics(g *gameLayer) {
	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Set the position of the Body
	px := g.fallingTriNode.Position().X()
	py := g.fallingTriNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	g.b2FallingTriBody = g.b2World.CreateBody(&bDef)

	tcc := g.fallingTriNode.(*custom.StaticTriangleNode)
	// Box2D expects polygon edges to be defined at full length, not
	// half-side
	scale := tcc.SideLength()
	verts := tcc.Vertices()

	vertices := []box2d.B2Vec2{}
	for i := 0; i < len(verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64(verts[i] * scale), Y: float64(verts[i+1] * scale)})
	}
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[0] * scale), Y: float64(verts[1] * scale)})
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[3] * scale), Y: float64(verts[4] * scale)})
	// vertices = append(vertices, box2d.B2Vec2{X: float64(verts[6] * scale), Y: float64(verts[7] * scale)})

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	b2Shape.Set(vertices, len(verts)/api.XYZComponentCount)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	g.b2FallingTriBody.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildCirclePhysics(g *gameLayer) {
	// -------------------------------------------
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	px := g.fallingCirNode.Position().X()
	py := g.fallingCirNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	g.b2FallingCirBody = g.b2World.CreateBody(&bDef)

	// Every Fixture has a shape
	circleShape := box2d.MakeB2CircleShape()
	circleShape.M_p.Set(0.0, 0.0) // Relative to body position
	tcc := g.fallingCirNode.(*custom.StaticCircleNode)
	radius := tcc.Radius()
	circleShape.M_radius = float64(radius)

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &circleShape
	fd.Density = 1.0
	g.b2FallingCirBody.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildGroundPhysics(g *gameLayer) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	wx := g.groundLineNode.Position().X()
	wy := g.groundLineNode.Position().Y()

	// The Ground = body + fixture + shape
	bDef.Type = box2d.B2BodyType.B2_staticBody
	bDef.Position.Set(
		float64(wx),
		float64(wy),
	)
	bDef.Angle = g.groundLineNode.Rotation()

	g.b2GroundBody = g.b2World.CreateBody(&bDef)

	tln := g.groundLineNode.(*custom.StaticHLineNode)
	halfLength := float64(tln.HalfLength())
	groundShape := box2d.MakeB2EdgeShape()
	groundShape.Set(
		box2d.MakeB2Vec2(
			-halfLength,
			0.0),
		box2d.MakeB2Vec2(
			halfLength,
			0.0),
	)

	fDef := box2d.MakeB2FixtureDef()
	fDef.Shape = &groundShape
	fDef.Density = 1.0
	g.b2GroundBody.CreateFixtureFromDef(&fDef) // attach Fixture to body
}
