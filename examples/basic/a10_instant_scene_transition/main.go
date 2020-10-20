package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
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

	// Add Atlas to the world so Scenes/Layers can obtain access to the atlas.
	world.AddAtlas(api.MonoAtlasName, monoAtlas)

	exitScene, err := newBasicExitScene("Exit", world)
	if err != nil {
		panic(err)
	}
	world.Push(exitScene)

	// This example uses the super basic Boot scene that does absolutely
	// nothing by exiting immediately.
	boot, err := newBasicBootScene("Boot", world)
	if err != nil {
		panic(err)
	}

	world.Push(boot)

	// -----------------------------------------------------
	// Now that Scene and Layers have added Shapes to the
	// Atlas we can now "Shake and Bake" it via the Burn().
	// -----------------------------------------------------
	err = monoAtlas.Burn()
	if err != nil {
		panic(err)
	}

	// And finally we can start the game.
	engine.Begin()
	if err != nil {
		panic(err)
	}
}
