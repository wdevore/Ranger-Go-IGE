package main

import (
	"math"

	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

// StarShipComponent is a triangle physics object
type StarShipComponent struct {
	hullVisual api.INode
	b2BodyHull *box2d.B2Body

	nacelScale      float32
	nacelLongLength float32

	rightNacelVisual  api.INode
	b2BodyRightNacel  *box2d.B2Body
	b2JointRightNacel *box2d.B2PrismaticJoint

	leftNacelVisual  api.INode
	b2BodyLeftNacel  *box2d.B2Body
	b2JointLeftNacel *box2d.B2PrismaticJoint

	scale float32

	maxMotorForce float64
	motorSpeed    float64
	minAxisRange  float64
	maxAxisRange  float64
	leftAnchor    api.IPoint
	rightAnchor   api.IPoint

	thrustEnabled  bool
	thrustStrength float64

	torqueEnabled bool

	targetingRate  float64
	targetPosition api.IPoint // rotates in front of ship
	targetAngle    float64
	yawEnabled     bool
	yawStrength    float64

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

// NewStarShipComponent constructs a component
func NewStarShipComponent(name string, parent api.INode) (*StarShipComponent, error) {
	o := new(StarShipComponent)

	o.nacelScale = 10.0
	// nacelHeightHalfLength := float32(0.4)
	// (4 * 0.4) / 2 = 1.6/2 = 0.8
	// o.nacelLongLength = (o.nacelScale * nacelHeightHalfLength) / 2.0
	o.nacelLongLength = 0.8

	// nacelWidthHalfLength := float32(0.15)

	o.maxMotorForce = 100.0
	o.motorSpeed = 10.0
	o.minAxisRange = -float64(o.nacelLongLength)
	o.maxAxisRange = float64(o.nacelLongLength)

	o.thrustEnabled = false
	o.thrustStrength = 2.0

	var err error

	o.hullVisual, err = shapes.NewMonoCircleNode("MainHull", api.FILLED, 16, parent.World(), parent)
	if err != nil {
		return nil, err
	}
	o.hullVisual.SetPosition(0.0, 0.0)
	gc := o.hullVisual.(*shapes.MonoCircleNode)
	gc.SetFilledColor(color.NewPaletteInt64(color.LightOrange))
	gc.SetFilledAlpha(0.5)

	o.torqueEnabled = true
	o.targetingRate = 30.0
	o.targetPosition = geometry.NewPointUsing(0.0, -1.0)
	o.yawEnabled = false
	o.yawStrength = 0.0

	// Note!! Typically you would parent the nacel's to the hull and let the
	// scenegraph handle the relationship, however that would be incorrect
	// because Box2D will handle the relationship via Joints. So visually it appears
	// as if the nacels are children of the hull but technically they are not.
	o.rightNacelVisual, err = shapes.NewMonoSquareNode("RightNacel", api.FILLED, true, parent.World(), parent)
	if err != nil {
		return nil, err
	}
	o.rightNacelVisual.SetScaleComps(o.nacelScale*0.125, o.nacelScale*float32(o.nacelLongLength/2.5))
	o.rightNacelVisual.SetPosition(0.0, 0.0)
	gsq := o.rightNacelVisual.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightNavyBlue))
	gsq.SetFilledAlpha(0.5)

	o.leftNacelVisual, err = shapes.NewMonoSquareNode("LeftNacel", api.FILLED, true, parent.World(), parent)
	if err != nil {
		return nil, err
	}
	o.leftNacelVisual.SetScaleComps(o.nacelScale*0.125, o.nacelScale*float32(o.nacelLongLength/2.5))
	o.leftNacelVisual.SetPosition(0.0, 0.0)
	gsq = o.leftNacelVisual.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Lime))
	gsq.SetFilledAlpha(0.5)

	return o, nil
}

// Configure component
func (s *StarShipComponent) Configure(scale float32, categoryBits, maskBits uint16, b2World *box2d.B2World) {
	s.scale = scale

	s.categoryBits = categoryBits
	s.maskBits = maskBits

	s.leftAnchor = geometry.NewPointUsing(-1.5, -float32(s.nacelLongLength*2.0))
	s.rightAnchor = geometry.NewPointUsing(1.5, -float32(s.nacelLongLength*2.0))

	buildStarShip(s, b2World)
}

