package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type settingsLayer struct {
	nodes.Node
}

func newSettingsLayer(name string, world api.IWorld, fontRenderer api.ITextureRenderer, parent api.INode) (api.INode, error) {
	o := new(settingsLayer)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world, fontRenderer); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *settingsLayer) build(world api.IWorld, fontRenderer api.ITextureRenderer) error {
	g.Node.Build(world)

	textureMan := world.TextureManager()

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.addLine("Type 'r' to return", -100.0, 0.0, 25, color.NewPaletteInt64(color.Aqua), world, textureAtlas, fontRenderer)

	return nil
}

func (g *settingsLayer) addLine(text string, x, y, s float32, color api.IPalette, world api.IWorld, textureAtlas api.ITextureAtlas, fontRenderer api.ITextureRenderer) {
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
