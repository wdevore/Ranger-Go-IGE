package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type sceneExit struct {
	nodes.Node
	nodes.Scene

	pretendWorkCnt  float64
	pretendWorkSpan float64

	delay          api.IDelay
	tweenOntoStage *gween.Tween
}

func newBasicExitScene(name string, world api.IWorld, fontRenderer api.ITextureRenderer) (api.INode, error) {
	o := new(sceneExit)
	o.Initialize(name)

	if err := o.build(world); err != nil {
		return nil, err
	}

	o.InitializeScene(api.SceneOffStage, api.SceneOffStage)

	o.pretendWorkSpan = float64(o.TransitionDuration())

	o.delay = nodes.NewDelay()

	textureMan := world.TextureManager()
	var err error

	textureAtlas := textureMan.GetAtlasByName("Font9x9")

	textureNode, err := extras.NewBitmapFont9x9Node("ExitText", textureAtlas, fontRenderer, world, o)
	if err != nil {
		panic(err)
	}
	textureNode.SetScale(25)
	textureNode.SetPosition(-100.0, 0.0)

	tn := textureNode.(*extras.BitmapFont9x9Node)
	tn.SetText("Exit Scene. Goodbye...")
	tn.SetColor(color.NewPaletteInt64(color.White).Array())

	return o, nil
}

func (s *sceneExit) build(world api.IWorld) error {
	s.Node.Build(world)

	dvr := s.World().Properties().Window.DeviceRes

	bg := newBackgroundNode("Background", world, s, color.NewPaletteInt64(color.DarkGray))
	bg.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

	newBasicExitLayer("Exit Layer", world, s)

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
			s.delay.SetPauseTime(float64(s.TransitionDuration()))
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

// EnterScene called when a node is entering the stage
func (s *sceneExit) EnterScene(man api.INodeManager) {
	// fmt.Println("sceneExit EnterScene")
	s.SetVisible(true)
	man.RegisterTarget(s)
}

// ExitScene called when a node is exiting stage
func (s *sceneExit) ExitScene(man api.INodeManager) bool {
	// fmt.Println("sceneExit ExitScene")
	man.UnRegisterTarget(s)
	s.setState("ExitScene: ", api.SceneOffStage)
	return false
}
