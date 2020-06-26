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

const (
	entityBoundary      = uint16(1)
	entityCircle        = uint16(1 << 2)
	entityTriangle      = uint16(1 << 3)
	entityRectangle     = uint16(1 << 4)
	entityStarShip      = uint16(1 << 5)
	entityStarShipRight = uint16(1 << 6)
	entityStarShipLeft  = uint16(1 << 7)
	entityLand          = uint16(1 << 8)
	entityMountain      = uint16(1 << 9)
	entityLava          = uint16(1 << 10)
)

const (
	objectRightZone = 2000
	objectLeftZone  = 2001
)

type gameLayer struct {
	nodes.Node

	starShipComp *StarShipComponent
	landComp     *landPhysicsComponent

	sqrPhyComp *boxPhysicsComponent

	trackerComp *TrackingComponent
	gamePoint   api.IPoint

	fencePhyComp *fencePhysicsComponent

	lavaPhyComp *lavaPhysicsComponent

	sqrNode api.INode
	triNode api.INode

	rayNode api.INode

	// Box 2D system
	b2Gravity box2d.B2Vec2
	b2World   box2d.B2World

	b2GroundBody *box2d.B2Body

	zoneMan *zoneManager

	// Scrolling
	scrollRing     *geometry.Circle
	softness       float32 // Controls how soft the scroll reacts
	scrollRingNode api.INode
	mx, my         int32 // Mouse-space position

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

	g.zoneMan = newZoneManager(g)
	g.zoneMan.Build(world)
	zoom := g.zoneMan.GetZoom()

	if err := g.addLava(zoom); err != nil {
		return err
	}

	if err := g.addSquare(zoom); err != nil {
		return err
	}

	if err := g.addRay(zoom); err != nil {
		return err
	}

	targetSize := float32(3.5)

	g.gamePoint = geometry.NewPoint()

	g.trackerComp = NewTrackingComponent("TriTrackerComp", zoom)
	g.trackerComp.Configure(float64(targetSize),
		entityTriangle, entityLava|entityStarShip|entityBoundary,
		&g.b2World)
	g.trackerComp.SetPosition(10.0, 0.0)
	g.trackerComp.EnableGravity(false)

	g.starShipComp = NewStarShipComponent("StarShip", zoom)
	g.starShipComp.Configure(5.0, entityStarShip,
		entityLava|entityMountain|entityLand|entityTriangle|entityRectangle|entityBoundary,
		&g.b2World)
	g.starShipComp.Reset(-15.0, -10.0)

	pos := geometry.NewPointUsing(0.0, -30.0)
	g.landComp = newLandPhysicsComponent()
	g.landComp.ConfigureFilter(entityLand,
		entityLava|entityRectangle|entityTriangle|entityStarShip|entityStarShipRight|entityStarShipLeft)
	g.landComp.Build(&g.b2World, zoom, pos)

	filter := newFilterListener()
	g.b2World.SetContactFilter(filter)

	// ---------------------------------------------------------
	if err := g.addFence(zoom); err != nil {
		return err
	}

	g.scrollRing = geometry.NewCircle()
	g.scrollRing.SetRadius(15.0)
	if err := g.addScrollRing(zoom, g.scrollRing); err != nil {
		return err
	}
	g.softness = 10.0

	return nil
}

func (g *gameLayer) addFence(parent api.INode) error {
	position := geometry.NewPoint()
	g.fencePhyComp = newFencePhysicsComponent()
	g.fencePhyComp.ConfigureFilter(entityBoundary,
		entityTriangle|entityStarShip|entityStarShipLeft|entityStarShipRight|entityRectangle)

	g.fencePhyComp.Build(g.World(), parent, &g.b2World, position)

	return nil
}

func (g *gameLayer) addLava(parent api.INode) error {
	g.lavaPhyComp = newLavaPhysicsComponent()
	g.lavaPhyComp.ConfigureFilter(entityMountain,
		entityTriangle|entityStarShip|entityStarShipLeft|entityStarShipRight)

	position := geometry.NewPointUsing(10.0, -28.0)
	g.lavaPhyComp.Build(&g.b2World, parent, position, 7.0, 1.5, 15)

	return nil
}

