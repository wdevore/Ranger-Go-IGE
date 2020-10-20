package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// This scene slides upwards onto the stage and downwards off the stage
type gameMenu struct {
	nodes.Node
	nodes.Scene

	atlas api.IAtlasX

	delay          api.IDelay
	tweenOntoStage *gween.Tween
	tweenOffStage  *gween.Tween

	enterExitState int
}

func newGameScene(name string, atlas api.IAtlasX, world api.IWorld) (api.INode, error) {
	o := new(gameMenu)
	o.Initialize(name)
	o.atlas = atlas

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.enterExitState = 0
	o.delay = nodes.NewDelay()

	return o, nil
}

func (s *gameMenu) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	// This is an example of a custom background node.
	bg, err := newBackgroundNode("Background", world, s)
	if err != nil {
		return err
	}
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	bn := bg.(*backgroundNode)
	bn.setColor(color.NewPaletteInt64(color.DarkGray))

	newGameLayer("Game Layer", s.atlas, world, s)

	return nil
}

func (s *gameMenu) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		value, isFinished := s.tweenOntoStage.Update(float32(msPerUpdate))

		// Update animation properties
		if isFinished {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.SetPosition(s.Position().X(), value)
	case api.SceneOnStage:
		// Stay on stage until user makes a choice.
	case api.SceneTransitioningOut:
		// Update animation
		value, isFinished := s.tweenOffStage.Update(float32(msPerUpdate))

		// Update animation properties
		if isFinished {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.SetPosition(s.Position().X(), value)
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *gameMenu) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

func (s *gameMenu) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Animate from bottom to top
		vrs := s.World().Properties().Window.DeviceRes
		s.tweenOntoStage = gween.New(-float32(vrs.Height), 0.0, s.TransitionDuration(), ease.OutCubic)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		// drop downwards
		vrs := s.World().Properties().Window.DeviceRes
		s.tweenOffStage = gween.New(0.0, -float32(vrs.Height), s.TransitionDuration(), ease.OutCubic)
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterScene called when a node is entering the stage
func (s *gameMenu) EnterScene(man api.INodeManager) {
	// fmt.Println("settingsMenu EnterNode")
	vrs := s.World().Properties().Window.DeviceRes
	s.SetPosition(0.0, -float32(vrs.Height))
	s.SetVisible(true)
	man.RegisterTarget(s)
	man.RegisterEventTarget(s)
}

// ExitScene called when a node is exiting stage
func (s *gameMenu) ExitScene(man api.INodeManager) bool {
	// fmt.Println("settingsMenu ExitNode")
	man.UnRegisterTarget(s)
	man.UnRegisterEventTarget(s)
	s.setState("ExitScene: ", api.SceneOffStage)
	return false
}

func (s *gameMenu) Handle(event api.IEvent) bool {
	if s.CurrentState() != api.SceneOnStage {
		return false
	}

	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println("settingsMenu ", event.GetKeyScan())
		// fmt.Println(event)

		if event.GetState() == 1 {
			switch event.GetKeyScan() {
			case 82: // r
				//
				// Signal NM that this scene wants to transition out.
				s.setState("Notify T: ", api.SceneTransitionStartOut)
			}
		}
	}

	return false
}
