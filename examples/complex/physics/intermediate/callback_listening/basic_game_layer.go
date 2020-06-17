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

	sqrPhyComp *boxPhysicsComponent
	cirPhyComp *cirPhysicsComponent

	trackerComp *TrackingComponent

	fencePhyComp *fencePhysicsComponent

	sqrNode api.INode
	cirNode api.INode
	triNode api.INode

	// Box 2D system
	b2Gravity box2d.B2Vec2
	b2World   box2d.B2World

	b2GroundBody *box2d.B2Body

	// IO
	downKeyDown  bool
	leftKeyDown  bool
	upKeyDown    bool
	rightKeyDown bool
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

	if err := g.addSquare(); err != nil {
		return err
	}

	if err := g.addCircle(); err != nil {
		return err
	}

	targetSize := 5.0

	g.trackerComp = NewTrackingComponent("TriTrackerComp", g)
	g.trackerComp.Configure(targetSize, entityTriangle, entityCircle|entityBoundary, &g.b2World)
	g.trackerComp.SetPosition(0.0, 0.0)

	// Contacts
	listener := newContactListener()
	lr := listener.(*contactListener)
	lr.addListener(g.trackerComp)
	// lr.addListener(g.boxComp)
	// lr.addListener(g.circleComp)

	g.b2World.SetContactListener(listener)

	// ---------------------------------------------------------
	if err := g.addFence(); err != nil {
		return err
	}

	return nil
}

func (g *gameLayer) addFence() error {
	position := geometry.NewPoint()
	g.fencePhyComp = newFencePhysicsComponent()
	g.fencePhyComp.Build(g.World(), g, &g.b2World, position)

	return nil
}

func (g *gameLayer) addSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(0.0, 5.0)
	g.sqrNode, err = custom.NewStaticSquareNode("Square", true, true, g.World(), g)
	if err != nil {
		return err
	}
	g.sqrNode.SetScale(3.0)
	g.sqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gol2 := g.sqrNode.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.Aqua))

	g.sqrPhyComp = newBoxPhysicsComponent()
	g.sqrPhyComp.Build(&g.b2World, g.sqrNode, fallingSqrPos)

	return nil
}

func (g *gameLayer) addCircle() error {
	var err error

	fallingCirPos := geometry.NewPointUsing(0.0, 15.0)
	g.cirNode, err = custom.NewStaticCircleNode("Circle", true, g.World(), g)
	if err != nil {
		return err
	}
	g.cirNode.SetScale(3.0)
	g.cirNode.SetPosition(fallingCirPos.X(), fallingCirPos.Y())
	gol2 := g.cirNode.(*custom.StaticCircleNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	g.cirPhyComp = newCirPhysicsComponent()
	g.cirPhyComp.Build(&g.b2World, g.cirNode, fallingCirPos)

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	if g.downKeyDown {
		g.sqrPhyComp.MoveDown()
	}
	if g.rightKeyDown {
		g.sqrPhyComp.MoveRight()
	}
	if g.upKeyDown {
		g.sqrPhyComp.MoveUp()
	}
	if g.leftKeyDown {
		g.sqrPhyComp.MoveLeft()
	}
	// Box2D expects a fractional number of dt not ms/frame which is
	// why I use secPerUpdate.
	g.trackerComp.Update()

	// Update Ray
	// gl := g.rayNode.(*custom.LineNode)
	// bodyPos := g.trackerComp.GetPosition()
	// gl.SetPoints(bodyPos.X, bodyPos.Y, g.targetPosition.X(), g.targetPosition.Y())

	// Instruct the world to perform a single step of simulation.
	// It is generally best to keep the time step and iterations fixed.
	g.b2World.Step(secPerUpdate, api.VelocityIterations, api.PositionIterations)

	// -----------------------------------------------------------
	g.sqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.cirPhyComp.Update(msPerUpdate, secPerUpdate)
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
		// fmt.Println(event)
		if event.GetState() == 0 {
			switch event.GetKeyScan() {
			case 65: // A = left
				g.leftKeyDown = false
			case 87: // W = up
				g.upKeyDown = false
			case 68: // D = right
				g.rightKeyDown = false
			case 83: // S = down
				g.downKeyDown = false
			}
		}

		if event.GetState() == 1 || event.GetState() == 2 {
			switch event.GetKeyScan() {
			case 65: // A = left
				g.leftKeyDown = true
			case 87: // W = up
				g.upKeyDown = true
			case 68: // D = right
				g.rightKeyDown = true
			case 83: // S = down
				g.downKeyDown = true
			case 82: // R
				g.sqrPhyComp.Reset()
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
	// g.b2Gravity = box2d.MakeB2Vec2(0.0, -9.8)

	// Construct a world object, which will hold and simulate the rigid bodies.
	g.b2World = box2d.MakeB2World(g.b2Gravity)
}
