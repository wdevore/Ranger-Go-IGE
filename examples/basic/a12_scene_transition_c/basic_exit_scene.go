package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type sceneExit struct {
	nodes.Node
	nodes.Scene

	pretendWorkCnt  float64
	pretendWorkSpan float64

	delay          api.IDelay
	tweenOntoStage *gween.Tween
}

func newBasicExitScene(name string, world api.IWorld) (api.INode, error) {
	o := new(sceneExit)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 1000.0

	o.delay = nodes.NewDelay()

	return o, nil
}

func (s *sceneExit) build(world api.IWorld) error {
	s.Node.Build(world)

	// This is an example of a custom background node.
	bg, err := newBackgroundNode("Background", world, s)
	if err != nil {
		return err
	}
	dvr := world.Properties().Window.DeviceRes
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	bn := bg.(*backgroundNode)
	bn.setColor(color.NewPaletteInt64(color.DarkGray))

	if err := s.addText(world); err != nil {
		return err
	}

	newBasicExitLayer("Exit Layer", world, s)

	return nil
}

func (s *sceneExit) addText(world api.IWorld) error {
	var err error

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
		return err
	}

	// #3 INode
	textureNode, err := shapes.NewBitmapFont9x9Node(name, atlas, world, s)
	if err != nil {
		return err
	}
	textureNode.SetScale(50)
	textureNode.SetPosition(-300.0, 0.0)
	bf := textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(color.NewPaletteInt64(color.Lime).Array())
	bf.SetText("Exit Scene. Goodbye...")

	return nil
}

func (s *sceneExit) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneTransitioningIn:
		value, isFinished := s.tweenOntoStage.Update(float32(msPerUpdate))

		// Update animation properties
		if isFinished {
			s.setState("Update: ", api.SceneOnStage)
		}
		s.SetPosition(value, s.Position().Y())
	case api.SceneOnStage:
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
			s.delay.SetPauseTime(1000.0)
			s.delay.Reset()
		}

		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningOut:
		// Update animation
		if s.delay.ReadyToTransition() {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.delay.UpdateTransition(msPerUpdate)
	}
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneExit) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

func (s *sceneExit) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Create an animation that drags the scene onto the stage
		// in the +X direction (enters from right)
		vrs := s.World().Properties().Window.DeviceRes
		s.SetPosition(-float32(vrs.Width), 0.0)
		s.tweenOntoStage = gween.New(s.Position().X(), 0.0, s.TransitionDuration(), ease.OutCubic)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneExit) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneExit EnterNode")
	s.SetVisible(true)
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage
func (s *sceneExit) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneExit ExitNode")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)
	return false
}
