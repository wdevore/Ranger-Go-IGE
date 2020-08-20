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

	// dvr := world.Properties().Window.DeviceRes

	// bg := newBackgroundNode("Background", world, s)
	// bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	game, err := newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	newOverlayLayer("Overlay Layer", world, game, s)

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