// Position gets the Hull's position
func (s *StarShipComponent) Position() api.IPoint {
	return s.hullVisual.Position()
}

// SetPosition sets component's location.
func (s *StarShipComponent) SetPosition(x, y float64) {
	s.hullVisual.SetPosition(float32(x), float32(y))
	s.b2BodyHull.SetTransform(box2d.MakeB2Vec2(x, y), s.b2BodyHull.GetAngle())

	s.rightNacelVisual.SetPosition(float32(x), float32(y))
	s.b2BodyRightNacel.SetTransform(box2d.MakeB2Vec2(x, y), s.b2BodyRightNacel.GetAngle())

	s.leftNacelVisual.SetPosition(float32(x), float32(y))
	s.b2BodyLeftNacel.SetTransform(box2d.MakeB2Vec2(x, y), s.b2BodyLeftNacel.GetAngle())
}

// EnableGravity enables/disables gravity for this component
func (s *StarShipComponent) EnableGravity(enable bool) {
	if enable {
		s.b2BodyHull.SetGravityScale(1.0)
	} else {
		s.b2BodyHull.SetGravityScale(0.0)
	}
}

// Reset configures the component back to defaults
func (s *StarShipComponent) Reset(x, y float64) {
	s.SetPosition(x, y)

	s.b2BodyHull.SetTransform(box2d.MakeB2Vec2(x, y), 0.0)

	// Typically you want to position bodyB (aka nacels) at the resting
	// position of the range of motion so that the motor isn't forced
	// to move bodyB into the resting position on the first few time steps.
	// If the motor is forced to do so then the nacels jerk as bodies are
	// moved which causes rapid motion.
	// Thus we need to calc the resting position based on the anchor points
	// AND the range of motion.
	//
	// X position is simply the bodyA's anchor coord.
	// Y is more complex because of the prismatic motor. The motor wants to
	// move bodyB to its resting position when no forces are acting on it.
	// Therefore the resting position we be the "anchor position" + "motion range"
	// The range of motion is: abs(minAxisRange) + abs(maxAxisRange)

	motionRange := math.Abs(s.minAxisRange) + math.Abs(s.maxAxisRange)
	s.b2BodyRightNacel.SetTransform(
		box2d.MakeB2Vec2(x+float64(s.rightAnchor.X()), y+float64(s.rightAnchor.Y())-motionRange), 0.0)
	s.b2BodyLeftNacel.SetTransform(
		box2d.MakeB2Vec2(x+float64(s.leftAnchor.X()), y+float64(s.leftAnchor.Y())-motionRange), 0.0)

	s.b2BodyHull.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	s.b2BodyHull.SetAngularVelocity(0.0)
	s.b2BodyHull.SetAwake(true)

	s.b2BodyRightNacel.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	s.b2BodyRightNacel.SetAngularVelocity(0.0)
	s.b2BodyRightNacel.SetAwake(true)

	s.b2BodyLeftNacel.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	s.b2BodyLeftNacel.SetAngularVelocity(0.0)
	s.b2BodyLeftNacel.SetAwake(true)
}

// EnableYaw enables/disables rotation
func (s *StarShipComponent) EnableYaw(enable bool, strength float64) {
	s.yawEnabled = enable
	s.yawStrength = strength
}

// SetThrust enables/disables thrust
func (s *StarShipComponent) SetThrust(enable bool) {
	s.thrustEnabled = enable
}

// ToggleThrust toggles thrust
func (s *StarShipComponent) ToggleThrust() {
	s.thrustEnabled = !s.thrustEnabled
}

// ApplyYaw calculates the next angle to rotate towards
func (s *StarShipComponent) ApplyYaw(dAngle float64) {
	// Take current angle and inc/dec angle to find new angle
	s.targetAngle = math.Mod(s.b2BodyHull.GetAngle()+dAngle, math.Pi*2.0)
}

