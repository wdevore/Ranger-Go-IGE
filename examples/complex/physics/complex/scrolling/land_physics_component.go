package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type landPhysicsComponent struct {
	physicsComponent

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...
}

func newLandPhysicsComponent() *landPhysicsComponent {
	o := new(landPhysicsComponent)
	return o
}

func (p *landPhysicsComponent) Build(phyWorld *box2d.B2World, parent api.INode, position api.IPoint) {
	p.position = position

	var err error

	err = p.buildPolygon(parent.World(), parent)
	if err != nil {
		panic(err)
	}

	// A body def used to create bodies
	bDef := box2d.MakeB2BodyDef()
	bDef.Type = box2d.B2BodyType.B2_staticBody

	px := p.phyNode.Position().X()
	py := p.phyNode.Position().Y()
	bDef.Position.Set(
		float64(px),
		float64(py),
	)
	// An instance of a body to contain Fixtures
	p.b2Body = phyWorld.CreateBody(&bDef)

	// Every Fixture has a shape
	b2ChainShape := box2d.MakeB2ChainShape()

	vertices := []box2d.B2Vec2{}
	gla := p.phyNode.(*custom.StaticLineStripNode)
	verts := gla.Vertices()
	scale := p.phyNode.Scale()

	for i := 0; i < len(verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64(verts[i] * scale), Y: float64(verts[i+1] * scale)})
	}

	b2ChainShape.CreateChain(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2ChainShape
	// fd.UserData = b.land
	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *landPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

func (p *landPhysicsComponent) buildPolygon(world api.IWorld, parent api.INode) error {
	var err error

	// --------------------------------------------------------------
	p.phyNode, err = custom.NewStaticLineStripNode("Land", world, parent)
	if err != nil {
		return err
	}
	p.phyNode.SetScale(2.0)
	p.phyNode.SetPosition(p.position.X(), p.position.Y())
	gpol := p.phyNode.(*custom.StaticLineStripNode)
	gpol.SetColor(color.NewPaletteInt64(color.White))

	p.phyNode.Build(world)

	vertices := []float32{
		-30.0, 2.5, 0.0,
		-25.0, 2.5, 0.0,
		-20.0, 7.5, 0.0,
		-10.0, 7.5, 0.0,
		-10.0, 5.5, 0.0,
		-5.0, 2.5, 0.0,
		1.0, 2.5, 0.0,
		1.5, 1.0, 0.0,
		7.0, 1.0, 0.0,
		7.5, 0.0, 0.0,
		10.5, 0.0, 0.0,
		11.0, 2.0, 0.0,
		11.0, 5.0, 0.0,
		15.0, 10.0, 0.0,
		20.0, 10.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2, 3, 4,
		5, 6, 7, 8, 9,
		10, 11, 12, 13, 14,
	}

	gpol.Populate(vertices, indices)

	return nil
}
