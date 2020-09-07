package main

import (
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

func main() {
	engine, err := engine.Construct("../../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	textureMan := world.TextureManager()

	// ------------------------------------------------------------
	// Load any textures we need
	// ------------------------------------------------------------
	textureMan.AddAtlas("Font9x9", "../../../assets/images/atlas/", "font9x9_texture_manifest.json", true)

	fontTextureRenderer := rendering.NewTextureRenderer(textureMan, world.TextureShader())
	fontTextureRenderer.Build("Font9x9")

	transitionDuration := float32(500.0)

	exitScene, err := newBasicExitScene("Exit", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	sceneE := exitScene.(api.IScene)
	sceneE.SetTransitionDuration(transitionDuration)
	exitScene.SetVisible(false)
	world.Push(exitScene)

	menu, err := newMenuScene("Menu", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	menuS := menu.(api.IScene)
	menuS.SetTransitionDuration(transitionDuration)
	menu.SetVisible(false)
	world.Push(menu)

	settings, err := newSettingsScene("Settings", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	settingsS := settings.(api.IScene)
	settingsS.SetTransitionDuration(transitionDuration)
	settings.SetVisible(false)
	menuSU := menu.(*sceneMenu)
	menuSU.AddSubMenu(settings)

	highscore, err := newHighscoreScene("Highscore", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	highscoreS := highscore.(api.IScene)
	highscoreS.SetTransitionDuration(transitionDuration)
	highscore.SetVisible(false)
	menuSU.AddSubMenu(highscore)

	game, err := newGameScene("Game", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	gameS := game.(api.IScene)
	gameS.SetTransitionDuration(transitionDuration)
	game.SetVisible(false)
	menuSU.AddSubMenu(game)

	// -----------------------------------------------
	splash, err := newBasicSplashScene("Splash", world, fontTextureRenderer)
	if err != nil {
		panic(err)
	}
	sceneS := splash.(api.IScene)
	sceneS.SetTransitionDuration(transitionDuration)
	splash.SetVisible(false)
	world.Push(splash)

	engine.Begin()

}
