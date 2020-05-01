package engine

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

const (
	second       = 1000000000
	renderMaxCnt = int64(50)
)

func init() {
	// GLFW is only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// engine is the main component of ranger
type engine struct {
	world api.IWorld

	// -----------------------------------------
	// Engine properties
	// -----------------------------------------
	running bool

	// ---------------------------------------------------------------------
	// Display
	// ---------------------------------------------------------------------
	windowDisplay *display.GlfwDisplay
	camera        *display.Camera
	viewport      *display.Viewport

	// -----------------------------------------
	// Scene graph is a node manager
	// -----------------------------------------
	sceneGraph api.INodeManager

	// -----------------------------------------
	// Debug
	// -----------------------------------------
	stepEnabled bool
}

// Construct creates a new Engine
func Construct(relativePath string) api.IEngine {
	o := new(engine)
	o.world = newWorld(relativePath)

	o.sceneGraph = nodes.NewNodeManager(o.world)

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

	wpc := e.world.Properties().Window.ClearColor
	e.windowDisplay.SetClearColor(wpc.R, wpc.G, wpc.B, wpc.A)

	e.configureProjections(e.world)

	e.running = true

	e.loop()

	return nil
}

func (e *engine) loop() {
	display := e.windowDisplay
	engProps := e.world.Properties().Engine

	lag := int64(0)
	updatePeriod := float64(second) / engProps.UPSRate
	frameToUpdateRatio := engProps.FPSRate / engProps.UPSRate
	frameScaler := frameToUpdateRatio / 1000000000.0

	nsPerUpdate := int64(math.Round(updatePeriod))
	msPerUpdate := float64(nsPerUpdate) / 1000000.0 // <-- milliseconds, Ex: 33.33333 or 16.6666666
	upsCnt := 0
	ups := 0
	fpsCnt := 0
	fps := 0
	previousT := time.Now()
	secondCnt := int64(0)
	renderElapsedTime := int64(0)
	renderCnt := int64(0)
	avgRender := 0.0

	for !display.Closed() && e.running {
		currentT := time.Now()

		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// Pump IO events
		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		display.Poll()

		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// Update
		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		elapsedNano := (currentT.Sub(previousT)).Nanoseconds()

		// Note: This update is based on:
		// https://gameprogrammingpatterns.com/game-loop.html

		if !e.stepEnabled {
			lag += elapsedNano
			lagging := true
			for lagging {
				if lag >= nsPerUpdate {
					e.sceneGraph.Update(msPerUpdate, float64(elapsedNano)*frameScaler)
					lag -= nsPerUpdate
					upsCnt++
				} else {
					lagging = false
				}
			}
		}

		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// Render Scenegraph by visiting the nodes
		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		renderT := time.Now()

		e.sceneGraph.PreVisit()

		// **** Any rendering must occur AFTER this point ****
		display.Pre()

		// Calc interpolation for nodes that need it.
		interpolation := float64(lag) / float64(nsPerUpdate)

		// Once the last scene has exited the stage we stop running.
		moreScenes := e.sceneGraph.Visit(interpolation)

		if !moreScenes {
			// e.running = false
			// continue
		}

		if renderCnt >= renderMaxCnt {
			avgRender = float64(renderElapsedTime) / float64(renderMaxCnt) / 1000.0
			renderCnt = 0
			renderElapsedTime = 0
		} else {
			renderElapsedTime += (time.Now().Sub(renderT)).Microseconds()
			renderCnt++
		}

		secondCnt += elapsedNano
		if secondCnt >= second {
			if engProps.ShowTimingInfo {
				e.drawStats(fps, ups, avgRender)
			}

			fps = fpsCnt
			ups = upsCnt
			upsCnt = 0
			fpsCnt = 0
			secondCnt = 0
		}

		// time.Sleep(time.Millisecond * 1)

		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// Finish rendering
		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// presentT := time.Now()
		// SDL present's elapsed time is different if vsync is on or off.
		e.sceneGraph.PostVisit()

		fpsCnt++
		previousT = currentT

		display.Swap()
	}
}

func (e *engine) End() {
	fmt.Println("Engine shutting down...")
	e.windowDisplay.Shutdown()
}

func (e *engine) World() api.IWorld {
	return e.world
}

func (e *engine) configureProjections(world api.IWorld) {
	wp := world.Properties().Window

	e.viewport = display.NewViewport()

	e.viewport.SetDimensions(0, 0, wp.DeviceRes.Width, wp.DeviceRes.Height)
	e.viewport.Apply()

	// Calc the aspect ratio between the physical (aka device) dimensions and the
	// the virtual (aka user's design choice) dimensions.

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

func (e *engine) drawStats(fps, ups int, avgRend float64) {
	fmt.Printf("fps (%2d), ups (%2d), rend (%2.4f)\n", fps, ups, avgRend)
	// fmt.Printf("secCnt %d, fpsCnt %d, presC %d\n", secondCnt, fpsCnt, presentElapsedCnt)
}

// ---------------- Update BEGIN -----------------------------
// e.currentUpdateTime = glfw.GetTime()
// e.deltaUpdateTime = glfw.GetTime() - e.currentUpdateTime
// ---------------- Update END -----------------------------
