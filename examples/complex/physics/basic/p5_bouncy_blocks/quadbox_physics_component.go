package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type quadBoxPhysicsComponent struct {
	physicsComponent

	boxes []api.INode
}

func newQuadBoxPhysicsComponent() *quadBoxPhysicsComponent {
	o := new(quadBoxPhysicsComponent)
	return o
}

func (p *quadBoxPhysicsComponent) Update(msPerUpdate, secPerUpdate float64) {
	if p.b2Body.IsActive() {
		pos := p.b2Body.GetPosition()
		p.phyNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot := p.b2Body.GetAngle()
		p.phyNode.SetRotation(rot)
	}
}

func (p *quadBoxPhysicsComponent) Reset() {
	x := p.position.X()
	y := p.position.Y()

	p.phyNode.SetPosition(float32(x), float32(y))
	p.b2Body.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
	p.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	p.b2Body.SetAngularVelocity(0.0)
	p.b2Body.SetAwake(true)
}

func (p *quadBoxPhysicsComponent) Build(world api.IWorld, parent api.INode, phyWorld *box2d.B2World, position api.IPoint) {
	// phyNode is the Anchor node.
	p.phyNode, _ = shapes.NewMonoSquareNode("Anchor", api.OUTLINED, true, world, parent)
	// Make the anchor invisible but not the children.
	gsq := p.phyNode.(*shapes.MonoSquareNode)
	gsq.SetOutlineAlpha(0.0)
	p.phyNode.SetPosition(position.X(), position.Y())

	p.buildBoxes(world, p.phyNode)
	p.buildPhysics(phyWorld, position)
}

func (p *quadBoxPhysicsComponent) buildPhysics(phyWorld *box2d.B2World, position api.IPoint) {
	p.position = position

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

	for i := 0.0; i < 4.0; i++ {
		// Every Fixture has a shape
		b2Shape := box2d.MakeB2PolygonShape()

		// Box2D assumes the same is defined in unit-space which
		// means if the object is defined otherwise we need the object
		// to return the correct value
		box := p.boxes[int(i)]
		tcc := box.(*shapes.MonoSquareNode)
		pos := box.Position()
		offsetFromBodyOriginX := float64(pos.X())
		offsetFromBodyOriginY := float64(pos.Y())

		b2Shape.SetAsBoxFromCenterAndAngle(
			float64(tcc.HalfSide()), float64(tcc.HalfSide()),
			box2d.B2Vec2{X: offsetFromBodyOriginX, Y: offsetFromBodyOriginY}, 0.0)

		fd := box2d.MakeB2FixtureDef()
		fd.Shape = &b2Shape
		fd.Density = 1.0
		fd.Friction = i / 4.0
		fd.Restitution = (4.0 - i) / 4.0
		p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
	}
}

func (p *quadBoxPhysicsComponent) buildBoxes(world api.IWorld, parent api.INode) error {
	var err error

	// --------------------------------------------------------------
	node, err := shapes.NewMonoSquareNode("Square1", api.FILLED, true, world, parent)
	if err != nil {
		return err
	}
	node.SetScale(3.0)
	node.SetPosition(1.5, 0.0)
	gsq := node.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Aqua))
	gsq.SetFilledAlpha(0.5)

	p.boxes = append(p.boxes, node)

	// --------------------------------------------------------------
	node, err = shapes.NewMonoSquareNode("Square2", api.FILLED, true, world, parent)
	if err != nil {
		return err
	}
	node.SetScale(3.0)
	node.SetPosition(0.0, 1.5)
	gsq = node.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Aqua))
	gsq.SetFilledAlpha(0.5)

	p.boxes = append(p.boxes, node)

	// --------------------------------------------------------------
	node, err = shapes.NewMonoSquareNode("Square3", api.FILLED, true, world, parent)
	if err != nil {
		return err
	}
	node.SetScale(3.0)
	node.SetPosition(-1.5, 0.0)
	gsq = node.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Aqua))
	gsq.SetFilledAlpha(0.5)

	p.boxes = append(p.boxes, node)

	// --------------------------------------------------------------
	node, err = shapes.NewMonoSquareNode("Square4", api.FILLED, true, world, parent)
	if err != nil {
		return err
	}
	node.SetScale(3.0)
	node.SetPosition(0.0, -1.5)
	gsq = node.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.Aqua))
	gsq.SetFilledAlpha(0.5)

	p.boxes = append(p.boxes, node)

	return nil
}
