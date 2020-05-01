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
	VirtualRes   dimensionJSON
	DeviceRes    dimensionJSON
	FullScreen   bool
	Orientation  string
	Position     int2DCoordJSON
	Title        string
}

type cameraJSON struct {
	Centered bool
	View     float3DCoordJSON
}

type fontJSON struct {
	Path         string
	Name         string
	Size         int
	Scale        float64
	CharsFromSet int
}

// Properties reflects config.json
type Properties struct {
	Engine engineJSON
	Window windowJSON
	Camera cameraJSON
	Font   fontJSON
}
