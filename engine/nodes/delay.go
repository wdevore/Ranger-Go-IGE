package nodes

import "github.com/wdevore/Ranger-Go-IGE/api"

type delay struct {
	pauseTime     float64
	pauseCnt      float64
	canTransition bool // true = ready to transition
}

// NewDelay creates a new transition for scene transitions.
func NewDelay() api.IDelay {
	o := new(delay)
	return o
}

// Reset resets the internal timing properties and readies for another cycle.
func (t *delay) Reset() {
	t.pauseCnt = 0.0
	t.canTransition = false
}

// SetPauseTime sets how long the pause lasts, in milliseconds
func (t *delay) SetPauseTime(milliseconds float64) {
	t.pauseTime = milliseconds
}

// Inc used rarely to manually increment the internal counter
func (t *delay) Inc(dt float64) {
	t.pauseCnt += dt
}

// UpdateTransition the internal timer
func (t *delay) UpdateTransition(dt float64) {
	if t.pauseCnt >= t.pauseTime {
		t.canTransition = true
	}
	t.pauseCnt += dt
}

// ReadyToTransition indicates if the node can transition to another scene
func (t *delay) ReadyToTransition() bool {
	return t.canTransition
}
