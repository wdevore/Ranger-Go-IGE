package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	crow  api.INode
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	if err := o.Build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) Build(world api.IWorld) error {
	g.Node.Build(world)

	dvr := world.Properties().Window.DeviceRes

	hline, err := custom.NewStaticAtlasNode("HLine", "HLine", world, g)
	if err != nil {
		return err
	}
	hline.SetScale(float32(dvr.Width))
	ghl := hline.(*custom.StaticAtlasNode)
	ghl.SetColor(color.NewPaletteInt64(color.LightGray))

	vline, err := custom.NewStaticAtlasNode("VLine", "VLine", world, g)
	if err != nil {
		return err
	}
	vline.SetScale(float32(dvr.Height))
	gvl := vline.(*custom.StaticAtlasNode)
	gvl.SetColor(color.NewPaletteInt64(color.LightGray))

	// Square ----------------------------------------------------
	square, err := custom.NewStaticAtlasNode("Square", "CenteredSquare", world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(90.0, 80.0)
	gsq := square.(*custom.StaticAtlasNode)
	gsq.SetColor(color.NewPaletteInt64(color.LightPurple))

	outSquare, err := custom.NewStaticAtlasNode("OutSquare", "CenteredOutlineSquare", world, g)
	if err != nil {
		return err
	}
	outSquare.SetScale(100.0)
	outSquare.SetPosition(90.0, 80.0)
	gsq = outSquare.(*custom.StaticAtlasNode)
	gsq.SetColor(color.NewPaletteInt64(color.White))

	// Circle ----------------------------------------------------
	circle, err := custom.NewStaticAtlasNode("Circle", "Circle12Segments", world, g)
	if err != nil {
		return err
	}
	circle.SetScale(100.0)
	circle.SetPosition(30.0, 80.0)
	gc := circle.(*custom.StaticAtlasNode)
	gc.SetColor(color.NewPaletteInt64(color.GoldYellow))
	gc.SetAlpha(0.5)

	outCircle, err := custom.NewStaticAtlasNode("OutCircle", "Circle12SegmentsOutline", world, g)
	if err != nil {
		return err
	}
	outCircle.SetScale(100.0)
	outCircle.SetPosition(30.0, 80.0)
	gsq = outCircle.(*custom.StaticAtlasNode)
	gsq.SetColor(color.NewPaletteInt64(color.White))

	// Triangle ----------------------------------------------------
	tri, err := custom.NewStaticAtlasNode("Triangle", "CenteredTriangle", world, g)
	if err != nil {
		return err
	}
	tri.SetScale(100.0)
	tri.SetPosition(55.0, 45.0)
	gt := tri.(*custom.StaticAtlasNode)
	gt.SetColor(color.NewPaletteInt64WithAlpha(color.Pink, 0.75))
	outTri, err := custom.NewStaticAtlasNode("OutTri", "CenteredOutlineTriangle", world, g)
	if err != nil {
		return err
	}
	outTri.SetScale(100.0)
	outTri.SetPosition(55.0, 45.0)
	gsq = outTri.(*custom.StaticAtlasNode)
	gsq.SetColor(color.NewPaletteInt64(color.White))

	// Crowbar ----------------------------------------------------
	g.crow, err = custom.NewStaticAtlasNode("Crow", "CrowBar", world, g)
	if err != nil {
		return err
	}
	g.crow.SetScale(100.0)
	g.crow.SetPosition(100.0, 150.0)
	gb := g.crow.(*custom.StaticAtlasNode)
	gb.SetColor(color.NewPaletteInt64WithAlpha(color.White, 0.25))

	return nil
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.crow.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle += 1.25
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
}
