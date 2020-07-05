package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	zbar  api.INode
	ozbar api.INode
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

	// g.addBar(world)

	_, err := world.TextureManager().LoadTexture("../../assets/ctype0001.png")
	if err != nil {
		panic(err)
	}

	g.addTexture(world)

	return nil
}

func (g *gameLayer) addTexture(world api.IWorld) {

	imNode, err := custom.NewStaticQuadNode("CTypeShip", world, g)
	if err != nil {
		panic(err)
	}
	imNode.SetScale(10)
	imNode.SetPosition(-100.0, 100.0)
	gzr := imNode.(*custom.StaticQuadNode)
	gzr.SetColor(color.NewPaletteInt64(color.Black))
}

func (g *gameLayer) addBar(world api.IWorld) {
	var err error

	// ---------------------------------------------------------
	g.zbar, err = custom.NewStaticZBarNode("FilledZBar", true, world, g)
	if err != nil {
		panic(err)
	}
	g.zbar.SetScale(100)
	g.zbar.SetPosition(300.0, 100.0)
	gzr := g.zbar.(*custom.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.LightNavyBlue))

	g.ozbar, err = custom.NewStaticZBarNode("OutlineZBar", false, world, g)
	if err != nil {
		panic(err)
	}
	g.ozbar.SetScale(100)
	g.ozbar.SetPosition(300.0, 100.0)
	gzr = g.ozbar.(*custom.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.White))
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	// g.zbar.SetRotation(maths.DegreeToRadians * g.angle)
	// g.ozbar.SetRotation(maths.DegreeToRadians * g.angle)
	// g.angle += 1.25
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
