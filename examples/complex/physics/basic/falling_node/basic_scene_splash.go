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
	nodes.Transition

	// Tanema's framework
	tween *gween.Tween
}

func newBasicSplashScene(name string, replacement api.INode) api.INode {
	o := new(sceneSplash)
	o.Initialize(name)
	o.SetReplacement(replacement)
	return o
}

func (s *sceneSplash) Build(world api.IWorld) error {
	s.Node.Build(world)

	game, err := newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	game.SetScale(5.0)
	newOverlayLayer("Overlay Layer", world, game, s)

	return nil
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneSplash) TransitionAction() int {
	// Basically this scene never transitions to any node.
	return api.SceneNoAction
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneSplash) EnterNode(man api.INodeManager) {
	man.RegisterEventTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneSplash) ExitNode(man api.INodeManager) {
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
