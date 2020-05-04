package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
)

func main() {
	engine, err := engine.Construct("../../..")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	// Override some of the world properties for this example
	world.PropertiesOverride("config.json")

	splash := newBasicSplashScene("Splash", nil)
	err = splash.Build(world)
	if err != nil {
		log.Fatal(err)
	}

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := custom.NewBasicBootScene("Boot", splash)

	nodes.PrintTree(splash)

	engine.PushStart(boot)

	engine.Begin()

}
