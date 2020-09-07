package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type gameLayer struct {
	nodes.Node

	sqr api.INode

	angle float64
}

func newGameLayer(name string, world api.IWorld, fontRenderer api.ITextureRenderer, parent api.INode) (api.INode, error) {
	o := new(gameLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world, fontRenderer); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *gameLayer) build(world api.IWorld, fontRenderer api.ITextureRenderer) error {
	g.Node.Build(world)

	textureMan := world.TextureManager()

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.addLine("Type 'r' to return", -100.0, 0.0, 25, color.NewPaletteInt64(color.GoldYellow), world, textureAtlas, fontRenderer)

	var err error
	g.sqr, err = extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}

	g.sqr.SetScale(100.0)
	g.sqr.SetPosition(100.0, 100.0)
	gol2 := g.sqr.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.LightOrange))

	return nil
}

func (g *gameLayer) addLine(text string, x, y, s float32, color api.IPalette, world api.IWorld, textureAtlas api.ITextureAtlas, fontRenderer api.ITextureRenderer) {
	textureNode, err := extras.NewBitmapFont9x9Node("Ranger", textureAtlas, fontRenderer, world, g)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(s)
	textureNode.SetPosition(x, y)

	tn := textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText(text)
	tn.SetColor(color.Array())
}

func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.sqr.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 1.5
}

func (g *gameLayer) EnterNode(man api.INodeManager) {
	// fmt.Println("gameLayer EnterNode")
	man.RegisterTarget(g)
}

// ExitScene called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	// fmt.Println("gameLayer ExitNode")
	man.UnRegisterTarget(g)
}
