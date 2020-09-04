package nodes

import "github.com/wdevore/Ranger-Go-IGE/api"

// Scene is an embedded type for nodes of type IScene.
// Boot scenes and Splash scenes are typical examples.
type Scene struct {
	replacement api.INode

	currentState, previousState int
}

// InitializeScene setups composite
func (s *Scene) InitializeScene(current, previous int) {
	s.currentState = current
	s.previousState = previous
}

// SetReplacement sets this node's replacment during transitions.
func (s *Scene) SetReplacement(replacement api.INode) {
	s.replacement = replacement
}

// GetReplacement returns the replacement node for transitions.
func (s *Scene) GetReplacement() api.INode {
	return s.replacement
}

// State returns both the current and previous state.
func (s *Scene) State() (current, previous int) {
	return s.currentState, s.previousState
}

// Notify is the channel NodeManager uses to cmd the Scene
func (s *Scene) Notify(state int) {
}

// CurrentState returns scene's current transition state.
func (s *Scene) CurrentState() int {
	return s.currentState
}

// SetCurrentState sets the current state
func (s *Scene) SetCurrentState(current int) {
	s.currentState = current
}
