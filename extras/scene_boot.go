package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// Note: this is a very basic boot Node used pretty much for just
// engine development. You should actually supply your own boot node,
// and example can be found in the examples folder.
type sceneBoot struct {
	nodes.Node
	nodes.Scene
}

// NewBasicBootScene returns an IScene node of base type INode
func NewBasicBootScene(name string) api.INode {
	o := new(sceneBoot)
	o.Initialize(name)
	return o
}

// --------------------------------------------------------
// Transitioning
// --------------------------------------------------------

func (s *sceneBoot) Notify(state int) {
	// This boot scene does absolutely nothing other than
	// satisfying the NodeManager requirement for 2 scenes.
	// So just simply exit the stage immdediately.
	s.SetCurrentState(api.SceneExitedStage)
}
