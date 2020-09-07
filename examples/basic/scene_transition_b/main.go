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

	exitScene, err := newBasicExitScene("Exit", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	exitScene.SetVisible(false)

	world.Push(exitScene)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot, err := newBasicBootScene("Boot", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}

	world.Push(boot)

	engine.Begin()

}
