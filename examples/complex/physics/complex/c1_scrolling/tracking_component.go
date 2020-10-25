package main

import (
	"fmt"
	"math"

	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

// TrackingComponent is a box
type TrackingComponent struct {
	parent api.INode
	visual api.INode
	b2Body *box2d.B2Body

	beginContactColor api.IPalette
	endContactColor   api.IPalette

	scale float64

	algorithm         int
	trackingAlgorithm int
	stopping          bool
	targetingRate     float64
	targetPosition    api.IPoint

	thrustEnabled  bool
	thrustStrength float64

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

// NewTrackingComponent constructs a component
func NewTrackingComponent(name string, parent api.INode) *TrackingComponent {
	o := new(TrackingComponent)
	o.parent = parent
	o.algorithm = 1
	o.trackingAlgorithm = 5
	o.targetingRate = 10.0 // default is medium
	o.thrustEnabled = false
	o.thrustStrength = 5.0

	return o
}

// Configure component
func (t *TrackingComponent) Configure(scale float64, categoryBits, maskBits uint16, b2World *box2d.B2World) error {
	t.scale = scale

	t.targetPosition = geometry.NewPoint()

	t.categoryBits = categoryBits
	t.maskBits = maskBits

	t.beginContactColor = color.NewPaletteInt64(color.White)
	t.endContactColor = color.NewPaletteInt64(color.Pink)

	var err error

	t.visual, err = shapes.NewMonoTriangleNode("Triangle", api.OUTLINED, t.parent.World(), t.parent)
	if err != nil {
		return err
	}
	t.visual.SetID(1001)
	t.visual.SetScale(float32(t.scale))
	t.visual.SetPosition(0.0, 0.0)
	gt := t.visual.(*shapes.MonoTriangleNode)
	gt.SetOutlineColor(t.endContactColor)

	buildComp(t, b2World)

	return nil
}

// SetColor sets the visual's color
func (t *TrackingComponent) SetColor(color api.IPalette) {
	gr := t.visual.(*shapes.MonoTriangleNode)
	gr.SetOutlineColor(color)
}

// SetTargetPosition sets the target position. The position is in
// view-space.
func (t *TrackingComponent) SetTargetPosition(p api.IPoint) {
	t.targetPosition.SetByPoint(p)
}

// SetPosition sets component's location.
func (t *TrackingComponent) SetPosition(x, y float32) {
	t.visual.SetPosition(x, y)
	t.b2Body.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), t.b2Body.GetAngle())
}

// GetPosition returns the body's position
func (t *TrackingComponent) GetPosition() box2d.B2Vec2 {
	return t.b2Body.GetPosition()
}

// EnableGravity enables/disables gravity for this component
func (t *TrackingComponent) EnableGravity(enable bool) {
	if enable {
		t.b2Body.SetGravityScale(1.0)
	} else {
		t.b2Body.SetGravityScale(0.0)
	}
}

// Reset configures the component back to defaults
func (t *TrackingComponent) Reset(x, y float32) {
	t.visual.SetPosition(x, y)
	t.b2Body.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
	t.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	t.b2Body.SetAngularVelocity(0.0)
	t.b2Body.SetAwake(true)
}

// SetTrackingAlgo changes the tracking style
func (t *TrackingComponent) SetTrackingAlgo(style int) {
	t.trackingAlgorithm = style
}

// SetTargetingRate changes the how fast the tracking locks on.
func (t *TrackingComponent) SetTargetingRate(rate float64) {
	t.targetingRate = rate
}

// SetVelocityAlgo changes the velocity style
func (t *TrackingComponent) SetVelocityAlgo(style int) {
	t.algorithm = style
}

// Stop set linear velocity to 0
func (t *TrackingComponent) Stop() {
	t.stopping = !t.stopping
	if t.stopping {
		fmt.Println("Stopping...")
	}
}

// MoveLeft applies linear force to box center
func (t *TrackingComponent) MoveLeft(dx float64) {
	velocity := t.b2Body.GetLinearVelocity()
	switch t.algorithm {
	case 1:
		velocity.X -= dx
	case 2:
		velocity.X = math.Max(velocity.X-0.1, -5.0)
	}
	t.b2Body.SetLinearVelocity(velocity)
}

