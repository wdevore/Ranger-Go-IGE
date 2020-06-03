package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene
	nodes.Transition
}

func newBasicSplashScene(name string, replacement api.INode) api.INode {
	o := new(sceneSplash)
	o.Initialize(name)
	o.SetReplacement(replacement)
	return o
}

func (s *sceneSplash) Build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := world.Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s)
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	// newBasicGameLayer("Game Layer", world, s)
	_, err := newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	return nil
}

func (s *sceneSplash) Use() {
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneSplash) TransitionAction() int {
	// Basically this scene never transitions to any node.
	return api.SceneNoAction
}
