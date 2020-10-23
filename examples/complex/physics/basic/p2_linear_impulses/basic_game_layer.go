package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	orangeSqrPhyComp *boxPhysicsComponent
	blueSqrPhyComp   *boxPhysicsComponent
	limeSqrPhyComp   *boxPhysicsComponent
	purpleSqrPhyComp *boxPhysicsComponent
	yellowSqrPhyComp *boxPhysicsComponent

	fencePhyComp *fencePhysicsComponent

	orangeSqrNode api.INode
	blueSqrNode   api.INode
	limeSqrNode   api.INode
	purpleSqrNode api.INode
	yellowSqrNode api.INode

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

	if err := o.build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	// Instead of using two node: vline and hline, I'm using one "+ node.
	xyAxis, err := shapes.NewMonoPlusNode("XYAxis", world, world.Underlay())
	if err != nil {
		return err
	}
	dvr := world.Properties().Window.DeviceRes
	xyAxis.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	ghl := xyAxis.(*shapes.MonoPlusNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	setupPhysicsWorld(g)

	if err := g.addOrangeSquare(); err != nil {
		return err
	}

	if err := g.addBlueSquare(); err != nil {
		return err
	}

	if err := g.addLimeSquare(); err != nil {
		return err
	}

	if err := g.addPurpleSquare(); err != nil {
		return err
	}

	if err := g.addYellowSquare(); err != nil {
		return err
	}

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

func (g *gameLayer) addOrangeSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(-20.0, 5.0)

	g.orangeSqrNode, err = shapes.NewMonoSquareNode("OrangeSquare", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.orangeSqrNode.SetScale(3.0)
	g.orangeSqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.orangeSqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightOrange))

	g.orangeSqrPhyComp = newBoxPhysicsComponent()
	g.orangeSqrPhyComp.Build(&g.b2World, g.orangeSqrNode, fallingSqrPos)

	return nil
}

func (g *gameLayer) addBlueSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(-10.0, 5.0)

	g.blueSqrNode, err = shapes.NewMonoSquareNode("BlueSquare", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.blueSqrNode.SetScale(3.0)
	g.blueSqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.blueSqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightNavyBlue))

	g.blueSqrPhyComp = newBoxPhysicsComponent()
	g.blueSqrPhyComp.Build(&g.b2World, g.blueSqrNode, fallingSqrPos)

	return nil
}

func (g *gameLayer) addLimeSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(0.0, 5.0)

	g.limeSqrNode, err = shapes.NewMonoSquareNode("LimeSquare", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.limeSqrNode.SetScale(3.0)
	g.limeSqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.limeSqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Lime))

	g.limeSqrPhyComp = newBoxPhysicsComponent()
	g.limeSqrPhyComp.Build(&g.b2World, g.limeSqrNode, fallingSqrPos)

	return nil
}

func (g *gameLayer) addPurpleSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(10.0, 5.0)

	g.purpleSqrNode, err = shapes.NewMonoSquareNode("PurpleSquare", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.purpleSqrNode.SetScale(3.0)
	g.purpleSqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.purpleSqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightPurple))

	g.purpleSqrPhyComp = newBoxPhysicsComponent()
	g.purpleSqrPhyComp.Build(&g.b2World, g.purpleSqrNode, fallingSqrPos)
	g.purpleSqrPhyComp.EnableGravity(false)

	return nil
}

func (g *gameLayer) addYellowSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(20.0, 5.0)

	g.yellowSqrNode, err = shapes.NewMonoSquareNode("YellowSquare", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.yellowSqrNode.SetScale(3.0)
	g.yellowSqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.yellowSqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Yellow))

	g.yellowSqrPhyComp = newBoxPhysicsComponent()
	g.yellowSqrPhyComp.Build(&g.b2World, g.yellowSqrNode, fallingSqrPos)

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
	g.orangeSqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.blueSqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.limeSqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.purpleSqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.yellowSqrPhyComp.Update(msPerUpdate, secPerUpdate)
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

		if event.GetState() == 1 || event.GetState() == 2 {
			switch event.GetKeyScan() {
			case 68: // d
				g.orangeSqrPhyComp.ApplyForce(0.0, 500.0)
			case 70: // f
				g.blueSqrPhyComp.ApplyImpulse(0.0, 200.0)
			case 90: // z
				g.limeSqrPhyComp.ApplyImpulseToCorner(0.0, 200.0)
			case 65: // a
				g.purpleSqrPhyComp.ApplyTorque(150.0)
			case 83: // s
				g.yellowSqrPhyComp.ApplyAngularImpulse(50.0)
			case 82: // R
				g.orangeSqrPhyComp.Reset()
				g.blueSqrPhyComp.Reset()
				g.limeSqrPhyComp.Reset()
				g.purpleSqrPhyComp.Reset()
				g.yellowSqrPhyComp.Reset()
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