// MoveRight applies linear force to box center
func (t *TrackingComponent) MoveRight(dx float64) {
	velocity := t.b2Body.GetLinearVelocity()
	switch t.algorithm {
	case 1:
		velocity.X += dx
	case 2:
		velocity.X = math.Max(velocity.X+0.1, 5.0)
	}
	t.b2Body.SetLinearVelocity(velocity)
}

// MoveUp applies linear force to box center
func (t *TrackingComponent) MoveUp(dy float64) {
	velocity := t.b2Body.GetLinearVelocity()
	switch t.algorithm {
	case 1:
		velocity.Y -= dy
	case 2:
		velocity.Y = math.Max(velocity.Y-0.1, -5.0)
	}
	t.b2Body.SetLinearVelocity(velocity)
}

// MoveDown applies linear force to box center
func (t *TrackingComponent) MoveDown(dy float64) {
	velocity := t.b2Body.GetLinearVelocity()
	switch t.algorithm {
	case 1:
		velocity.Y += dy
	case 2:
		velocity.Y = math.Max(velocity.Y+0.1, 5.0)
	}
	t.b2Body.SetLinearVelocity(velocity)
}

// Thrust applies force in the direction of the ray
func (t *TrackingComponent) Thrust() {
	t.thrustEnabled = !t.thrustEnabled
}

// ApplyForce applies linear force to box center
func (t *TrackingComponent) ApplyForce(dirX, dirY float64) {
	t.b2Body.ApplyForce(box2d.B2Vec2{X: dirX, Y: dirY}, t.b2Body.GetWorldCenter(), true)
}

// ApplyImpulse applies linear impulse to box center
func (t *TrackingComponent) ApplyImpulse(dirX, dirY float64) {
	t.b2Body.ApplyLinearImpulse(box2d.B2Vec2{X: dirX, Y: dirY}, t.b2Body.GetWorldCenter(), true)
}

// ApplyImpulseToCorner applies linear impulse to 1,1 box corner
// As the box rotates the 1,1 corner rotates which means impulses
// could change the rotation to either CW or CCW.
func (t *TrackingComponent) ApplyImpulseToCorner(dirX, dirY float64) {
	t.b2Body.ApplyLinearImpulse(box2d.B2Vec2{X: dirX, Y: dirY}, t.b2Body.GetWorldPoint(box2d.B2Vec2{X: 1.0, Y: 1.0}), true)
}

// ApplyTorque applies torgue to box center
func (t *TrackingComponent) ApplyTorque(torgue float64) {
	t.b2Body.ApplyTorque(torgue, true)
}

// ApplyAngularImpulse applies angular impulse to box center
func (t *TrackingComponent) ApplyAngularImpulse(impulse float64) {
	t.b2Body.ApplyAngularImpulse(impulse, true)
}

