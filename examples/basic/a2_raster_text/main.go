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
	// This example only needs the provided basic Static-Mono atlas.
	// You are free to create your own Atlases btw.
	// -----------------------------------------------------
	monoAtlas := world.GetAtlas(api.MonoAtlasName)
	if monoAtlas == nil {
		monoAtlas = atlas.NewStaticMonoAtlas(world)
		world.AddAtlas(api.MonoAtlasName, monoAtlas)
	}

	// -----------------------------------------------------
	// Setup scenes and layers of the game.
	// -----------------------------------------------------
	splash, err := newBasicSplashScene("Splash", world)
	if err != nil {
		panic(err)
	}
	world.Push(splash)

	// This example uses the super basic Boot scene that does absolutely
	// nothing by exiting immediately.
	boot := extras.NewBasicBootScene("Boot")
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
	err = engine.Begin()
	if err != nil {
		panic(err)
	}
}
