package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// This example shows the most basic of scene transitions.
// The transitions are instant which means there is no animations
// onto or off of the stage.
// This also means that the boot scene will be covered by the
// Exit scene while the Boot scene is exiting, so it will appear
// as if the boot scene was only visible for 2 secs instead of 3.

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
	world.Push(exitScene)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := NewBasicBootScene("Boot", world, fontTextureRenderer)

	// nodes.PrintTree(splash)

	world.Push(boot)

	engine.Begin()

}
