package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type gameExitLayer struct {
	nodes.Node
}

func newBasicExitGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameExitLayer)

	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (g *gameExitLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	if err := g.addText(world); err != nil {
		return err
	}

	square, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, g)
	if err != nil {
		return err
	}
	square.SetScale(100.0)
	square.SetPosition(100.0, 0.0)
	gsq := square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightPurple))

	return nil
}

func (g *gameExitLayer) addText(world api.IWorld) error {
	var err error

	// Note: To render text you need 3 objects:
	// SpriteSheet contains the manifest and font image.
	// SingleTextureAtlas renders a single sub-texture (i.e. character).
	// INode will render strings using the atlas.

	name := "Font9x9"
	// #1 SpriteSheet
	spriteSheet := fonts.NewFont9x9SpriteSheet(name, "font9x9_sprite_sheet_manifest.json")
	spriteSheet.Load("../../assets/", true)

	// #2 TextureAtlas
	atlas := atlas.NewSingleTextureAtlas(name, spriteSheet, world)
	err = atlas.Burn()
	if err != nil {
		return err
	}

	// #3 INode
	textureNode, err := shapes.NewBitmapFont9x9Node(name, atlas, world, g)
	if err != nil {
		return err
	}
	textureNode.SetScale(40)
	textureNode.SetPosition(-300.0, 0.0)
	textureNode.SetRotation(20.0 * maths.DegreeToRadians)
	bf := textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(color.NewPaletteInt64(color.GreenYellow).Array())
	bf.SetText("Exit Scene. Goodbye...")

	return nil
}
