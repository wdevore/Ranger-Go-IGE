package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

const (
	centerSquareName = "CenteredSquareShape"
	triangleName     = "TriangleShape"
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

	// dvr := world.Properties().Window.DeviceRes

	atlas := world.GetAtlas(api.MonoAtlasName).(api.IStaticAtlasX)

	// -----------------------------------------------------
	// Preload any shapes the game needs.
	// -----------------------------------------------------
	vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(true, true)
	atlas.AddShape(api.CenteredFilledSquareShapeName, vertices, indices, mode)

	vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(true)
	atlas.AddShape(api.FilledTriangleShapeName, vertices, indices, mode)

	// In this example, the local config.json file specifies a default
	// solid background via BackgroundColor and ClearStyle="SingleColor"
	// The Engine.Begin appends a background node to the UnderLay node.

	_, err := newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	return nil
}
