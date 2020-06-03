package configuration

type engineJSON struct {
	Enabled          bool
	LoopFor          int
	ShowConfig       bool
	ShowGLInfo       bool
	ShowMonitorInfo  bool
	ShowTimingInfo   bool
	ShowJoystickInfo bool
	GLMajorVersion   int
	GLMinorVersion   int
	FPSRate          float64
	UPSRate          float64
}

type colorJSON struct {
	R float32
	G float32
	B float32
	A float32
}

type dimensionJSON struct {
	Height int
	Width  int
}

type int2DCoordJSON struct {
	X int
	Y int
}

type float3DCoordJSON struct {
	X float64
	Y float64
	Z float64
}

type windowJSON struct {
	BitsPerPixel int
	LockToVSync  bool
	ClearColor   colorJSON
	ClearStyle   string // "None", "SingleColor", "Checkerboard", "Custom"
	VirtualRes   dimensionJSON
	DeviceRes    dimensionJSON
	FullScreen   bool
	Orientation  string
	ViewScale    float64
	Position     int2DCoordJSON
	Title        string
}

type depthJSON struct {
	Near float32
	Far  float32
}

type cameraJSON struct {
	Centered bool
	View     float3DCoordJSON
	Depth    depthJSON
}

type fontJSON struct {
	Path         string
	Name         string
	Size         int
	Scale        float64
	CharsFromSet int
}

type shadersJSON struct {
	UseDefault        bool
	VertexShaderSrc   string
	FragmentShaderSrc string

	// These properties are populated dynamically
	VertexShaderCode   string
	FragmentShaderCode string
}

// Properties reflects config.json
type Properties struct {
	Engine  engineJSON
	Window  windowJSON
	Camera  cameraJSON
	Font    fontJSON
	Shaders shadersJSON
}
