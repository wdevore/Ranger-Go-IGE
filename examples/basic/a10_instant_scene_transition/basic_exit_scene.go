package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type sceneExit struct {
	nodes.Node
	nodes.Scene

	pretendWorkCnt  float64
	pretendWorkSpan float64

	delay api.IDelay
}

func newBasicExitScene(name string, world api.IWorld) (api.INode, error) {
	o := new(sceneExit)
	o.Initialize(name)

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.pretendWorkSpan = 1000.0

	o.delay = nodes.NewDelay()

	return o, nil
}

func (s *sceneExit) build(world api.IWorld) error {
	s.Node.Build(world)

	// This is an example of a custom background node.
	bg, err := newBackgroundNode("Background", world, s)
	if err != nil {
		return err
	}
	dvr := world.Properties().Window.DeviceRes
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	bn := bg.(*backgroundNode)
	bn.setColor(color.NewPaletteInt64(color.LighterGray))

	newBasicExitGameLayer("Game Layer", world, s)

	return nil
}

func (s *sceneExit) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		if s.delay.ReadyToTransition() {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.delay.UpdateTransition(msPerUpdate)
		// Update animation properties
	case api.SceneOnStage:
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
			s.delay.SetPauseTime(1000.0)
			s.delay.Reset()
		}

		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningOut:
		// Update animation
		if s.delay.ReadyToTransition() {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.delay.UpdateTransition(msPerUpdate)
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneExit) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

func (s *sceneExit) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Configure animation properties for entering the stage.
		s.delay.SetPauseTime(1000.0)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneExit) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneExit EnterNode")
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneExit) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneExit ExitNode")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)
	return false
}
