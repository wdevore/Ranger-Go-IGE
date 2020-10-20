package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type sceneBoot struct {
	nodes.Node
	nodes.Scene

	pretendWorkCnt  float64
	pretendWorkSpan float64

	delay api.IDelay

	scanCnt   float64
	scanDelay float64

	textureNode api.INode

	dotScale float32
	dots     []api.INode
	colors   []api.IPalette
}

func newBasicBootScene(name string, world api.IWorld) (api.INode, error) {
	o := new(sceneBoot)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 1000.0
	o.scanDelay = 75
	o.dotScale = 15.0

	o.delay = nodes.NewDelay()

	// This is an example of a custom background node.
	bg, err := newBackgroundNode("Background", world, o)
	if err != nil {
		return nil, err
	}
	dvr := world.Properties().Window.DeviceRes
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))
	bn := bg.(*backgroundNode)
	bn.setColor(color.NewPaletteInt64(color.LightGray))

	if err := o.addText(world); err != nil {
		return nil, err
	}

	o.colors = []api.IPalette{
		color.NewPaletteInt64(0xFBD872FF),
		color.NewPaletteInt64(0xFFC845FF),
		color.NewPaletteInt64(0xFFB81CFF),
		color.NewPaletteInt64(0xC69214FF),
		color.NewPaletteInt64(0xAD841FFF),
	}

	if err := o.buildScanThingy(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (s *sceneBoot) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg, err := newBackgroundNode("Background", world, s)
	if err != nil {
		return err
	}
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	return nil
}

func (s *sceneBoot) addText(world api.IWorld) error {
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
	s.textureNode, err = shapes.NewBitmapFont9x9Node(name, atlas, world, s)
	if err != nil {
		return err
	}
	s.textureNode.SetScale(50)
	s.textureNode.SetPosition(-300.0, 0.0)
	// s.textureNode.SetRotation(20.0 * maths.DegreeToRadians)
	bf := s.textureNode.(*shapes.BitmapFont9x9Node)
	bf.SetColor(color.NewPaletteInt64(color.LightOrange).Array())
	bf.SetText("Loading")

	return nil
}

func (s *sceneBoot) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		if s.delay.ReadyToTransition() {
			tn := s.textureNode.(*shapes.BitmapFont9x9Node)
			tn.SetText("OnStage")
			s.setState("Update: ", api.SceneOnStage)
		}
		s.delay.UpdateTransition(msPerUpdate)
		// Update animation properties
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

	s.animate(msPerUpdate)
}

func (s *sceneBoot) animate(msPerUpdate float64) {
	if s.scanCnt > s.scanDelay {
		s.scanCnt = 0.0
		// Shift colors
		c := s.colors[4]
		s.colors[4] = s.colors[3]
		s.colors[3] = s.colors[2]
		s.colors[2] = s.colors[1]
		s.colors[1] = s.colors[0]
		s.colors[0] = c

		for i, dot := range s.dots {
			gol2 := dot.(*shapes.MonoSquareNode)
			gol2.SetFilledColor(s.colors[i])
		}
	}
	s.scanCnt += msPerUpdate
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneBoot) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		// Configure animation properties for entering the stage.
		s.delay.SetPauseTime(1000.0)
		s.setState("Notify T: ", api.SceneTransitioningIn)
	case api.SceneTransitionStartOut:
		s.setState("Notify T: ", api.SceneTransitioningOut)
	}
}

func (s *sceneBoot) setState(header string, state int) {
	s.SetCurrentState(state)
	// nodes.ShowState(header, s, "")
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneBoot) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneboot EnterNode")
	man.RegisterTarget(s)
}

// ExitNode called when a node is exiting stage.
// Return true if this node is to be "repooled" to avoid
// being destroyed.
func (s *sceneBoot) ExitScene(man api.INodeManager) bool {
	// fmt.Println("Sceneboot exit")
	// s.SetVisible(false)
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)

	return false
}

func (s *sceneBoot) buildScanThingy(world api.IWorld) error {
	x := float32(75.0)
	for i := 0; i < 5; i++ {
		dot, err := shapes.NewMonoSquareNode("Square", api.FILLED, true, world, s)
		if err != nil {
			return err
		}

		dot.SetScale(s.dotScale)
		dot.SetPosition(x, -10.0)
		gsq := dot.(*shapes.MonoSquareNode)
		gsq.SetFilledColor(s.colors[i])
		s.dots = append(s.dots, dot)

		x += s.dotScale
	}

	return nil
}
