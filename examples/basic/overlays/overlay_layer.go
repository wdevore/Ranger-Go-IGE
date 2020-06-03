package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
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

	var err error

	g.text, err = custom.NewDynamicTextNode("Text", 500, world, g)
	if err != nil {
		return err
	}
	gt := g.text.(*custom.DynamicTextNode)
	g.text.SetScale(1.25)
	gt.SetText("Ranger")
	gt.SetPixelSize(10.0)
	gt.SetColor(color.NewPaletteInt64(color.Lime).Array())

	if world.Properties().Engine.ShowTimingInfo {
		g.timing, err = custom.NewDynamicTextNode("TimingInfo", 500, world, g)
		if err != nil {
			return err
		}
		g.timing.SetScale(1.0)
		// Set position to lower-left corner
		dvr := world.Properties().Window.DeviceRes
		g.timing.SetPosition(float32(-dvr.Width/2+10.0), float32(-dvr.Height/2)+10.0)

		gt2 := g.timing.(*custom.DynamicTextNode)
		gt2.SetText("-")
		gt2.SetPixelSize(2.0)
		gt2.SetColor(color.NewPaletteInt64(color.LightGray).Array())
	}

	return nil
}

// Update updates the time properties of a node.
func (g *overlayLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.text.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 0.25

	w := g.World()
	if w.Properties().Engine.ShowTimingInfo {
		gt := g.timing.(*custom.DynamicTextNode)
		s := fmt.Sprintf("f:%d u:%d r:%2.3f", w.Fps(), w.Ups(), w.AvgRender())
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
