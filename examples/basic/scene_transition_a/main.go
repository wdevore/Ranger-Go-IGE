package main

import (
	"log"

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

	exitScene, err := newBasicExitScene("Exit", world, fontTextureRenderer, nil)
	if err != nil {
		panic(err)
	}
	engine.Push(exitScene)

	splash, err := newBasicSplashScene("Splash", world, fontTextureRenderer, nil)
	if err != nil {
		panic(err)
	}
	engine.Push(splash)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := NewBasicBootScene("Boot", world, fontTextureRenderer, nil)

	// nodes.PrintTree(splash)

	engine.Push(boot)

	engine.Begin()

}
