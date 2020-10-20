package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
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

	// Note: To render text you need 3 objects:
	// SpriteSheet contains the manifest and font image.
	// SingleTextureAtlas renders a single sub-texture (i.e. character).
	// INode will render strings using the atlas.

	name := "Font9x9"
	// #1 SpriteSheet
	spriteSheet := fonts.NewFont9x9SpriteSheet(name, "font9x9_sprite_sheet_manifest.json")
	spriteSheet.Load("../../assets/", true)

	// #2 TextureAtlas
	atlas := atlas.NewSingleTextureAtlas(name, spriteSheet, world)
	err = atlas.Burn()
	if err != nil {
		panic(err)
	}

	transitionDuration := float32(500.0)

	exitScene, err := newBasicExitScene("Exit", atlas, world)
	if err != nil {
		panic(err)
	}
	sceneE := exitScene.(api.IScene)
	sceneE.SetTransitionDuration(transitionDuration)
	exitScene.SetVisible(false)
	world.Push(exitScene)

	menu, err := newMenuScene("Menu", atlas, world)
	if err != nil {
		panic(err)
	}
	menuS := menu.(api.IScene)
	menuS.SetTransitionDuration(transitionDuration)
	menu.SetVisible(false)
	world.Push(menu)

	settings, err := newSettingsScene("Settings", atlas, world)
	if err != nil {
		panic(err)
	}
	settingsS := settings.(api.IScene)
	settingsS.SetTransitionDuration(transitionDuration)
	settings.SetVisible(false)
	menuSU := menu.(*sceneMenu)
	menuSU.AddSubMenu(settings)

	highscore, err := newHighscoreScene("Highscore", atlas, world)
	if err != nil {
		panic(err)
	}
	highscoreS := highscore.(api.IScene)
	highscoreS.SetTransitionDuration(transitionDuration)
	highscore.SetVisible(false)
	menuSU.AddSubMenu(highscore)

	game, err := newGameScene("Game", atlas, world)
	if err != nil {
		panic(err)
	}
	gameS := game.(api.IScene)
	gameS.SetTransitionDuration(transitionDuration)
	game.SetVisible(false)
	menuSU.AddSubMenu(game)

	// -----------------------------------------------
	splash, err := newBasicSplashScene("Splash", atlas, world)
	if err != nil {
		panic(err)
	}
	sceneS := splash.(api.IScene)
	sceneS.SetTransitionDuration(transitionDuration)
	splash.SetVisible(false)
	world.Push(splash)

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
