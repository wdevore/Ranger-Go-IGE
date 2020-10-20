package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type highscoreLayer struct {
	nodes.Node

	atlas api.IAtlasX
}

func newHighscoreLayer(name string, atlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(highscoreLayer)
	o.Initialize(name)
	o.atlas = atlas
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *highscoreLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	g.addLine("Type 'r' to return", -100.0, 0.0, 25, color.NewPaletteInt64(color.White), world)

	return nil
}

func (g *highscoreLayer) addLine(text string, x, y, s float32, textColor api.IPalette, world api.IWorld) error {
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