// Update component
func (t *TrackingComponent) Update() {
	if t.b2Body.IsActive() {
		pos := t.b2Body.GetPosition()
		t.visual.SetPosition(float32(pos.X), float32(pos.Y))

		rot := t.b2Body.GetAngle()
		t.visual.SetRotation(rot)
	}

	bodyPos := t.b2Body.GetPosition()
	ray := box2d.MakeB2Vec2(float64(t.targetPosition.X())-bodyPos.X, float64(t.targetPosition.Y())-bodyPos.Y)
	targetAngle := math.Atan2(-ray.X, ray.Y)

	switch t.trackingAlgorithm {
	case 1:
		// Instant targeting.
		t.b2Body.SetTransform(t.b2Body.GetPosition(), targetAngle)
	case 2:
		// This is slow and constantly overshoots
		// Torque unchecked
		nAngle := math.Mod(t.b2Body.GetAngle(), math.Pi*2.0)
		totalRotation := targetAngle - nAngle

		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}

		torque := 10.0
		if totalRotation < 0 {
			torque = -10.0
		}
		t.b2Body.ApplyTorque(torque, true)
	case 3:
		// This is slow but eventually locks on.
		// look ahead more than one time step to adjust the rate
		// at which the correct angle is reached. Here I'm looking ahead
		// 1/3 second forward in time.
		nAngle := math.Mod(t.b2Body.GetAngle(), math.Pi*2.0)
		nextAngle := nAngle + t.b2Body.GetAngularVelocity()/t.targetingRate

		totalRotation := targetAngle - nextAngle
		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}

		torque := 10.0
		if totalRotation < 0 {
			torque = -10.0
		}
		t.b2Body.ApplyTorque(torque, true)
	case 4:
		// This locks on instantly but includes some jittering in the
		// last few frames.
		nAngle := math.Mod(t.b2Body.GetAngle(), math.Pi*2.0)
		nextAngle := nAngle + t.b2Body.GetAngularVelocity()/t.targetingRate

		totalRotation := targetAngle - nextAngle
		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}

		desiredAngularVelocity := totalRotation * t.targetingRate
		torque := t.b2Body.GetInertia() * desiredAngularVelocity / (1 / t.targetingRate)
		t.b2Body.ApplyTorque(torque, true)
	case 5:
		// This is equivalent as case 4 but without the "time" element
		// This locks on instantly but includes some jittering in the
		// last few frames.
		nAngle := math.Mod(t.b2Body.GetAngle(), math.Pi*2.0)
		nextAngle := nAngle + t.b2Body.GetAngularVelocity()/t.targetingRate

		totalRotation := targetAngle - nextAngle
		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}

		desiredAngularVelocity := totalRotation * t.targetingRate
		impulse := t.b2Body.GetInertia() * desiredAngularVelocity
		t.b2Body.ApplyAngularImpulse(impulse, true)
	case 6:
		nAngle := math.Mod(t.b2Body.GetAngle(), math.Pi*2.0)
		nextAngle := nAngle + t.b2Body.GetAngularVelocity()/t.targetingRate
		totalRotation := targetAngle - nextAngle
		if totalRotation < -math.Pi {
			totalRotation += math.Pi * 2.0
		} else if totalRotation > math.Pi {
			totalRotation -= math.Pi * 2.0
		}
		desiredAngularVelocity := totalRotation * t.targetingRate
		change := 5.0 * maths.DegreeToRadians //allow 1 degree rotation per time step
		desiredAngularVelocity = math.Min(change, math.Max(-change, desiredAngularVelocity))
		impulse := t.b2Body.GetInertia() * desiredAngularVelocity
		t.b2Body.ApplyAngularImpulse(impulse, true)
	}

	if t.thrustEnabled {
		bodyPos := t.b2Body.GetPosition()
		ray := box2d.MakeB2Vec2(float64(t.targetPosition.X())-bodyPos.X, float64(t.targetPosition.Y())-bodyPos.Y)
		t.ApplyForce(ray.X*t.thrustStrength, ray.Y*t.thrustStrength)
	}
}

// HandleBeginContact processes BeginContact events
func (t *TrackingComponent) HandleBeginContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*shapes.MonoTriangleNode)

	if !ok {
		n, ok = nodeB.(*shapes.MonoTriangleNode)
	}

	if ok {
		n.SetOutlineColor(t.beginContactColor)
	}

	return false
}

// HandleEndContact processes EndContact events
func (t *TrackingComponent) HandleEndContact(nodeA, nodeB api.INode) bool {
	n, ok := nodeA.(*shapes.MonoTriangleNode)

	if !ok {
		n, ok = nodeB.(*shapes.MonoTriangleNode)
	}

	if ok {
		n.SetOutlineColor(t.endContactColor)
	}

	return false
}

func buildComp(t *TrackingComponent, b2World *box2d.B2World) {
	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_dynamicBody

	// Note: +Y points down in Ranger verses Upward in Box2D's GUI.
	vertices := []box2d.B2Vec2{}

	// Sync the visual's vertices to this physic object
	tr := t.visual.(*shapes.MonoTriangleNode)
	// Box2D expects polygon edges to be defined at full length, not
	// half-side
	scale := tr.SideLength()
	verts := tr.Vertices()

	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scale), Y: float64((*verts)[i+1] * scale)})
	}

	// An instance of a body to contain Fixture
	t.b2Body = b2World.CreateBody(&bDef)

	// Every Fixture has a shape
	b2Shape := box2d.MakeB2PolygonShape()
	b2Shape.Set(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2Shape
	fd.Density = 1.0
	fd.UserData = t.visual
	fd.Filter.CategoryBits = t.categoryBits
	fd.Filter.MaskBits = t.maskBits

	t.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}
