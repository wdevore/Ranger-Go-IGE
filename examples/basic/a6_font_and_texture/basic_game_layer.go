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

	angle    float64
	msCnt    float64
	timeSpan float64

	zbar  api.INode
	ozbar api.INode

	shipNode    api.INode
	textureNode api.INode
	textureIdx  int

	shipTextureRenderer api.ITextureRenderer
	textureShipIdx      int
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

	g.timeSpan = 1000.0

	err := g.addShip(world)
	if err != nil {
		return err
	}

	err = g.addFont(world)
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
	g.textureNode.SetScale(50.0)
	g.textureNode.SetPosition(20.0, 20.0)
	g.textureNode.SetRotation(45.0 * maths.DegreeToRadians)
	bf := g.textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(color.NewPaletteInt64(color.GreenYellow).Array())
	bf.SetText("Ranger is a Go!")

	return nil
}

func (g *gameLayer) addShip(world api.IWorld) error {
	var err error

	g.spriteSheet = fonts.NewFont9x9SpriteSheet("StarShip", "starship_sprite_sheet_manifest.json")
	g.spriteSheet.Load("../../assets/", true)

	atlas := atlas.NewSingleTextureAtlas("StarShip", g.spriteSheet, world)
	err = atlas.Burn()
	if err != nil {
		return err
	}

	// #3 INode
	g.shipNode, err = shapes.NewBitmapNode("StarShip", atlas, world, g)
	if err != nil {
		return err
	}
	g.shipNode.SetScale(300)
	g.shipNode.SetPosition(-200.0, 0.0)

	return nil
}

func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.textureNode.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle -= 0.5

	if g.msCnt > g.timeSpan {
		g.msCnt = 0.0
		g.incTextureID()
	}
	g.msCnt += msPerUpdate
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (g *gameLayer) EnterNode(man api.INodeManager) {
	man.RegisterTarget(g)
	man.RegisterEventTarget(g)
}

// ExitNode called when a node is exiting stage
func (g *gameLayer) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(g)
	man.UnRegisterEventTarget(g)
}

func (g *gameLayer) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		// fmt.Println(event)

		if event.GetState() == 1 || event.GetState() == 2 {
			switch event.GetKeyScan() {
			case 68: // d
			case 70: // f
			case 90: // z
			case 65: // a
			case 264: // Down
				g.decTextureID()
			case 265: // Up
				g.incTextureID()
			case 263: // Left
				g.decShipTextureID()
			case 262: // Right
				g.incShipTextureID()
			}
			// fmt.Println(g.textureIdx)
		}
	}

	return false
}

func (g *gameLayer) incShipTextureID() {
	g.textureShipIdx = (g.textureShipIdx + 1) % 35
	tn := g.shipNode.(*shapes.BitmapNode)
	tn.SetIndex(g.textureShipIdx)
}

func (g *gameLayer) decShipTextureID() {
	g.textureShipIdx = (g.textureShipIdx - 1) % 35
	if g.textureShipIdx < 0 {
		g.textureShipIdx = 35 - 1
	}

	tn := g.shipNode.(*shapes.BitmapNode)
	tn.SetIndex(g.textureShipIdx)
}

func (g *gameLayer) incTextureID() {
	g.textureIdx = (g.textureIdx + 1) % 94
	tn := g.textureNode.(*shapes.BitmapFont9x9Node)
	tn.SetIndex(g.textureIdx)
}

func (g *gameLayer) decTextureID() {
	g.textureIdx = (g.textureIdx - 1) % 94
	if g.textureIdx < 0 {
		g.textureIdx = 94 - 1
	}
	tn := g.textureNode.(*shapes.BitmapFont9x9Node)
	tn.SetIndex(g.textureIdx)
}
