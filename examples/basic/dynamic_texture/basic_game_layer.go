package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type gameLayer struct {
	nodes.Node

	angle float64
	zbar  api.INode
	ozbar api.INode

	textureNode api.INode
	textureIdx  int
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

	g.addShip(world)

	g.addBar(world)

	return nil
}

func (g *gameLayer) addShip(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.textureNode, err = custom.NewDynamicTextureNode("StarShip", 0, textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(300)
	g.textureNode.SetPosition(0.0, 0.0)

	indexes := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 29, 30, 31, 32, 33, 34, 35,
	}

	tn := g.textureNode.(*custom.DynamicTextureNode)
	tn.SetIndexes(indexes)
	tn.Populate()

	// Use render graphic to bind image
	renG := world.GetRenderGraphic(api.TextureRenderGraphic)
	textureAtlas := textureMan.GetAtlasByName("StarShip")

	renG.ConstructWithImage(textureAtlas.AtlasImage(), false, world.ShapeAtlas())
}

func (g *gameLayer) addBar(world api.IWorld) {
	var err error

	xPos := float32(100.0)
	yPos := float32(100.0)
	scale := float32(200.0)

	// ---------------------------------------------------------
	g.zbar, err = custom.NewStaticZBarNode("FilledZBar", true, world, g)
	if err != nil {
		panic(err)
	}
	g.zbar.SetScale(scale)
	g.zbar.SetPosition(xPos, yPos)
	gzr := g.zbar.(*custom.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.LightOrange))

	g.ozbar, err = custom.NewStaticZBarNode("OutlineZBar", false, world, g)
	if err != nil {
		panic(err)
	}
	g.ozbar.SetScale(scale)
	g.ozbar.SetPosition(xPos, yPos)
	gzr = g.ozbar.(*custom.StaticZBarNode)
	gzr.SetColor(color.NewPaletteInt64(color.White))
}

// Update updates the time properties of a node.
func (g *gameLayer) Update(msPerUpdate, secPerUpdate float64) {
	g.zbar.SetRotation(maths.DegreeToRadians * g.angle)
	g.ozbar.SetRotation(maths.DegreeToRadians * g.angle)
	g.angle += 1.25
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
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(1)
			case 70: // f
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(0)
			case 90: // z
			case 65: // a
			case 83: // s
			case 82: // R
			case 263: // Left
				g.textureIdx = (g.textureIdx - 1) % 35
				if g.textureIdx < 0 {
					g.textureIdx = 34
				}
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			case 262: // Right
				g.textureIdx = (g.textureIdx + 1) % 35
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			}
		}
	}

	return false
}