func (g *gameLayer) addSquare(parent api.INode) error {
	var err error

	fallingSqrPos := geometry.NewPointUsing(25.0, -10.0)
	g.sqrNode, err = custom.NewStaticSquareNode("Square", true, true, g.World(), parent)
	if err != nil {
		return err
	}
	g.sqrNode.SetScale(3.0)
	g.sqrNode.SetPosition(fallingSqrPos.X(), fallingSqrPos.Y())
	gol2 := g.sqrNode.(*custom.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.GoldYellow))

	g.sqrPhyComp = newBoxPhysicsComponent()
	g.sqrPhyComp.ConfigureFilter(entityRectangle, entityLava|entityLand|entityStarShip|entityBoundary)
	g.sqrPhyComp.Build(&g.b2World, g.sqrNode, fallingSqrPos)
	g.sqrPhyComp.EnableGravity(false)

	return nil
}

func (g *gameLayer) addRay(parent api.INode) error {
	var err error

	g.rayNode, err = newDynamicLineNode("DynoLin", g.World(), parent)
	if err != nil {
		return err
	}
	glc := g.rayNode.(*DynamicLineNode)
	glc.SetColor(color.NewPaletteInt64(color.White))

	return nil
}

func (g *gameLayer) addScrollRing(parent api.INode, ring *geometry.Circle) error {
	var err error

	g.scrollRingNode, err = custom.NewStaticCircleNode("ScrollRing", false, g.World(), g)
	if err != nil {
		return err
	}
	g.scrollRingNode.SetScale(ring.Radius() * 2)
	g.scrollRingNode.SetVisible(false)
	pos := geometry.NewPointUsing(0.0, 0.0)
	g.scrollRingNode.SetPosition(pos.X(), pos.Y())
	gol2 := g.scrollRingNode.(*custom.StaticCircleNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	return nil
}

var _nodePoint = geometry.NewPoint()

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

	g.lavaPhyComp.Update(float32(msPerUpdate))

	// Update Ray
	ray := g.rayNode.(*DynamicLineNode)
	bodyPos := g.trackerComp.GetPosition()
	ray.SetPoints(float32(bodyPos.X), float32(bodyPos.Y), g.gamePoint.X(), g.gamePoint.Y())

	// -----------------------------------------------------------
	g.sqrPhyComp.Update(msPerUpdate, secPerUpdate)
	g.starShipComp.Update()

	shipPos := g.starShipComp.Position()

	g.zoneMan.UpdateCheck(shipPos, msPerUpdate)

	nodes.MapNodeToNode(g.starShipComp.hullVisual, g, _nodePoint, g)
	dRange := g.scrollRing.DistanceFromEdge(_nodePoint)
	// fmt.Println(fmt.Sprintf("Range(%2.3f) %s, %v", dRange, shipPos, _nodePoint))

	if dRange > 0 {
		gc := g.scrollRingNode.(*custom.StaticCircleNode)
		gc.SetColor(color.NewPaletteInt64(color.Aqua))

		ray := maths.NewVectorUsing(_nodePoint.X(), _nodePoint.Y())
		ray.Normalize()
		ray.Scale(-dRange / g.softness)

		zoom := g.zoneMan.GetZoom()
		gz := zoom.(*custom.ZoomNode)
		gz.TranslateBy(ray.X(), ray.Y())

	} else {
		gc := g.scrollRingNode.(*custom.StaticCircleNode)
		gc.SetColor(color.NewPaletteInt64(color.LightestGray))
	}

	g.cursorTrack(g.mx, g.my)
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
			case 262: // right arrow
				g.starShipComp.EnableYaw(true, 0.0)
			case 263: // left arrow
				g.starShipComp.EnableYaw(true, 0.0)
			case 90:
				g.starShipComp.SetThrust(false)
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
				g.starShipComp.Reset(10.0, -30.0)
			case 262: // right arrow
				g.starShipComp.EnableYaw(true, -2.0)
			case 263: // left arrow
				g.starShipComp.EnableYaw(true, 2.0)
			case 90:
				g.starShipComp.SetThrust(true)
			case 70:
				g.trackerComp.Thrust()
			case 86:
				visible := g.scrollRingNode.IsVisible()
				g.scrollRingNode.SetVisible(!visible)
			case 84: // T
				g.lavaPhyComp.TriggerOneshot()
			}
		}
	case api.IOTypeMouseMotion:
		g.mx, g.my = event.GetMousePosition()
		g.cursorTrack(g.mx, g.my)
	}

	return false
}

func (g *gameLayer) cursorTrack(x, y int32) {
	// Map mouse-space to gamelayer space
	nodes.MapDeviceToNode(x, y, g.zoneMan.zoom, g.gamePoint)

	g.trackerComp.SetTargetPosition(g.gamePoint)
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
