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

	// ------------------------------------------------------------
	// Load any textures we need
	// ------------------------------------------------------------
	textureMan := world.TextureManager()
	textureMan.AddAtlas("Font9x9", "../../../assets/images/atlas/", "font9x9_texture_manifest.json", true)

	// ------------------------------------------------------------
	// Set a custom background clear effect
	preNode, err := extras.NewStaticCheckerboardNode("CheckBackground", world, nil)
	if err != nil {
		panic(err)
	}
	grn := preNode.(*extras.StaticCheckerboardNode)
	grn.SetSize(200)

	engine.SetPreNode(preNode)
	// ------------------------------------------------------------

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