// ApplyImpulseThrust applies linear impulse opposite of the ship heading
func (s *StarShipComponent) ApplyImpulseThrust() {
	a := s.b2BodyHull.GetAngle() + math.Pi/2.0
	dir := box2d.MakeB2Vec2(math.Cos(a)*s.thrustStrength, math.Sin(a)*s.thrustStrength)

	s.b2BodyHull.ApplyLinearImpulse(dir, s.b2BodyHull.GetWorldCenter(), true)
}

// Update component
func (s *StarShipComponent) Update() {
	if s.b2BodyHull.IsActive() {
		pos := s.b2BodyHull.GetPosition()
		s.hullVisual.SetPosition(float32(pos.X), float32(pos.Y))

		rot := s.b2BodyHull.GetAngle()
		s.hullVisual.SetRotation(rot)
	}

	if s.b2BodyRightNacel.IsActive() {
		pos := s.b2BodyRightNacel.GetPosition()
		s.rightNacelVisual.SetPosition(float32(pos.X), float32(pos.Y))

		rot := s.b2BodyRightNacel.GetAngle()
		s.rightNacelVisual.SetRotation(rot)
	}

	if s.b2BodyLeftNacel.IsActive() {
		pos := s.b2BodyLeftNacel.GetPosition()
		s.leftNacelVisual.SetPosition(float32(pos.X), float32(pos.Y))

		rot := s.b2BodyLeftNacel.GetAngle()
		s.leftNacelVisual.SetRotation(rot)
	}

	if s.thrustEnabled {
		s.ApplyImpulseThrust()
	}

	if s.torqueEnabled {
		s.ApplyYaw(s.yawStrength * maths.DegreeToRadians)

		nAngle := math.Mod(s.b2BodyHull.GetAngle(), math.Pi*2.0)
		nextAngle := nAngle + s.b2BodyHull.GetAngularVelocity()/s.targetingRate

		totalRotation := s.targetAngle - nextAngle
		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}

		desiredAngularVelocity := totalRotation * s.targetingRate
		impulse := s.b2BodyHull.GetInertia() * desiredAngularVelocity
		s.b2BodyHull.ApplyAngularImpulse(impulse, true)
	}
}

func buildStarShip(s *StarShipComponent, b2World *box2d.B2World) {
	buildMainHull(s, b2World)

	buildRightNacel(s, b2World)

	buildRightJoint(s, b2World)

	buildLeftNacel(s, b2World)

	buildLeftJoint(s, b2World)
}

func buildMainHull(s *StarShipComponent, b2World *box2d.B2World) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// An instance of a body to contain main hull
	s.b2BodyHull = b2World.CreateBody(&bDef)

	s.hullVisual.SetScale(s.scale)

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2CircleShape()
	tcc := s.hullVisual.(*shapes.MonoCircleNode)
	radius := tcc.Radius()
	b2Shape.SetRadius(float64(radius))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 0.2

	fd.UserData = s.hullVisual

	fd.Filter.CategoryBits = s.categoryBits
	fd.Filter.MaskBits = s.maskBits

	s.b2BodyHull.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildRightNacel(s *StarShipComponent, b2World *box2d.B2World) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	s.b2BodyRightNacel = b2World.CreateBody(&bDef)

	sr := s.rightNacelVisual.(*shapes.MonoSquareNode)
	scaleX, scaleY := s.rightNacelVisual.ScaleComps()

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	vertices := []box2d.B2Vec2{}
	verts := sr.Vertices()

	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scaleX), Y: float64((*verts)[i+1] * scaleY)})
	}

	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 0.5
	fd.UserData = s.rightNacelVisual

	fd.Filter.CategoryBits = s.categoryBits
	fd.Filter.MaskBits = s.maskBits

	s.b2BodyRightNacel.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildRightJoint(s *StarShipComponent, b2World *box2d.B2World) {
	// --------------------------------------------------------
	// Configure a prismatic joint between the Hull and the right nacel.
	// --------------------------------------------------------
	b2JointDef := box2d.MakeB2PrismaticJointDef()

	// bodyA's local directional vertical axis
	b2JointDef.LocalAxisA.Set(0.0, -1.0) // -1 for +Y pointing upward, +1 for +Y pointing downward

	// The main hull is the reference frame work both nacels
	b2JointDef.BodyA = s.b2BodyHull

	// Set the bodyA's anchor point in prep for bodyB's relative anchor point
	// Note: The coords are in local-space not parent or world-space.
	b2JointDef.LocalAnchorA.Set(float64(s.rightAnchor.X()), float64(s.rightAnchor.Y()))

	b2JointDef.BodyB = s.b2BodyRightNacel

	// Because of the way I will set the range limits, bodyB's anchor point
	// is right in the center of the nacel or (0,0)
	b2JointDef.LocalAnchorB.Set(0.0, 0.0)

	// Now we set the range of motion for this joint. Remember the motion is
	// along the local Y axis relative to the hull
	b2JointDef.LowerTranslation = s.minAxisRange
	b2JointDef.UpperTranslation = s.maxAxisRange

	b2JointDef.EnableLimit = true

	// Set the max force allowed before the motor starts a "breaking" effect.
	// Almost like damping to keep the motor speed limited.
	b2JointDef.MaxMotorForce = s.maxMotorForce

	// Set the motor's functioning speed. The motion of the body could be higher
	// if something is hit but a max is enforced by the MaxMotorForce.
	b2JointDef.MotorSpeed = s.motorSpeed
	b2JointDef.EnableMotor = true

	// We don't want the Hull and the Nacel to collide with each other.
	b2JointDef.CollideConnected = false

	// Finally we create the joint which also adds it to the physics world
	b2World.CreateJoint(&b2JointDef)
}

