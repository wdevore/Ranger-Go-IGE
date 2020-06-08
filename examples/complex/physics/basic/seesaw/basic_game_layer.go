package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	seesawPhyComp  *seesawPhysicsComponent
	seesawPosition api.IPoint

	fallingSqrNode api.INode

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

	g.seesawPosition = geometry.NewPointUsing(0.0, 30.0)

	if err := g.addSeesaw(); err != nil {
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
	g.groundLineNode.SetScale(100.0)
	g.groundLineNode.SetPosition(0.0, -5.0)
	gl2 := g.groundLineNode.(*custom.StaticHLineNode)
	gl2.SetColor(color.NewPaletteInt64(color.Yellow))

	return nil
}

func (g *gameLayer) addSeesaw() error {

	g.seesawPhyComp = newSeesawPhysicsComponent()
	g.seesawPhyComp.Build(g.World(), g, &g.b2World, g.seesawPosition)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	// Box2D expects a fractional number of dt not ms/frame which is
	// why I use secPerUpdate.

	// Instruct the world to perform a single step of simulation.
	// It is generally best to keep the time step and iterations fixed.
	g.b2World.Step(secPerUpdate, api.VelocityIterations, api.PositionIterations)

	// -----------------------------------------------------------
	g.seesawPhyComp.Update(msPerUpdate, secPerUpdate)
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
				g.seesawPhyComp.Reset()
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

func buildGroundPhysics(g *gameLayer) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()

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
