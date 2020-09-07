package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

func main() {
	engine, err := engine.Construct("../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	textureMan := world.TextureManager()

	// ------------------------------------------------------------
	// Load any textures we need
	// ------------------------------------------------------------
	textureMan.AddAtlas("Font9x9", "../../../assets/images/atlas/", "font9x9_texture_manifest.json", true)

	fontTextureRenderer := rendering.NewTextureRenderer(textureMan, world.TextureShader())
	fontTextureRenderer.Build("Font9x9")

	transitionDuration := float32(500.0)

	exitScene, err := newBasicExitScene("Exit", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	sceneE := exitScene.(api.IScene)
	sceneE.SetTransitionDuration(transitionDuration)
	exitScene.SetVisible(false)
	world.Push(exitScene)

	splash, err := newBasicSplashScene("Splash", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	sceneS := splash.(api.IScene)
	sceneS.SetTransitionDuration(transitionDuration)
	splash.SetVisible(false)
	world.Push(splash)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot, err := NewBasicBootScene("Boot", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	sceneB := boot.(api.IScene)
	sceneB.SetTransitionDuration(transitionDuration)

	// nodes.PrintTree(splash)

	world.Push(boot)

	engine.Begin()

}
