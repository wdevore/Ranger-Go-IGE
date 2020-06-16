package main

import (
	"github.com/ByteArena/box2d"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

type physicsComponent struct {
	position api.IPoint
	phyNode  api.INode
	b2Body   *box2d.B2Body
}

func (p *physicsComponent) Update(msPerUpdate, secPerUpdate float64) {
	if p.b2Body.IsActive() {
		pos := p.b2Body.GetPosition()
		p.phyNode.SetPosition(float32(pos.X), float32(pos.Y))

		rot := p.b2Body.GetAngle()
		p.phyNode.SetRotation(rot)
	}
}

func (p *physicsComponent) Reset() {
	x := p.position.X()
	y := p.position.Y()

	p.phyNode.SetPosition(float32(x), float32(y))
	p.b2Body.SetTransform(box2d.MakeB2Vec2(float64(x), float64(y)), 0.0)
	p.b2Body.SetLinearVelocity(box2d.MakeB2Vec2(0.0, 0.0))
	p.b2Body.SetAngularVelocity(0.0)
	p.b2Body.SetAwake(true)
}
