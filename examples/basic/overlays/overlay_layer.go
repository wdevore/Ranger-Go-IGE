package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

type overlayLayer struct {
	nodes.Node

	angle float64
	text  api.INode

	timing api.INode
}

func newOverlayLayer(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(overlayLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

func (g *overlayLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	g.text = custom.NewRasterTextNode("Text", world, g)
	g.text.SetScale(1.25)
	gt := g.text.(*custom.RasterTextNode)
	gt.SetText("Ranger")
	gt.SetPixelSize(15.0)
	gt.SetColor(rendering.NewPaletteInt64(rendering.LightOrange))

	if world.Properties().Engine.ShowTimingInfo {
		g.timing = custom.NewRasterTextNode("TimingInfo", world, g)
		g.timing.SetScale(1.0)
		// Set position to lower-left corner
		dvr := world.Properties().Window.DeviceRes
		g.timing.SetPosition(float32(-dvr.Width/2+10.0), float32(-dvr.Height/2)+10.0)

		gt = g.timing.(*custom.RasterTextNode)
		gt.SetText("-")
		gt.SetPixelSize(2.0)
		gt.SetColor(rendering.NewPaletteInt64(rendering.LightGray))
	}

	return nil
}

// Update updates the time properties of a node.
func (g *overlayLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.text.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 0.25

	if g.World().Properties().Engine.ShowTimingInfo {
		gt := g.timing.(*custom.RasterTextNode)
		s := fmt.Sprintf("f:%d u:%d r:%2.3f", g.World().Fps(), g.World().Ups(), g.World().AvgRender())
		gt.SetText(s)
	}

}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *overlayLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *overlayLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
}
