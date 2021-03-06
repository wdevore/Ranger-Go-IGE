package api

const (
	// SceneOffStage means a scene off the stage either destroyed,
	// on the stack or in a pool.
	SceneOffStage = iota

	// SceneTransitionStartIn : the scene is beginning to transition
	SceneTransitionStartIn

	// SceneTransitioningIn : the scene is busy transitioning onto the stage
	SceneTransitioningIn

	// SceneOnStage means a scene is actively doing on stage.
	SceneOnStage

	// SceneTransitionStartOut : the scene is beginning to transition
	SceneTransitionStartOut

	// SceneTransitioningOut : the scene is busy transitioning off the stage
	SceneTransitioningOut

	// SceneExitedStage : the scene has finished transitioning off the stage
	SceneExitedStage

	// SceneFinished means a scene is done. Destroy it and/or remove from pool.
	SceneFinished
)

// IScene scene management
type IScene interface {
	Notify(int)

	State() (int, int)
	CurrentState() int
	SetCurrentState(current int)

	TransitionDuration() float32
	SetTransitionDuration(duration float32)

	EnterScene(INodeManager)
	ExitScene(INodeManager) bool
}
