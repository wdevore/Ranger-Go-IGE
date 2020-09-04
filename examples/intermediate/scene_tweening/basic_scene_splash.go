package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene

	// Tanema's framework
	tween *gween.Tween
}

func newBasicSplashScene(name string, world api.IWorld, replacement api.INode) (api.INode, error) {
	o := new(sceneSplash)
	o.Initialize(name)
	o.SetReplacement(replacement)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *sceneSplash) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes
	bg := newBackgroundNode("Background", world, s)
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newBasicGameLayer("Game Layer", world, s)

	newOverlayLayer("Overlay Layer", world, s)

	s.tween = gween.New(float32(-dvr.Width), 0.0, 2000, ease.OutExpo)

	return nil
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneSplash) TransitionAction() int {
	// Basically this scene never transitions to any node.
	return api.SceneNoAction
}

// Update updates the time properties of a node.
func (s *sceneSplash) Update(msPerUpdate, secPerUpdate float64) {
	value, isFinished := s.tween.Update(float32(msPerUpdate))

	if !isFinished {
		s.SetPosition(value, s.Position().Y())
	}
	// else {
	// 	s.tween.Reset()
	// }
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneSplash) EnterNode(man api.INodeManager) {
	man.RegisterTarget(s)
	man.RegisterEventTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneSplash) ExitNode(man api.INodeManager) {
	man.UnRegisterTarget(s)
	man.UnRegisterEventTarget(s)
}

// -----------------------------------------------------
// IO events
// -----------------------------------------------------

func (s *sceneSplash) Handle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeKeyboard {
		// fmt.Println(event.GetKeyScan())
		switch event.GetKeyScan() {
		case 65: // A
			if event.GetState() == 1 {
				// Transition out
				s.tween.Reset()
				dvr := s.World().Properties().Window.DeviceRes
				s.tween = gween.New(0.0, float32(-dvr.Width), 2000, ease.OutExpo)
			}
		case 84: // T
			if event.GetState() == 1 {
				// Transition in
				s.tween.Reset()
				dvr := s.World().Properties().Window.DeviceRes
				s.tween = gween.New(float32(-dvr.Width), 0.0, 2000, ease.OutExpo)
			}
		}
	}

	return false
}
