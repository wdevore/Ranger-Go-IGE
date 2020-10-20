package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
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

	// Add Atlas to the world so Scenes/Layers can obtain access to the atlas.
	world.AddAtlas(api.MonoAtlasName, monoAtlas)

	transitionDuration := float32(500.0)

	exitScene, err := newBasicExitScene("Exit", world)
	if err != nil {
		panic(err)
	}
	sceneE := exitScene.(api.IScene)
	sceneE.SetTransitionDuration(transitionDuration)
	exitScene.SetVisible(false)
	world.Push(exitScene)

	splash, err := newBasicSplashScene("Splash", world)
	if err != nil {
		panic(err)
	}
	sceneS := splash.(api.IScene)
	sceneS.SetTransitionDuration(transitionDuration)
	splash.SetVisible(false)
	world.Push(splash)

	boot, err := newBasicBootScene("Boot", world)
	if err != nil {
		panic(err)
	}
	sceneB := boot.(api.IScene)
	sceneB.SetTransitionDuration(transitionDuration)

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
