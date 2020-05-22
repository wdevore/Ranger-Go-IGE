package api

const (
	// FILLED polygon
	FILLED = 0
	// OUTLINED polygon
	OUTLINED = 1
	// FILLOUTLINED both fill and outlined
	FILLOUTLINED = 2

	// CLOSED indicates a polygon should be rendered closed
	CLOSED = 0
	// OPEN indicates a polygon should be rendered open
	OPEN = 1
)

// XYZComponentCount indicates how many parts to a vertex
const XYZComponentCount int = 3

// OpenGL Object types
const (
	GLLines = 0
)
