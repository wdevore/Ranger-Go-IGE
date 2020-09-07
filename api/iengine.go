package api

// IEngine is the main engine API
type IEngine interface {
	// Start launches the game loop
	Begin()

	// Ends shuts down the engine
	End()

	// World provides access to the engine's world properties
	World() IWorld
}
