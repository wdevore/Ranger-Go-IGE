package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type settingsLayer struct {
	nodes.Node

	atlas api.IAtlasX
}

func newSettingsLayer(name string, atlas api.IAtlasX, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(settingsLayer)

	o.Initialize(name)
	o.SetParent(parent)
	o.atlas = atlas
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *settingsLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	err := g.addLine("Type 'r' to return", -100.0, 0.0, 25, color.NewPaletteInt64(color.Aqua), world)
	if err != nil {
		return err
	}

	return nil
}

func (g *settingsLayer) addLine(text string, x, y, s float32, textColor api.IPalette, world api.IWorld) error {
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
