package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type menuLayer struct {
	nodes.Node

	atlas api.IAtlasX
}

func newMenuLayer(name string, atlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(menuLayer)
	o.Initialize(name)
	o.SetParent(parent)
	o.atlas = atlas
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *menuLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	square, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(0.0, 200.0)
	gsq := square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.GoldYellow))

	g.addLine("Select Choice", -100.0, 100.0, 25, color.NewPaletteInt64(color.White), world)

	g.addLine("1 = Settings", -100.0, 65.0, 25, color.NewPaletteInt64(color.Lime), world)
	g.addLine("2 = Highscore", -100.0, 40.0, 25, color.NewPaletteInt64(color.Pink), world)
	g.addLine("3 = Game", -100.0, 15.0, 25, color.NewPaletteInt64(color.Yellow), world)
	g.addLine("x = To Exit", -100.0, -10.0, 25, color.NewPaletteInt64(color.Red), world)

	return nil
}

func (g *menuLayer) addLine(text string, x, y, s float32, textColor api.IPalette, world api.IWorld) error {
	textureNode, err := shapes.NewBitmapFont9x9Node("SettingsLine", g.atlas, world, g)
	if err != nil {
		return err
	}
	textureNode.SetScale(s)
	textureNode.SetPosition(x, y)
	bf := textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(textColor.Array())
	bf.SetText(text)

	return nil
}
