package api

// ------------------------------------------------------
// Paths
// ------------------------------------------------------
const (
	RelativeShaderPath = "/engine/assets/shaders/"
)

const (
	// MonoAtlasName is the Map name for StaticMono Atlas
	MonoAtlasName = "MonoAtlas"
	// DynamicMonoAtlasName is the Map name for a DynamicMono Atlas
	DynamicMonoAtlasName = "DynamicMonoAtlas"
	// DynamicPixelAtlasName is the Map name for DynamicPixel Atlas
	DynamicPixelAtlasName = "DynamicPixelAtlas"
)

// These shape names are provided for convience. Use them for
// learning, but typically you would create your own shapes
// specific to your game which means you would have your own
// shape names.
const (
	CenteredFilledSquareShapeName     = "CenteredFilledSquareShape"
	CenteredOutlinedSquareShapeName   = "CenteredOutlinedSquareShape"
	UnCenteredFilledSquareShapeName   = "UnCenteredFilledSquareShape"
	UnCenteredOutlinedSquareShapeName = "UnCenteredOutlinedSquareShape"

	FilledTriangleShapeName   = "FilledTriangleShape"
	OutlinedTriangleShapeName = "OutlinedTriangleShape"

	FilledCircleShapeName   = "FilledCircleShape"
	OutlinedCircleShapeName = "OutlinedCircleShape"

	FilledArcShapeName   = "FilledArcShape"
	OutlinedArcShapeName = "OutlinedArcShape"

	FilledZBarShapeName   = "FilledZBarShape"
	OutlinedZBarShapeName = "OutlinedZBarShape"

	HLineShapeName   = "HLineShape"
	VLineShapeName   = "VLineShape"
	PlusShapeName    = "PlusShape"
	PolygonShapeName = "PolygonShape"

	PointsShapeName = "PointShape"

	// Dynamic shapes

	LineShapeName   = "LineShape"
	SquareShapeName = "SquareShape"
)

const (
	// FILLED polygon
	FILLED = 0
	// OUTLINED polygon
	OUTLINED = 1
	// OPENOUTLINED line strip
	OPENOUTLINED = 2
	// FILLOUTLINED both fill and outlined
	FILLOUTLINED = 3

	// CLOSED indicates a polygon should be rendered closed
	CLOSED = 0
	// OPEN indicates a polygon should be rendered open
	OPEN = 1
)

// XYZComponentCount indicates how many parts to a vertex
const XYZComponentCount int = 3

// XYZWComponentCount is a composite of 2D vertex and 2D texture coords
const XYZWComponentCount int = 4

// OpenGL Object types
const (
	GLLines = 0
)

const (
	// MeshStatic represents static VBO buffers
	MeshStatic = 0
	// MeshDynamic represents dynamic single mesh buffers,
	// for example, PixelBuffer
	MeshDynamic = 1
	// MeshDynamicMulti represent dynamic multi mesh buffers,
	// for example, lines
	MeshDynamicMulti = 2
)

const (
	// ----------------------------------------------
	// Physics
	// ----------------------------------------------

	// PTM is Pixels-to-Meters which isn't used in Ranger. It is
	// here as an example from pixel based engines. I wouldn't
	// use it, but instead use STM below.
	// Box2D uses the MKS(meters/kilograms/seconds) unit system.
	PTM = 1.0 / 30.0 // 1 MKS = 30 GUs

	// RangerScale is a value you change according to your desires.
	// The default is 30.0. For example
	RangerScale = 30.0

	// STM is the Scale-to-MKS ratio.
	// Because Ranger uses transforms we don't think in terms of
	// pixels but rather in terms of spaces. Ranger's View-space
	// --without any scaling--is equal to physic-space (aka Box2D-space)
	// Thus if we want, for example, everything is ranger scaled up
	// then we need to scale it back down to physic-space and that
	// is what STM is for.
	STM = 1.0 / RangerScale // 1 MKS = 30 GUs

	// VelocityIterations is a resolution adjustment
	VelocityIterations = 8

	// PositionIterations is a resolution adjustment
	PositionIterations = 3
)

// TextSetter is a functor for clients to what to notify objects of new text
type TextSetter func(string)
