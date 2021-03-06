package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

func main() {
	engine, err := engine.Construct("../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	// ------------------------------------------------------------
	// Set a custom background clear effect
	// preNode, err := extras.NewStaticCheckerboardNode("CheckBackground", world, nil)
	// if err != nil {
	// 	panic(err)
	// }
	// grn := preNode.(*extras.StaticCheckerboardNode)
	// grn.SetSize(200)

	// world.Underlay().AddChild(preNode)
	// ------------------------------------------------------------

	splash, err := newBasicSplashScene("Splash", world)
	if err != nil {
		panic(err)
	}
	world.Push(splash)

	// This example uses the super basic Boot scene that does absolutely
	// nothing by exiting immediately.
	boot := extras.NewBasicBootScene("Boot")

	world.Push(boot)

	// And finally we can start the game.
	engine.Begin()
}
