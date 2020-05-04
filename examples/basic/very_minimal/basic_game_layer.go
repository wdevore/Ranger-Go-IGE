package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

type gameLayer struct {
	nodes.Node

	backgroundColor api.IPalette
	backgroundMin   api.IPoint
	backgroundMax   api.IPoint

	textColor api.IPalette
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

func (g *gameLayer) Build(world api.IWorld) error {
	// vw, vh := world.ViewSize().Components()
	// x := -vw / 2.0
	// y := -vh / 2.0

	// g.backgroundMin = geometry.NewPointUsing(x, y)
	// g.backgroundMax = geometry.NewPointUsing(x+vw, y+vh)

	g.backgroundColor = rendering.NewPaletteInt64(rendering.DarkGray)
	g.textColor = rendering.NewPaletteInt64(rendering.LightNavyBlue)

	return nil
}

// -----------------------------------------------------
// Visuals
// -----------------------------------------------------

var o1 = geometry.NewPoint()
var o2 = geometry.NewPoint()

func (g *gameLayer) Draw(m4 api.IMatrix4) {
	// Transform vertices if anything has changed.
	if g.IsDirty() {
		// // Transform this node's vertices using the context
		// context.TransformPoint(g.backgroundMin, o1)
		// context.TransformPoint(g.backgroundMax, o2)
		g.SetDirty(false) // Node is no longer dirty
	}

	// Draw background first.
	// context.SetDrawColor(g.backgroundColor)
	// context.RenderAARectangle(o1, o2, api.FILLED)

	// context.SetDrawColor(g.textColor)
	// context.DrawText(450.0, 250.0, "Game Layer", 6, 1, false)

	// g.Node.Draw(context)
}
