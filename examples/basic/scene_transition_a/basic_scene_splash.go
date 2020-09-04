package main

import (
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

	transition api.ITransition
}

func newBasicSplashScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer, replacement api.INode) (api.INode, error) {
	o := new(sceneSplash)
	o.Initialize(name)
	o.SetReplacement(replacement)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 1000.0

	o.transition = nodes.NewTransition()

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

	bg := newBackgroundNode("Background", world, s)
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newBasicGameLayer("Game Layer", world, s)

	return nil
}

func (s *sceneSplash) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneOnStage:
		// fmt.Println("sceneSplash Update busy")
		if s.pretendWorkCnt > s.pretendWorkSpan {
			s.pretendWorkCnt = 0.0
			s.setState("Update: ", api.SceneTransitioningOut)
			s.transition.SetPauseTime(1000.0)
			s.transition.Reset()
		}
		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningIn:
		// fmt.Println("sceneSplash Update trans IN")
		if s.transition.ReadyToTransition() {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.transition.UpdateTransition(msPerUpdate)
		// Update animation properties
	case api.SceneTransitioningOut:
		// Update animation
		if s.transition.ReadyToTransition() {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.transition.UpdateTransition(msPerUpdate)
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
	// Technically the boot scene never cares about Notify
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Configure animation properties for entering the stage.
		s.transition.SetPauseTime(1000.0)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneSplash) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneSplash EnterNode")
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneSplash) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneSplash ExitNode")
	man.UnRegisterTarget(s)
	return false
}
