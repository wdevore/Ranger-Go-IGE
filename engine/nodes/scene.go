package nodes

import "github.com/wdevore/Ranger-Go-IGE/api"

// Scene is an embedded type for nodes of type IScene.
// Boot scenes and Splash scenes are typical examples.
type Scene struct {
	replacement api.INode

	currentState, previousState int

	transitionDuration float32
}

// InitializeScene setups composite
func (s *Scene) InitializeScene(current, previous int) {
	s.currentState = current
	s.previousState = previous
}

// State returns both the current and previous state.
func (s *Scene) State() (current, previous int) {
	return s.currentState, s.previousState
}

// Notify is the channel NodeManager uses to cmd the Scene.
// This is the minimal required for an instant transition
// from Boot to first scene.
func (s *Scene) Notify(state int) {
	s.SetCurrentState(state)

	switch s.CurrentState() {
	case api.SceneTransitionStartIn:
		s.SetCurrentState(api.SceneTransitionStartOut)
	case api.SceneTransitionStartOut:
		s.SetCurrentState(api.SceneOnStage)
	}
}

// CurrentState returns scene's current transition state.
func (s *Scene) CurrentState() int {
	return s.currentState
}

// SetCurrentState sets the current state
func (s *Scene) SetCurrentState(current int) {
	s.currentState = current
}

// EnterScene called when it is time to transition a scene
// onto the stage.
func (s *Scene) EnterScene(man api.INodeManager) {

}

// ExitScene called when a node is exiting stage.
// Return true if this node is to be "repooled" to avoid
// being destroyed.
func (s *Scene) ExitScene(man api.INodeManager) bool {
	return true
}

// TransitionDuration returns how long/fast a transition is
func (s *Scene) TransitionDuration() float32 {
	return s.transitionDuration
}

// SetTransitionDuration sets transition duration.
func (s *Scene) SetTransitionDuration(duration float32) {
	s.transitionDuration = duration
}
