package engine

import (
	"errors"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
)

func init() {
	// GLFW is only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// engine is the main component of ranger
type engine struct {
	world api.IWorld

	camera *display.Camera

	// ---------------------------------------------------------------------
	// Window
	// ---------------------------------------------------------------------
	windowDisplay *display.GlfwDisplay
}

// Construct creates a new Engine
func Construct(relativePath string) api.IEngine {
	o := new(engine)
	o.world = newWorld(relativePath)
	return o
}

func (e *engine) Start() error {
	if !e.world.Properties().Engine.Enabled {
		return errors.New("engine is NOT enabled in config file")
	}

	e.windowDisplay = display.New()

	err := e.windowDisplay.Initialize(e.world)

	if err != nil {
		return err
	}

	e.configureProjections(e.world)

	// e.renderContext.SetClearColors(graphics.Orange)

	e.loop()

	return nil

}

func (e *engine) loop() {
	for !e.windowDisplay.Closed() {
		e.windowDisplay.Poll()

		// ---------------- Update BEGIN -----------------------------
		// e.currentUpdateTime = glfw.GetTime()
		// e.deltaUpdateTime = glfw.GetTime() - e.currentUpdateTime
		// ---------------- Update END -----------------------------

		// This clear sync locked with the vertical refresh. The clear itself
		// takes ~30 microseconds on a mid-range mobile nvidia GPU.
		// e.renderContext.Clear()
		gl.ClearColor(1.0, 0.5, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		e.windowDisplay.Swap()
	}
}

func (e *engine) End() {

}

func (e *engine) World() api.IWorld {
	return e.world
}

func (e *engine) configureProjections(world api.IWorld) {
	// e.Viewport.SetDimensions(0, 0, config.Window.DeviceRes.Width, config.Window.DeviceRes.Height)
	// e.Viewport.Apply()

	// Calc the aspect ratio between the physical (aka device) dimensions and the
	// the virtual (aka user's design choice) dimensions.

	wp := world.Properties().Window

	deviceRatio := float64(wp.DeviceRes.Width) / float64(wp.DeviceRes.Height)
	virtualRatio := float64(wp.VirtualRes.Width) / float64(wp.VirtualRes.Height)

	xRatioCorrection := float64(wp.DeviceRes.Width) / float64(wp.VirtualRes.Width)
	yRatioCorrection := float64(wp.DeviceRes.Height) / float64(wp.VirtualRes.Height)

	var ratioCorrection float64

	if virtualRatio < deviceRatio {
		ratioCorrection = yRatioCorrection
	} else {
		ratioCorrection = xRatioCorrection
	}

	e.camera = display.NewCamera()

	if world.Properties().Camera.Centered {
		e.camera.SetCenteredProjection()
	} else {
		e.camera.SetProjection(
			float32(ratioCorrection),
			0.0, 0.0,
			float32(wp.DeviceRes.Height), float32(wp.DeviceRes.Width))
	}
}
