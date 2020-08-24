package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type gameLayer struct {
	nodes.Node

	angle float64
	zbar  api.INode
	ozbar api.INode

	fontTextureRenderer api.ITextureRenderer

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

	textureMan := world.TextureManager()

	g.fontTextureRenderer = rendering.NewTextureRenderer(textureMan, world.TextureShader())
	g.fontTextureRenderer.Build("Font9x9")

	g.addFont(world)

	return nil
}

func (g *gameLayer) addFont(world api.IWorld) {
	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	g.textureNode, err = extras.NewDynamicTextureNode("Font9x9", textureAtlas, g.fontTextureRenderer, world, g)
	if err != nil {
		panic(err)
	}
	g.textureNode.SetScale(50)
	g.textureNode.SetPosition(0.0, 0.0)
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
				tn.SetIndex(1)
			case 70: // f
				tn := g.textureNode.(*extras.DynamicTextureNode)
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
				tn := g.textureNode.(*extras.DynamicTextureNode)
				tn.SetIndex(g.textureIdx)
			case 262: // Right
				g.textureIdx = (g.textureIdx + 1) % 94
				tn := g.textureNode.(*extras.DynamicTextureNode)
				tn.SetIndex(g.textureIdx)
			}
			// fmt.Println(g.textureIdx)
		}
	}

	return false
}