func buildLeftNacel(s *StarShipComponent, b2World *box2d.B2World) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	s.b2BodyLeftNacel = b2World.CreateBody(&bDef)

	sr := s.leftNacelVisual.(*shapes.MonoSquareNode)
	scaleX, scaleY := s.leftNacelVisual.ScaleComps()

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	vertices := []box2d.B2Vec2{}
	verts := sr.Vertices()

	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scaleX), Y: float64((*verts)[i+1] * scaleY)})
	}

	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 0.5

	fd.UserData = s.leftNacelVisual

	fd.Filter.CategoryBits = s.categoryBits
	fd.Filter.MaskBits = s.maskBits

	s.b2BodyLeftNacel.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func buildLeftJoint(s *StarShipComponent, b2World *box2d.B2World) {
	// --------------------------------------------------------
	// Configure a prismatic joint between the Hull and the right nacel.
	// --------------------------------------------------------
	b2JointDef := box2d.MakeB2PrismaticJointDef()

	// bodyA's local directional vertical axis
	b2JointDef.LocalAxisA.Set(0.0, -1.0)

	// The main hull is the reference frame work both nacels
	b2JointDef.BodyA = s.b2BodyHull

	// Set the bodyA's anchor point in prep for bodyB's relative anchor point
	// Note: The coords are in local-space not parent or world-space.
	b2JointDef.LocalAnchorA.Set(float64(s.leftAnchor.X()), float64(s.leftAnchor.Y()))

	b2JointDef.BodyB = s.b2BodyLeftNacel

	// Because of the way I will set the range limits, bodyB's anchor point
	// is right in the center of the nacel or (0,0)
	b2JointDef.LocalAnchorB.Set(0.0, 0.0)

	// Now we set the range of motion for this joint. Remember the motion is
	// along the local Y axis
	b2JointDef.LowerTranslation = s.minAxisRange
	b2JointDef.UpperTranslation = s.maxAxisRange

	b2JointDef.EnableLimit = true

	// Set the max force allowed before the motor starts a "breaking" effect.
	// Almost like damping to keep the motor speed limited.
	b2JointDef.MaxMotorForce = s.maxMotorForce

	// Set the motor's functioning speed. The motion of the body could be higher
	// if something is hit but a max is enforced by the MaxMotorForce.
	b2JointDef.MotorSpeed = s.motorSpeed
	b2JointDef.EnableMotor = true

	// We don't want the Hull and the Nacel to collide with each other.
	b2JointDef.CollideConnected = false

	// Finally we create the joint which also adds it to the physics world
	b2World.CreateJoint(&b2JointDef)
}
