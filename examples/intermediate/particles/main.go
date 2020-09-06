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
	engine.Push(splash)

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := extras.NewBasicBootScene("Boot")

	nodes.PrintTree(splash)

	engine.Push(boot)

	engine.Begin()

}
