package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameLayer struct {
	nodes.Node

	sqrPhyComp *boxPhysicsComponent
	cirPhyComp *cirPhysicsComponent

	trackerComp *TrackingComponent
	gamePoint   api.IPoint

	fencePhyComp *fencePhysicsComponent

	sqrNode api.INode
	cirNode api.INode
	triNode api.INode

	radarNode    api.INode
	radarPhyComp *arcPhysicsComponent

	rayNode api.INode

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

	if err := g.addSquare(); err != nil {
		return err
	}

	if err := g.addCircle(); err != nil {
		return err
	}

	if err := g.addRay(); err != nil {
		return err
	}

	targetSize := 3.5

	g.gamePoint = geometry.NewPoint()

	g.trackerComp = NewTrackingComponent("TriTrackerComp", g)
	g.trackerComp.Configure(targetSize, entityTriangle, entityCircle|entityRadar|entityBoundary, &g.b2World)
	g.trackerComp.SetPosition(0.0, 0.0)

	if err := g.addRadarSensor(); err != nil {
		return err
	}

	// Contacts
	listener := newContactListener()
	lr := listener.(*contactListener)
	lr.addListener(g.trackerComp)
	lr.addListener(g.sqrPhyComp)
	lr.addListener(g.cirPhyComp)
	lr.addListener(g.radarPhyComp)

	g.b2World.SetContactListener(listener)

	filter := newFilterListener()
	g.b2World.SetContactFilter(filter)

	// ---------------------------------------------------------
	if err := g.addFence(); err != nil {
		return err
	}

	return nil
}

func (g *gameLayer) addFence() error {
	position := geometry.NewPoint()
	g.fencePhyComp = newFencePhysicsComponent()
	g.fencePhyComp.ConfigureFilter(entityBoundary, entityTriangle|entityCircle|entityRectangle)

	g.fencePhyComp.Build(g.World(), g, &g.b2World, position)

	return nil
}

func (g *gameLayer) addSquare() error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(0.0, 5.0)

	g.sqrNode, err = shapes.NewMonoSquareNode("Square", api.FILLED, true, g.World(), g)
	if err != nil {
		return err
	}
	g.sqrNode.SetScale(3.0)
	g.sqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gsq := g.sqrNode.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Aqua))

	g.sqrPhyComp = newBoxPhysicsComponent()
	g.sqrPhyComp.ConfigureFilter(entityRectangle, entityCircle|entityBoundary)
	g.sqrPhyComp.Build(&g.b2World, g.sqrNode, fallingSqrPos)

	return nil
}

func (g *gameLayer) addCircle() error {
	var err error

	fallingCirPos := geometry.NewPointUsing(0.0, 15.0)

	g.cirNode, err = shapes.NewMonoCircleNode("Circle", api.FILLED, 12, g.World(), g)
	if err != nil {
		return err
	}
	g.cirNode.SetScale(5.0)
	g.cirNode.SetPosition(fallingCirPos.X(), fallingCirPos.Y())
	gc := g.cirNode.(*shapes.MonoCircleNode)
	gc.SetFilledColor(color.NewPaletteInt64(color.LightOrange))

	g.cirPhyComp = newCirPhysicsComponent()
	g.cirPhyComp.ConfigureFilter(entityCircle, entityRectangle|entityTriangle|entityBoundary)
	g.cirPhyComp.Build(&g.b2World, g.cirNode, fallingCirPos)

	return nil
}

func (g *gameLayer) addRadarSensor() error {
	var err error

	pos := geometry.NewPointUsing(0.0, -15.0)

	startAngle := -45 * maths.DegreeToRadians
	endAngle := 45 * maths.DegreeToRadians
	g.radarNode, err = shapes.NewMonoArcNode("ArcSensor", api.FILLED, 12, startAngle, endAngle, g.World(), g)
	if err != nil {
		return err
	}
	g.radarNode.SetScale(15.0)
	g.radarNode.SetPosition(pos.X(), pos.Y())
	gac := g.radarNode.(*shapes.MonoArcNode)
	gac.SetFilledColor(color.NewPaletteInt64(color.LightPurple))
	gac.SetFilledAlpha(0.75)

	g.radarPhyComp = newArcPhysicsComponent()
	g.radarPhyComp.ConfigureFilter(entityRadar, entityTriangle)
	g.radarPhyComp.Build(&g.b2World, g.radarNode, pos)

	return nil
}

func (g *gameLayer) addRay() error {
	var err error

	g.rayNode, err = shapes.NewDynamicMonoLineNode("DynoLine", g.World(), g)
	if err != nil {
		return err
	}
	gl := g.rayNode.(*shapes.DynamicMonoLineNode)
	gl.SetColor(color.NewPaletteInt64(color.White))

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

	// Instruct the world to perform a single step of simulation.
	// It is generally best to keep the time step and iterations fixed.
	g.b2World.Step(secPerUpdate, api.VelocityIterations, api.PositionIterations)

	// Box2D expects a fractional number of dt not ms/frame which is
	// why I use secPerUpdate.
	g.trackerComp.Update()

	// Update Ray
	ray := g.rayNode.(*shapes.DynamicMonoLineNode)
	bodyPos := g.trackerComp.GetPosition()
	ray.SetVertex1(float32(bodyPos.X), float32(bodyPos.Y))
	ray.SetVertex2(g.gamePoint.X(), g.gamePoint.Y())

	// -----------------------------------------------------------
	g.sqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.cirPhyComp.Update(msPerUpdate, secPerUpdate)

	g.radarPhyComp.Update(msPerUpdate, secPerUpdate)

	dynoAtlas := g.World().GetAtlas(api.DynamicMonoAtlasName)
	dynoAtlas.(api.IDynamicAtlasX).Update()
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
	switch event.GetType() {
	case api.IOTypeKeyboard:
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
			case 70:
				g.trackerComp.Thrust()
			}
		}
	case api.IOTypeMouseMotion:
		mx, my := event.GetMousePosition()

		// Map mouse-space to gamelayer space
		nodes.MapDeviceToNode(mx, my, g, g.gamePoint)

		g.trackerComp.SetTargetPosition(g.gamePoint)
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
