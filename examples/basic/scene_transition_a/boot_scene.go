package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type sceneBoot struct {
	nodes.Node
	nodes.Scene

	fontTextureRenderer api.ITextureRenderer

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

// NewBasicBootScene returns an IScene node of base type INode
func NewBasicBootScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer) api.INode {
	o := new(sceneBoot)
	o.Initialize(name)

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = 1000.0
	o.scanDelay = 75
	o.dotScale = 15.0

	o.delay = nodes.NewDelay()

	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	o.textureNode, err = extras.NewBitmapFont9x9Node("Ranger", textureAtlas, fontRenderer, world, o)
	if err != nil {
		panic(err)
	}
	o.textureNode.SetScale(50)
	o.textureNode.SetPosition(-280.0, 0.0)

	tn := o.textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Loading")
	tn.SetColor(color.NewPaletteInt64(color.LightOrange).Array())

	o.colors = []api.IPalette{
		color.NewPaletteInt64(0xFBD872FF),
		color.NewPaletteInt64(0xFFC845FF),
		color.NewPaletteInt64(0xFFB81CFF),
		color.NewPaletteInt64(0xC69214FF),
		color.NewPaletteInt64(0xAD841FFF),
	}

	o.buildScanThingy(world)
	return o
}

func (s *sceneBoot) Update(msPerUpdate, secPerUpdate float64) {
	switch s.CurrentState() {
	case api.SceneOffStage:
		return
	case api.SceneTransitioningIn:
		if s.delay.ReadyToTransition() {
			tn := s.textureNode.(*extras.BitmapFont9x9Node)
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
			gol2 := dot.(*extras.StaticSquareNode)
			gol2.SetColor(s.colors[i])
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
	// fmt.Println("sceneboot exit")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)

	return false
}

func (s *sceneBoot) buildScanThingy(world api.IWorld) {
	x := float32(75.0)
	for i := 0; i < 5; i++ {
		dot, _ := extras.NewStaticSquareNode("FilledSqr", true, true, world, s)
		s.dots = append(s.dots, dot)
		dot.SetScale(s.dotScale)
		dot.SetPosition(x, -10.0)
		gol2 := dot.(*extras.StaticSquareNode)
		gol2.SetColor(s.colors[i])
		x += s.dotScale
	}
}
