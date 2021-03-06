package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

const (
	centerSquareName = "CenteredSquareShape"
)

type sceneSplash struct {
	nodes.Node
	nodes.Scene
}

func newBasicSplashScene(name string, world api.IWorld) (api.INode, error) {
	o := new(sceneSplash)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *sceneSplash) build(world api.IWorld) error {
	s.Node.Build(world)

	atlas := world.GetAtlas(api.MonoAtlasName)

	// -----------------------------------------------------
	// Preload any shapes the game needs.
	// This example needs two shapes.
	// -----------------------------------------------------
	vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(true, true)
	atlas.(api.IStaticAtlasX).AddShape(centerSquareName, vertices, indices, mode)

	_, err := newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	return nil
}
