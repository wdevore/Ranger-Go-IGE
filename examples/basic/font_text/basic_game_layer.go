package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type gameLayer struct {
	nodes.Node

	angle    float64
	msCnt    float64
	timeSpan float64

	zbar  api.INode
	ozbar api.INode

	shipNode         api.INode
	textureNode      api.INode
	textureNodeAlpha api.INode
	textureIdx       int

	fontTextureRenderer api.ITextureRenderer
	shipTextureRenderer api.ITextureRenderer
	textureShipIdx      int
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

	g.timeSpan = 1000.0
	textureMan := world.TextureManager()

	g.fontTextureRenderer = rendering.NewTextureRenderer(textureMan, world.TextureShader())
	g.fontTextureRenderer.Build("Font9x9")

	g.shipTextureRenderer = rendering.NewTextureRenderer(textureMan, world.TextureShader())
	g.shipTextureRenderer.Build("StarShip")

	g.addShip(world)
	g.addFont(world)
	g.addDynText(world)

	return nil
}

func (g *gameLayer) addFont(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.textureNode, err = extras.NewBitmapFont9x9Node("Ranger", textureAtlas, g.fontTextureRenderer, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(25)
	g.textureNode.SetPosition(20.0, 20.0)
	g.textureNode.SetRotation(45.0 * maths.DegreeToRadians)

	tn := g.textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Ranger is a Go!")
	tn.SetColor(color.NewPaletteInt64(color.LightOrange).Array())
}

func (g *gameLayer) addDynText(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.textureNodeAlpha, err = extras.NewDynamicTextureNode("Ranger", textureAtlas, g.fontTextureRenderer, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNodeAlpha.SetScale(800)
	g.textureNodeAlpha.SetPosition(0.0, 0.0)
	g.textureNodeAlpha.SetRotation(-20.0 * maths.DegreeToRadians)

	tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	c := color.NewPaletteInt64(color.PanSkin)
	c.SetAlpha(0.5)
	tn.SetColor(c.Array())

	textureNode, err := extras.NewBitmapFont9x9Node("StarCastle", textureAtlas, g.fontTextureRenderer, world, g)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(25)
	textureNode.SetPosition(-25.0, -25.0)

	tn2 := textureNode.(*extras.BitmapFont9x9Node)
	tn2.SetText("Star Castle")
}

func (g *gameLayer) addShip(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("StarShip")

	g.shipNode, err = extras.NewDynamicTextureNode("StarShip", textureAtlas, g.shipTextureRenderer, world, g)
	if err != nil {
		panic(err)
	}
	g.shipNode.SetScale(300)
	g.shipNode.SetPosition(-200.0, 0.0)
}

func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.textureNodeAlpha.SetRotation(maths.DegreeToRadians * -g.angle)
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
			case 83: // s
			case 82: // R
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
	tn := g.shipNode.(*extras.DynamicTextureNode)
	tn.SetIndex(g.textureShipIdx)
}

func (g *gameLayer) decShipTextureID() {
	g.textureShipIdx = (g.textureShipIdx - 1) % 35
	if g.textureShipIdx < 0 {
		g.textureShipIdx = 35 - 1
	}

	tn := g.shipNode.(*extras.DynamicTextureNode)
	tn.SetIndex(g.textureShipIdx)
}

func (g *gameLayer) incTextureID() {
	g.textureIdx = (g.textureIdx + 1) % 94
	tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	tn.SetIndex(g.textureIdx)
}

func (g *gameLayer) decTextureID() {
	g.textureIdx = (g.textureIdx - 1) % 94
	if g.textureIdx < 0 {
		g.textureIdx = 94 - 1
	}
	tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	tn.SetIndex(g.textureIdx)
}
