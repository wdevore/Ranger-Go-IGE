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
