package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

const (
	monoAtlasName = "MonoAtlas"
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
	monoAtlas := atlas.NewStaticMonoAtlas(world)

	// Add it to the world so Scenes/Layers can obtain access to the atlas.
	world.AddAtlas(monoAtlasName, monoAtlas)

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
	// Finally Shake and Bake any atlases we created
	// -----------------------------------------------------
	err = monoAtlas.Burn()
	if err != nil {
		panic(err)
	}

	nodes.PrintTree(splash)

	engine.Begin()

}
