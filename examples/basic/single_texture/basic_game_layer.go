package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras"
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

	return nil
}

func (g *gameLayer) addShip(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.textureNode, err = extras.NewDynamicTextureNode("StarShip", api.TextureRenderGraphic, 0, textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(300)
	g.textureNode.SetPosition(0.0, 0.0)

	indexes := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 29, 30, 31, 32, 33, 34, 35,
	}

	tn := g.textureNode.(*extras.DynamicTextureNode)
	tn.SetIndexes(indexes)
	tn.Populate(0)

	// Use render graphic to bind image
	renG := world.GetRenderGraphic(api.TextureRenderGraphic)
	textureAtlas := textureMan.GetAtlasByName("StarShip")

	renG.ConstructWithImage(textureAtlas.AtlasImage(), false)
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
			switch event.GetKeyScan() {
			case 68: // d
				tn := g.textureNode.(*extras.DynamicTextureNode)
				tn.SelectCoordsByIndex(1)
			case 70: // f
				tn := g.textureNode.(*extras.DynamicTextureNode)
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
				tn := g.textureNode.(*extras.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			case 262: // Right
				g.textureIdx = (g.textureIdx + 1) % 35
				tn := g.textureNode.(*extras.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			}
		}
	}

	return false
}
