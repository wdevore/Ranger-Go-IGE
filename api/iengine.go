package api

// IEngine is the main engine API
type IEngine interface {
	// Start launches the game loop
	Begin()

	// Ends shuts down the engine
	End()

	// World provides access to the engine's world properties
	World() IWorld

	SetPreNode(INode)
	SetPostNode(INode)

	// PushStart pushes the given node onto the stack as the
	// first scene to start once the engine's configuration in complete.
	PushStart(INode)

	RouteEvents(event IEvent)
}
