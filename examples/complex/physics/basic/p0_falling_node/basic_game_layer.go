package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	circle api.INode

	// Box 2D system
	b2Gravity box2d.B2Vec2
	b2World   box2d.B2World

	b2Body *box2d.B2Body
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

	var err error

	// ---------------------------------------------------------
	g.circle, err = shapes.NewMonoCircleNode("Circle", api.FILLOUTLINED, 6, world, g)
	if err != nil {
		return err
	}
	g.circle.SetScale(5.0)
	g.circle.SetPosition(30.0, 80.0)
	gc := g.circle.(*shapes.MonoCircleNode)
	gc.SetFilledColor(color.NewPaletteInt64(color.LightOrange))

	// Instead of using two node: vline and hline, I'm using one "+ node.
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	dvr := world.Properties().Window.DeviceRes
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

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

	// A body def used to create body
	bd := box2d.MakeB2BodyDef()
	bd.Type = box2d.B2BodyType.B2_dynamicBody
	bd.Position.Set(float64(g.circle.Position().X()), float64(g.circle.Position().Y()))

	// An instance of a body to contain the Fixtures
	g.b2Body = g.b2World.CreateBody(&bd)

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2CircleShape()
	b2Shape.M_p.Set(0.0, 0.0) // Relative to body position
	b2Shape.M_radius = float64(g.circle.Scale())

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	g.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	// Box2D expects a fractional number of dt not ms/frame which is
	// why I use secPerUpdate.

	// Instruct the world to perform a single step of simulation.
	// It is generally best to keep the time step and iterations fixed.
	g.b2World.Step(secPerUpdate, api.VelocityIterations, api.PositionIterations)

	pos := g.b2Body.GetPosition()
	if g.b2Body.IsActive() {
		g.circle.SetPosition(float32(pos.X), float32(pos.Y))
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
				x := 25.0
				y := 25.0

				g.circle.SetPosition(float32(x), float32(y))
				g.b2Body.SetTransform(box2d.MakeB2Vec2(x, y), 0.0)
				g.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
				g.b2Body.SetAngularVelocity(0.0)
			}
		}
	}

	return false
}
