package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type menuLayer struct {
	nodes.Node
}

func newMenuLayer(name string, world api.IWorld, fontRenderer api.ITextureRenderer, parent api.INode) (api.INode, error) {
	o := new(menuLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world, fontRenderer); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *menuLayer) build(world api.IWorld, fontRenderer api.ITextureRenderer) error {
	g.Node.Build(world)

	osql, err := extras.NewStaticSquareNode("FilledSqr", true, true, world, g)
	if err != nil {
		return err
	}

	osql.SetScale(100.0)
	osql.SetPosition(0.0, 200.0)
	gol2 := osql.(*extras.StaticSquareNode)
	gol2.SetColor(color.NewPaletteInt64(color.GoldYellow))

	textureMan := world.TextureManager()

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.addLine("Select Choice", -100.0, 100.0, 25, color.NewPaletteInt64(color.White), world, textureAtlas, fontRenderer)

	g.addLine("1 = Settings", -100.0, 65.0, 25, color.NewPaletteInt64(color.Lime), world, textureAtlas, fontRenderer)
	g.addLine("2 = Highscore", -100.0, 40.0, 25, color.NewPaletteInt64(color.Pink), world, textureAtlas, fontRenderer)
	g.addLine("3 = Game", -100.0, 15.0, 25, color.NewPaletteInt64(color.Yellow), world, textureAtlas, fontRenderer)
	g.addLine("x = To Exit", -100.0, -10.0, 25, color.NewPaletteInt64(color.Red), world, textureAtlas, fontRenderer)

	return nil
}

func (g *menuLayer) addLine(text string, x, y, s float32, color api.IPalette, world api.IWorld, textureAtlas api.ITextureAtlas, fontRenderer api.ITextureRenderer) {
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
