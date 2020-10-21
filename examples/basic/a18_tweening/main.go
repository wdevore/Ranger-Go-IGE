package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

func main() {
	engine, err := engine.Construct("../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	splash, err := newBasicSplashScene("Splash", world)
	if err != nil {
		panic(err)
	}
	world.Push(splash)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := extras.NewBasicBootScene("Boot")

	nodes.PrintTree(splash)

	world.Push(boot)

	// -----------------------------------------------------
	// We don't need to burn the MonoAtlas because the config.json commands
	// the engine to supply a background and as such the atlas will be
	// burnt automagically.
	// Note: that most of the time "you" will supplying your own
	// backgrounds and as such you will need to remember to burn
	// the atlas.
	// -----------------------------------------------------

	// And finally we can start the game.
	engine.Begin()
	if err != nil {
		panic(err)
	}
}
