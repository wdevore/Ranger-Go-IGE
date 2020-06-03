package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
)

func main() {
	engine, err := engine.Construct("../../../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	// Set a custom background clear effect
	preNode, err := custom.NewStaticCheckerboardNode("CheckBackground", world, nil)
	if err != nil {
		panic(err)
	}
	grn := preNode.(*custom.StaticCheckerboardNode)
	grn.SetSize(200)

	engine.SetPreNode(preNode)

	splash := newBasicSplashScene("Splash", nil)
	err = splash.Build(world)
	if err != nil {
		panic(err)
	}

	// This example uses the super basic Boot scene that does absolutely nothing.
	boot := custom.NewBasicBootScene("Boot", splash)

	nodes.PrintTree(splash)

	engine.PushStart(boot)

	engine.Begin()

}
