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

	angle    float64
	msCnt    float64
	timeSpan float64

	zbar  api.INode
	ozbar api.INode

	shipNode         api.INode
	textureNode      api.INode
	textureNodeAlpha api.INode
	textureIdx       int
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

	g.timeSpan = 500.0
	textureMan := world.TextureManager()

	g.addShip(world)
	// g.addFont(world)
	// g.addDynText(world)

	// ---------------------------------------------------------
	// Bind atlas images for text above
	// textureAtlas := textureMan.GetAtlasByName("Font9x9")
	// renG := world.GetRenderGraphic(api.TextureRenderGraphic)
	// renG.ConstructWithImage(textureAtlas.AtlasImage(), false, world.ShapeAtlas())

	textureAtlasShip := textureMan.GetAtlasByName("StarShip")
	renG2 := world.GetRenderGraphic(api.TextureRenderGraphic)
	renG2.ConstructWithImage(textureAtlasShip.AtlasImage(), false)
	// renG2 := world.GetRenderGraphic(api.TextureRenderGraphic)
	// renG2.ConstructWithImage(textureAtlas2.AtlasImage(), false)

	return nil
}

func (g *gameLayer) addFont(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.textureNode, err = extras.NewBitmapFont9x9Node("Ranger", "Font9x9", textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(25)
	g.textureNode.SetPosition(20.0, 20.0)
	g.textureNode.SetRotation(45.0 * maths.DegreeToRadians)

	tn := g.textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Ranger is a Go!")
	tn.SetColor(color.NewPaletteInt64(color.LightPink).Array())
	tn.Populate()
}

func (g *gameLayer) addDynText(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.textureNodeAlpha, err = extras.NewDynamicTextureNode("Font9x9", api.Texture2RenderGraphic, 0, textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNodeAlpha.SetScale(500)
	g.textureNodeAlpha.SetPosition(0.0, -150.0)
	g.textureNodeAlpha.SetRotation(-20.0 * maths.DegreeToRadians)

	indexes := []int{}
	for i := 0; i < 94; i++ {
		indexes = append(indexes, i)
	}

	tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	tn.SetIndexes(indexes)
	c := color.NewPaletteInt64(color.PanSkin)
	c.SetAlpha(0.5)
	tn.SetColor(c.Array())
	tn.Populate(0)

	textureNode, err := extras.NewBitmapFont9x9Node("StarCastle", "Font9x9", textureMan, world, g)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(25)
	textureNode.SetPosition(0.0, -50.0)

	tn2 := textureNode.(*extras.BitmapFont9x9Node)
	tn2.SetText("Star Castle")
	tn2.Populate()
}

func (g *gameLayer) addShip(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.shipNode, err = extras.NewDynamicTextureNode("StarShip", api.TextureRenderGraphic, 0, textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.shipNode.SetScale(300)
	g.shipNode.SetPosition(-200.0, 0.0)

	indexes := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 29, 30, 31, 32, 33, 34, 35,
	}

	tn := g.shipNode.(*extras.DynamicTextureNode)
	tn.SetIndexes(indexes)
	tn.SetColor(color.NewPaletteInt64(color.Transparent).Array())
	tn.Populate(1)
}

func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	// g.textureNode.SetRotation(maths.DegreeToRadians * g.angle)
	// g.angle -= 0.5

	// if g.msCnt > g.timeSpan {
	// 	g.msCnt = 0.0
	// 	g.incTextureID()
	// }
	// g.msCnt += msPerUpdate
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
				g.decTextureID()
			case 262: // Right
				g.incTextureID()
			}
			// fmt.Println(g.textureIdx)
		}
	}

	return false
}

func (g *gameLayer) incTextureID() {
	g.textureIdx = (g.textureIdx + 1) % 35
	tn := g.shipNode.(*extras.DynamicTextureNode)
	tn.SelectCoordsByIndex(g.textureIdx)

	// g.textureIdx = (g.textureIdx + 1) % 94
	// tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	// tn.SelectCoordsByIndex(g.textureIdx)
}

func (g *gameLayer) decTextureID() {
	g.textureIdx = (g.textureIdx - 1) % 35
	if g.textureIdx < 0 {
		g.textureIdx = 35 - 1
	}

	tn := g.shipNode.(*extras.DynamicTextureNode)
	tn.SelectCoordsByIndex(g.textureIdx)

	// g.textureIdx = (g.textureIdx - 1) % 94
	// if g.textureIdx < 0 {
	// 	g.textureIdx = 94 - 1
	// }
	// tn := g.textureNodeAlpha.(*extras.DynamicTextureNode)
	// tn.SelectCoordsByIndex(g.textureIdx)
}
