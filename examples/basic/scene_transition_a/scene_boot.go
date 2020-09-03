package main

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

// Note: this is a very basic boot Node used pretty much for just
// engine development. You should actually supply your own boot node,
// and example can be found in the examples folder.
type sceneBoot struct {
	nodes.Node
	nodes.Scene

	fontTextureRenderer api.ITextureRenderer

	pretendWorkCnt  float64
	pretendWorkSpan float64

	currentState, previousState int
	transitionInCnt             float64
	transitionInDelay           float64
	transitionOutCnt            float64
	transitionOutDelay          float64

	scanCnt   float64
	scanDelay float64

	textureNode api.INode

	dotScale float32
	dots     []api.INode
	colors   []api.IPalette
}

// NewBasicBootScene returns an IScene node of base type INode
func NewBasicBootScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer, replacement api.INode) api.INode {
	o := new(sceneBoot)
	o.Initialize(name)
	o.SetReplacement(replacement)

	o.currentState = api.SceneOffStage
	o.previousState = o.currentState

	o.pretendWorkSpan = 1000.0
	o.scanDelay = 75
	o.dotScale = 15.0

	o.transitionInDelay = 1000.0
	o.transitionOutDelay = 1000.0

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
	switch s.currentState {
	case api.SceneOffStage:
		return
	case api.SceneOnStage:
		if s.pretendWorkCnt > s.pretendWorkSpan {
			// Tell NM that we want to transition off the stage.
			s.setState("Update: ", api.SceneTransitionStartOut)
		}

		s.pretendWorkCnt += msPerUpdate
	case api.SceneTransitioningIn:
		if s.transitionInCnt > s.transitionInDelay {
			tn := s.textureNode.(*extras.BitmapFont9x9Node)
			tn.SetText("OnStage")
			s.setState("Update: ", api.SceneOnStage)
		}
		s.transitionInCnt += msPerUpdate
		// Update animation properties
	case api.SceneTransitioningOut:
		// Update animation
		if s.transitionOutCnt > s.transitionOutDelay {
			s.setState("Update: ", api.SceneExitedStage)
		}
		s.transitionOutCnt += msPerUpdate
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

// Transition indicates what to transition to next
func (s *sceneBoot) State() (current, previous int) {
	return s.currentState, s.previousState
}

func (s *sceneBoot) Notify(state int) {
	s.setState("Notify: ", state)

	switch s.currentState {
	case api.SceneTransitionStartIn:
		// Configure animation properties for entering the stage.
		s.setState("Notify T: ", api.SceneTransitioningIn)
	}
}

func (s *sceneBoot) setState(header string, state int) {
	s.previousState = s.currentState
	s.currentState = state
	nodes.ShowState(header, s, "")
}

// -----------------------------------------------------
// Node lifecycles
// -----------------------------------------------------

// EnterNode called when a node is entering the stage
func (s *sceneBoot) EnterNode(man api.INodeManager) {
	fmt.Println("sceneboot EnterNode")
	man.RegisterTarget(s)
}

// EnterStageNode ---
func (s *sceneBoot) EnterStageNode(man api.INodeManager) {
	// fmt.Println("sceneboot EnterStageNode")
	// Setup transition animation and start it.
	// However, for this example the boot scene appears immediately.
	// So we respond with transition complete.
	// s.setState(api.SceneTransitioningIn)
}

// ExitNode called when a node is exiting stage
func (s *sceneBoot) ExitNode(man api.INodeManager) {
	fmt.Println("sceneboot exit")
	man.UnRegisterTarget(s)
	s.setState("ExitNode: ", api.SceneOffStage)
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
