package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene
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
