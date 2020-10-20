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

type gameLayer struct {
	nodes.Node
	spriteSheet api.ISpriteSheet

	textureNode api.INode
	textureIdx  int
}

func newBasicGameLayer(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(gameLayer)

	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}
	return o, nil
}

func (g *gameLayer) build(world api.IWorld) error {
	g.Node.Build(world)

	err := g.addFont(world)
	if err != nil {
		return err
	}

	return nil
}

func (g *gameLayer) addFont(world api.IWorld) error {
	var err error

	// Note: To render text you need 3 objects:
	// SpriteSheet contains the manifest and font image.
	// SingleTextureAtlas renders a single sub-texture (i.e. character).
	// INode will render strings using the atlas.

	// #1 SpriteSheet
	g.spriteSheet = fonts.NewFont9x9SpriteSheet("Font9x9", "font9x9_sprite_sheet_manifest.json")
	g.spriteSheet.Load("../../assets/", true)

	// #2 TextureAtlas
	atlas := atlas.NewSingleTextureAtlas("Font9x9", g.spriteSheet, world)
	err = atlas.Burn()
	if err != nil {
		return err
	}

	// #3 INode
	g.textureNode, err = shapes.NewBitmapFont9x9Node("Font9x9", atlas, world, g)
	if err != nil {
		return err
	}
	g.textureNode.SetScale(40)
	g.textureNode.SetPosition(-300.0, 0.0)
	g.textureNode.SetRotation(-20.0 * maths.DegreeToRadians)
	bf := g.textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(color.NewPaletteInt64(color.GreenYellow).Array())
	bf.SetText("Use `<--` and `-->` arrows keys")

	return nil
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterEventTarget(g)
}

func (g *gameLayer) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		// fmt.Println(event)

		if event.GetState() == 1 || event.GetState() == 2 {
			bf := g.textureNode.(*shapes.BitmapFont9x9Node)
			bf.SetColor(color.NewPaletteInt64(color.GoldYellow).Array())
			switch event.GetKeyScan() {
			case 68: // d
				tn := g.textureNode.(*shapes.BitmapFont9x9Node)
				tn.SetIndex(1)
			case 70: // f
				tn := g.textureNode.(*shapes.BitmapFont9x9Node)
				tn.SetIndex(0)
			case 90: // z
			case 65: // a
			case 83: // s
			case 82: // R
			case 263: // Left
				g.textureIdx = (g.textureIdx - 1) % 94
				if g.textureIdx < 0 {
					g.textureIdx = 94 - 1
				}
				tn := g.textureNode.(*shapes.BitmapFont9x9Node)
				tn.SetIndex(g.textureIdx)
			case 262: // Right
				g.textureIdx = (g.textureIdx + 1) % 94
				tn := g.textureNode.(*shapes.BitmapFont9x9Node)
				tn.SetIndex(g.textureIdx)
			}
			// fmt.Println(g.textureIdx)
		}
	}

	return false
}
