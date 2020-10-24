package main

import (
	"math/rand"

	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/particles"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

// A Mountain polygon and collection of triangles

type lavaPhysicsComponent struct {
	physicsComponent

	categoryBits uint16 // I am a...
	maskBits     uint16 // I can collide with a...

	lava []*triPhysicsComponent

	particleSys  api.IParticleSystem
	particleActr api.IParticleActivator

	minDelay  float64
	maxDelay  float64
	delayTime float32

	delayCnt float32
}

func newLavaPhysicsComponent() *lavaPhysicsComponent {
	o := new(lavaPhysicsComponent)
	return o
}

func (p *lavaPhysicsComponent) Build(phyWorld *box2d.B2World, parent api.INode,
	position api.IPoint, mountainScale, lavaScale float32,
	lavaCount int) {

	p.particleActr = NewActivatorCone()
	p.particleSys = particles.NewParticleSystem(p.particleActr)
	p.particleSys.Activate(true)

	p.buildLava(phyWorld, parent, position, lavaScale, lavaCount)
	p.buildMountain(phyWorld, parent, position, mountainScale)

	p.delayCnt = 0.0
	p.delayTime = 2000.0
	p.minDelay = 500.0
	p.maxDelay = 2000.0
}

func (p *lavaPhysicsComponent) buildLava(
	phyWorld *box2d.B2World, parent api.INode,
	position api.IPoint, scale float32, lavaCount int) {
	p.position = position

	pPos := geometry.NewPoint()
	pPos.SetByComp(position.X(), position.Y()+scale*2.0)
	p.particleSys.SetPosition(pPos.X(), pPos.Y())

	for i := 0; i < lavaCount; i++ {
		particle := newTriPhysicsComponent()
		particle.Configure(phyWorld, parent,
			i, scale, color.NewPaletteInt64(color.YellowGreen),
			pPos)
		particle.ConfigureFilter(entityLava, entityTriangle|entityLava|entityRectangle|entityLand|entityStarShip)
		particle.Build(phyWorld, pPos)
		particle.EnableGravity(true)

		p.lava = append(p.lava, particle)

		p.particleSys.AddParticle(particle)
	}
}

func (p *lavaPhysicsComponent) buildMountain(
	phyWorld *box2d.B2World, parent api.INode,
	position api.IPoint, scale float32) {
	p.position = position

	var err error

	err = p.buildPolygon(parent.World(), parent, scale)
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
	gla := p.phyNode.(*shapes.MonoPolygonNode)
	verts := gla.Vertices()

	for i := 0; i < len(*verts); i += api.XYZComponentCount {
		vertices = append(vertices, box2d.B2Vec2{X: float64((*verts)[i] * scale), Y: float64((*verts)[i+1] * scale)})
	}

	b2ChainShape.CreateChain(vertices, len(vertices))

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2ChainShape
	fd.Filter.CategoryBits = p.categoryBits
	fd.Filter.MaskBits = p.maskBits

	p.b2Body.CreateFixtureFromDef(&fd) // attach Fixture to body
}

func (p *lavaPhysicsComponent) ConfigureFilter(categoryBits, maskBits uint16) {
	p.categoryBits = categoryBits
	p.maskBits = maskBits
}

func (p *lavaPhysicsComponent) buildPolygon(world api.IWorld, parent api.INode, scale float32) error {
	var err error

	// --------------------------------------------------------------
	vertices := []float32{
		0.5, 0.0, 0.0,
		0.2, 0.8, 0.0,
		0.0, 0.6, 0.0,
		-0.2, 0.7, 0.0,
		-0.5, 0.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2,
		0, 2, 4,
		2, 3, 4,
	}

	// --------------------------------------------------------------
	p.phyNode, err = shapes.NewMonoPolygonNode("Mountain", &vertices, &indices, api.FILLED, world, parent)
	if err != nil {
		return err
	}
	p.phyNode.SetScale(scale)
	p.phyNode.SetPosition(p.position.X(), p.position.Y())
	gpol := p.phyNode.(*shapes.MonoPolygonNode)
	gpol.SetColor(color.NewPaletteInt64(color.Brick))

	return nil
}

func (p *lavaPhysicsComponent) Update(dt float32) {
	p.particleSys.Update(dt)

	// Update physics
	for _, particle := range p.lava {
		particle.Update()
	}

	if p.delayCnt > p.delayTime {
		p.delayTime = float32(maths.Lerp(p.minDelay, p.maxDelay, rand.Float64()))
		p.delayCnt = 0.0
		p.TriggerOneshot()
	}

	p.delayCnt += dt
}

// TriggerOneshot activates a single particle
func (p *lavaPhysicsComponent) TriggerOneshot() {
	p.particleSys.TriggerOneshot()
}
