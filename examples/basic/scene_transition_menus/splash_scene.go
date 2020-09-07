package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene

	pretendWorkCnt  float64
	pretendWorkSpan float64

	delay          api.IDelay
	tweenOntoStage *gween.Tween
	tweenOffStage  *gween.Tween
}

func newBasicSplashScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer) (api.INode, error) {
	o := new(sceneSplash)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 100.0

	o.delay = nodes.NewDelay()

	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	textureNode, err := extras.NewBitmapFont9x9Node("Ranger", textureAtlas, fontRenderer, world, o)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(25)
	textureNode.SetPosition(-100.0, 0.0)

	tn := textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Splash Scene")
	tn.SetColor(color.NewPaletteInt64(color.White).Array())

	return o, nil
}

func (s *sceneSplash) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s, color.NewPaletteInt64(color.LightGray))
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newBasicSplashLayer("Game Layer", world, s)

	return nil
}

func (s *sceneSplash) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		s.setState("Update: ", api.SceneOnStage)
	case api.SceneOnStage:
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
		}

		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningOut:
		// Update animation
		value, isFinished := s.tweenOffStage.Update(float32(msPerUpdate))

		// Update animation properties
		if isFinished {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.SetPosition(value, s.Position().Y())
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneSplash) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

func (s *sceneSplash) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		// Create an animation that drags the scene onto the stage
		// in the +X direction (enters from right)
		vrs := s.World().Properties().Window.DeviceRes
		s.tweenOffStage = gween.New(s.Position().X(), float32(vrs.Width), s.TransitionDuration(), ease.OutCubic)
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterScene called when a node is entering the stage
func (s *sceneSplash) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneSplash EnterNode")
	s.SetVisible(true)
	man.RegisterTarget(s)
}

// ExitScene called when a node is exiting stage
func (s *sceneSplash) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneSplash ExitNode")
	man.UnRegisterTarget(s)
	s.setState("ExitScene: ", api.SceneOffStage)
	return false
}
