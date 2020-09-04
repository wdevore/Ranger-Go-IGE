package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene
	// nodes.Transition

	pretendWorkCnt  float64
	pretendWorkSpan float64

	currentState, previousState int
	transitionInCnt             float64
	transitionInDelay           float64
	transitionOutCnt            float64
	transitionOutDelay          float64
}

func newBasicSplashScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer, replacement api.INode) (api.INode, error) {
	o := new(sceneSplash)
	o.Initialize(name)
	o.SetReplacement(replacement)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.currentState = api.SceneOffStage
	o.previousState = o.currentState
	o.pretendWorkSpan = 1000.0

	o.transitionInDelay = 1000.0
	o.transitionOutDelay = 1000.0

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
	switch s.currentState {
	case api.SceneOffStage:
		return
	case api.SceneOnStage:
		// fmt.Println("sceneSplash Update busy")
		if s.pretendWorkCnt > s.pretendWorkSpan {
			s.pretendWorkCnt = 0.0
			s.setState("Update: ", api.SceneTransitioningOut)
		}
		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningIn:
		// fmt.Println("sceneSplash Update trans IN")
		if s.transitionInCnt > s.transitionInDelay {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.transitionInCnt += msPerUpdate
		// Update animation properties
	case api.SceneTransitioningOut:
		// Update animation
		if s.transitionOutCnt > s.transitionOutDelay {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.transitionOutCnt += msPerUpdate
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneSplash) State() (current, previous int) {
	return s.currentState, s.previousState
}

func (s *sceneSplash) setState(header string, state int) {
	s.previousState = s.currentState
	s.currentState = state
	nodes.ShowState(header, s, "")
}

func (s *sceneSplash) Notify(state int) {
	// Technically the boot scene never cares about Notify
	s.setState("Notify: ", state)

	switch s.currentState {
	case api.SceneTransitionStartIn:
		// Configure animation properties for entering the stage.
		s.setState("Notify T: ", api.SceneTransitioningIn)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneSplash) EnterNode(man api.INodeManager) {
	fmt.Println("sceneSplash EnterNode")
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneSplash) ExitNode(man api.INodeManager) {
	fmt.Println("sceneSplash ExitNode")
	man.UnRegisterTarget(s)
}
