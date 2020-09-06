package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

// Note: this is a very basic boot Node used pretty much for just
// engine development. You should actually supply your own boot node,
// and example can be found in the examples folder.
type sceneBoot struct {
	nodes.Node
	nodes.Scene

	delay         api.IDelay
	tweenOffStage *gween.Tween

	fontTextureRenderer api.ITextureRenderer

	pretendWorkCnt  float64
	pretendWorkSpan float64

	scanCnt   float64
	scanDelay float64

	textureNode api.INode

	dotScale float32
	dots     []api.INode
	colors   []api.IPalette
}

// NewBasicBootScene returns an IScene node of base type INode
func NewBasicBootScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer) (api.INode, error) {
	o := new(sceneBoot)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 1000.0
	o.scanDelay = 75
	o.dotScale = 15.0

	o.delay = nodes.NewDelay()

	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	o.textureNode, err = extras.NewBitmapFont9x9Node("Ranger", textureAtlas, fontRenderer, world, o)
	if err != nil {
		panic(err)
	}
	o.textureNode.SetScale(50)
	o.textureNode.SetPosition(-280.0, 0.0)

	tn := o.textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Loading")
	tn.SetColor(color.NewPaletteInt64(color.LightOrange).Array())

	o.colors = []api.IPalette{
		color.NewPaletteInt64(0xFBD872FF),
		color.NewPaletteInt64(0xFFC845FF),
		color.NewPaletteInt64(0xFFB81CFF),
		color.NewPaletteInt64(0xC69214FF),
		color.NewPaletteInt64(0xAD841FFF),
	}

	o.buildScanThingy(world)

	return o, nil
}

func (s *sceneBoot) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s, color.NewPaletteInt64(color.DarkGray))
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	return nil
}

func (s *sceneBoot) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		if s.delay.ReadyToTransition() {
			tn := s.textureNode.(*extras.BitmapFont9x9Node)
			tn.SetText("OnStage")
			s.setState("Update: ", api.SceneOnStage)
		}
		s.delay.UpdateTransition(msPerUpdate)
		// Update animation properties
	case api.SceneOnStage:
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
		}
		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningOut:
		// Transitioning out is nothing more than applying a transform
		// to the scene's position and/or rotation and/or scale.
		// This example animates only the position.
		value, isFinished := s.tweenOffStage.Update(float32(msPerUpdate))

		if isFinished {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.SetPosition(value, s.Position().Y())
	}

	s.animate(msPerUpdate)
}

func (s *sceneBoot) animate(msPerUpdate float64) {
	if s.scanCnt > s.scanDelay {
		s.scanCnt = 0.0
		// Shift colors
		c := s.colors[4]
		s.colors[4] = s.colors[3]
		s.colors[3] = s.colors[2]
		s.colors[2] = s.colors[1]
		s.colors[1] = s.colors[0]
		s.colors[0] = c

		for i, dot := range s.dots {
			gol2 := dot.(*extras.StaticSquareNode)
			gol2.SetColor(s.colors[i])
		}
	}
	s.scanCnt += msPerUpdate
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneBoot) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Create an animation that drags the scene off stage
		// in the +X direction (leaves moving to the right)
		vrs := s.World().Properties().Window.DeviceRes
		s.tweenOffStage = gween.New(s.Position().X(), float32(vrs.Width), s.TransitionDuration(), ease.OutCubic)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

func (s *sceneBoot) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterScene called when a node is entering the stage
func (s *sceneBoot) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneboot EnterNode")
	man.RegisterTarget(s)
}

// ExitScene called when a node is exiting stage.
// Return true if this node is to be "repooled" to avoid
// being destroyed.
func (s *sceneBoot) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneboot exit")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)

	return false
}

func (s *sceneBoot) buildScanThingy(world api.IWorld) {
	x := float32(75.0)
	for i := 0; i < 5; i++ {
		dot, _ := extras.NewStaticSquareNode("FilledSqr", true, true, world, s)
		s.dots = append(s.dots, dot)
		dot.SetScale(s.dotScale)
		dot.SetPosition(x, -10.0)
		gol2 := dot.(*extras.StaticSquareNode)
		gol2.SetColor(s.colors[i])
		x += s.dotScale
	}
}
