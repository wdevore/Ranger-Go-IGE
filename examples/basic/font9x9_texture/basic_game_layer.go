package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
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

	g.addFont(world)

	return nil
}

func (g *gameLayer) addFont(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	g.textureNode, err = custom.NewDynamicTextureNode("Font9x9", 0, textureMan, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(50)
	g.textureNode.SetPosition(0.0, 0.0)

	indexes := []int{}
	for i := 0; i < 94; i++ {
		indexes = append(indexes, i)
	}

	tn := g.textureNode.(*custom.DynamicTextureNode)
	tn.SetIndexes(indexes)
	tn.Populate()

	// Use render graphic to bind image
	renG := world.GetRenderGraphic(api.TextureRenderGraphic)
	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	renG.ConstructWithImage(textureAtlas.AtlasImage(), true, world.ShapeAtlas())
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
				g.textureIdx = (g.textureIdx - 1) % 94
				if g.textureIdx < 0 {
					g.textureIdx = 94 - 1
				}
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			case 262: // Right
				g.textureIdx = (g.textureIdx + 1) % 94
				tn := g.textureNode.(*custom.DynamicTextureNode)
				tn.SelectCoordsByIndex(g.textureIdx)
			}
			// fmt.Println(g.textureIdx)
		}
	}

	return false
}
