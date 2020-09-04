package nodes

import "github.com/wdevore/Ranger-Go-IGE/api"

// Transition holds properties for scene transitions.
type transition struct {
	pauseTime     float64
	pauseCnt      float64
	canTransition bool // true = ready to transition
}

// NewTransition creates a new transition for scene transitions.
func NewTransition() api.ITransition {
	o := new(transition)
	return o
}

// Reset resets the internal timing properties and readies for another cycle.
func (t *transition) Reset() {
	t.pauseCnt = 0.0
	t.canTransition = false
}

// SetPauseTime sets how long the pause lasts, in milliseconds
func (t *transition) SetPauseTime(milliseconds float64) {
	t.pauseTime = milliseconds
}

// Inc used rarely to manually increment the internal counter
func (t *transition) Inc(dt float64) {
	t.pauseCnt += dt
}

// UpdateTransition the internal timer
func (t *transition) UpdateTransition(dt float64) {
	if t.pauseCnt >= t.pauseTime {
		t.canTransition = true
	}
	t.pauseCnt += dt
}

// ReadyToTransition indicates if the node can transition to another scene
func (t *transition) ReadyToTransition() bool {
	return t.canTransition
}
