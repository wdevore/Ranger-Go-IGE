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

	dvr := world.Properties().Window.DeviceRes

	atlas := world.GetAtlas(monoAtlasName)

	// -----------------------------------------------------
	// Preload any shapes the game needs.
	// This example needs two shapes.
	// -----------------------------------------------------
	vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(true, true)
	atlas.AddShape(centerSquareName, vertices, indices, mode)

	vertices, indices, mode = generators.GenerateUnitTriangleVectorShape(true)
	atlas.AddShape(triangleName, vertices, indices, mode)

	bg, err := newBackgroundNode("Background", world, s)
	if err != nil {
		return err
	}

	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	_, err = newBasicGameLayer("Game Layer", world, s)
	if err != nil {
		return err
	}

	return nil
}
