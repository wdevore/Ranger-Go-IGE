package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type sceneMenu struct {
	nodes.Node
	nodes.Scene

	delay          api.IDelay
	tweenOntoStage *gween.Tween
	tweenOffStage  *gween.Tween

	// 0 = Initial entry onto the stage
	// 1 = Exit stage to another scene
	enterExitState int

	subMenus []api.INode
}

func newMenuScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer) (api.INode, error) {
	o := new(sceneMenu)
	o.Initialize(name)

	if err := o.build(world, fontRenderer); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.enterExitState = 0
	o.delay = nodes.NewDelay()

	return o, nil
}

func (s *sceneMenu) build(world api.IWorld, fontRenderer api.ITextureRenderer) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s, color.NewPaletteInt64(color.LighterGray))
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newMenuLayer("Menu Layer", world, fontRenderer, s)

	return nil
}

func (s *sceneMenu) AddSubMenu(subMenu api.INode) {
	s.subMenus = append(s.subMenus, subMenu)
}

func (s *sceneMenu) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		// Update animation properties
		value, isFinished := s.tweenOntoStage.Update(float32(msPerUpdate))

		if isFinished {
			s.setState("Update: ", api.SceneOnStage)
		}

		if s.enterExitState == 0 {
			s.SetPosition(value, s.Position().Y())
		} else {
			s.SetPosition(s.Position().X(), value)
		}
	case api.SceneOnStage:
		// Stay on stage until user makes a choice.
	case api.SceneTransitioningOut:
		// Update animation
		value, isFinished := s.tweenOffStage.Update(float32(msPerUpdate))

		if isFinished {
			s.setState("Update: ", api.SceneExitedStage)
		}
		if s.enterExitState == 0 {
			s.SetPosition(value, s.Position().Y())
		} else {
			s.SetPosition(s.Position().X(), value)
		}
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneMenu) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

func (s *sceneMenu) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		vrs := s.World().Properties().Window.DeviceRes
		if s.enterExitState == 0 {
			// Create an animation that drags the scene onto the stage
			// in the +X direction (enters from right)
			s.tweenOntoStage = gween.New(-float32(vrs.Width), 0.0, s.TransitionDuration(), ease.OutCubic)
		} else {
			// Vertical animation that comes from the top and moves downward
			s.tweenOntoStage = gween.New(float32(vrs.Height), 0.0, s.TransitionDuration(), ease.OutCubic)
		}

		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		vrs := s.World().Properties().Window.DeviceRes
		if s.enterExitState == 0 {
			// Create an animation that drags the scene onto the stage
			// in the +X direction (enters from right)
			s.tweenOffStage = gween.New(0.0, float32(vrs.Width), s.TransitionDuration(), ease.OutCubic)
		} else {
			// The user has selected a menu choice
			// Create a Vertical animation moving upwards off stage
			s.tweenOffStage = gween.New(0.0, float32(vrs.Height), s.TransitionDuration(), ease.OutCubic)
		}
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterScene called when a node is entering the stage
func (s *sceneMenu) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneMenu EnterScene")
	s.SetVisible(true)
	vrs := s.World().Properties().Window.DeviceRes
	if s.enterExitState == 0 {
		s.SetPosition(-float32(vrs.Width), 0.0)
	} else {
		s.SetPosition(0.0, float32(vrs.Height))
	}
	man.RegisterTarget(s)
	man.RegisterEventTarget(s)
}

// ExitScene called when a node is exiting stage
func (s *sceneMenu) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneMenu ExitScene")
	man.UnRegisterTarget(s)
	man.UnRegisterEventTarget(s)
	s.setState("ExitScene: ", api.SceneOffStage)
	return false
}

func (s *sceneMenu) Handle(event api.IEvent) bool {
	if s.CurrentState() != api.SceneOnStage {
		return false
	}

	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println("sceneMenu", event.GetKeyScan())
		// fmt.Println(event)

		if event.GetState() == 1 {
			switch event.GetKeyScan() {
			case 49: // 1 Settings
				s.enterExitState = 1
				// Push this Menu scene first
				s.World().Push(s)
				// Push Settings scene next
				for _, subMenu := range s.subMenus {
					if subMenu.Name() == "Settings" {
						s.World().Push(subMenu)
						break
					}
				}
				// Signal NM that this scene wants to transition out.
				s.setState("Notify T: ", api.SceneTransitionStartOut)
				return true
			case 50: // 2 HighScore
				s.enterExitState = 1
				// Push this Menu scene first
				s.World().Push(s)
				// Push Settings scene next
				for _, subMenu := range s.subMenus {
					if subMenu.Name() == "Highscore" {
						s.World().Push(subMenu)
						break
					}
				}
				// Signal NM that this scene wants to transition out.
				s.setState("Notify T: ", api.SceneTransitionStartOut)
				return true
			case 51: // 3 Game
				s.enterExitState = 1
				// Push this Menu scene first
				s.World().Push(s)
				// Push Settings scene next
				for _, subMenu := range s.subMenus {
					if subMenu.Name() == "Game" {
						s.World().Push(subMenu)
						break
					}
				}
				// Signal NM that this scene wants to transition out.
				s.setState("Notify T: ", api.SceneTransitionStartOut)
				return true
			case 88: // x
				s.enterExitState = 0
				s.setState("Handle: ", api.SceneTransitionStartOut)
			}
		}
	}

	return false
}
