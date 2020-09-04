package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type sceneExit struct {
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

func newBasicExitScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer, replacement api.INode) (api.INode, error) {
	o := new(sceneExit)
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

	textureNode, err := extras.NewBitmapFont9x9Node("ExitText", textureAtlas, fontRenderer, world, o)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(25)
	textureNode.SetPosition(-100.0, 0.0)

	tn := textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Exit Scene. Goodbye...")
	tn.SetColor(color.NewPaletteInt64(color.White).Array())

	return o, nil
}

func (s *sceneExit) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s)
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newBasicExitGameLayer("Game Layer", world, s)

	return nil
}

func (s *sceneExit) Update(msPerUpdate, secPerUpdate float64) {
	switch s.currentState {
	case api.SceneOffStage:
		return
	case api.SceneOnStage:
		// fmt.Println("sceneBoot Update busy")
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
		}

		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningIn:
		// fmt.Println("sceneBoot Update trans IN")
		if s.transitionInCnt > s.transitionInDelay {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.transitionInCnt += msPerUpdate
		// Update animation properties
	case api.SceneTransitioningOut:
		// fmt.Println("sceneBoot Update trans OUT ", s.transitionOutCnt)
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

func (s *sceneExit) State() (current, previous int) {
	return s.currentState, s.previousState
}

func (s *sceneExit) setState(header string, state int) {
	s.previousState = s.currentState
	s.currentState = state
	nodes.ShowState(header, s, "")
}

func (s *sceneExit) Notify(state int) {
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
func (s *sceneExit) EnterNode(man api.INodeManager) {
	fmt.Println("sceneExit EnterNode")
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneExit) ExitNode(man api.INodeManager) {
	fmt.Println("sceneExit ExitNode")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)
}
