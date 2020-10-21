package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

func main() {
	engine, err := engine.Construct("../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	// -----------------------------------------------------
	// Create any Atlases the game/example needs.
	// This example needs the provided basic Static-Mono atlas.
	// You are free to create your own Atlases btw.
	// -----------------------------------------------------
	monoAtlas := world.GetAtlas(api.MonoAtlasName)
	if monoAtlas == nil {
		monoAtlas = atlas.NewStaticMonoAtlas(world)
		world.AddAtlas(api.MonoAtlasName, monoAtlas)
	}

	// -----------------------------------------------------
	// Create any Atlases the game/example needs.
	// This example needs the provided basic Static-Mono atlas.
	// You are free to create your own Atlases btw.
	// -----------------------------------------------------
	dynoAtlas := world.GetAtlas(api.DynamicMonoAtlasName)
	if dynoAtlas == nil {
		dynoAtlas = atlas.NewDynamicMonoAtlas(world)
		world.AddAtlas(api.DynamicMonoAtlasName, dynoAtlas)
	}

	splash, err := newBasicSplashScene("Splash", world)
	if err != nil {
		panic(err)
	}
	world.Push(splash)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := extras.NewBasicBootScene("Boot")

	world.Push(boot)

	// -----------------------------------------------------
	// We don't need to burn the MonoAtlas because the config.json commands
	// the engine to supply a background and as such the atlas will be
	// burnt automagically.
	// Note: that most of the time "you" will supplying your own
	// backgrounds and as such you will need to remember to burn
	// the atlas.
	// -----------------------------------------------------

	// We do need to burn the dynamic atlas
	err = dynoAtlas.Burn()
	if err != nil {
		panic(err)
	}

	// And finally we can start the game.
	engine.Begin()
	if err != nil {
		panic(err)
	}
}
