package engine

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
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

	// -----------------------------------------
	// Debug
	// -----------------------------------------
	stepEnabled bool
	postNode    api.INode
}

// Construct creates a new Engine
func Construct(relativePath string, overrides string) (eng api.IEngine, err error) {
	o := new(engine)

	o.world = newWorld(relativePath)

	if !o.world.Properties().Engine.Enabled {
		return nil, errors.New("Engine is NOT enabled in config file")
	}

	// Apply overrides
	if overrides != "" {
		o.world.PropertiesOverride(overrides)
	}

	o.windowDisplay = display.NewDisplay(o)
	err = o.windowDisplay.Initialize(o.world)

	if err != nil {
		return nil, errors.New("Engine.Construct Display error: " + err.Error())
	}

	wpc := o.world.Properties().Window.ClearColor
	o.windowDisplay.SetClearColor(wpc.R, wpc.G, wpc.B, wpc.A)

	err = o.world.Configure()

	if err != nil {
		return nil, errors.New("Engine.Construct World Configure error: " + err.Error())
	}

	// The NodeManager can't configure itself until OpenGL has
	// been initialized which happens above during display
	// construction.
	o.world.NodeManager().Configure(o.world)

	// ---------------------------------------------------------
	if o.world.Properties().Engine.ShowTimingInfo {
		o.postNode, err = extras.NewDynamicTextNode("TimingInfo", 500, o.world, nil)
		if err != nil {
			return nil, err
		}

		gt2 := o.postNode.(api.IDynamicText)
		gt2.SetText("")
		gt2.SetPixelSize(2.0)
		gic := o.postNode.(api.IColor)
		gic.SetColor(color.NewPaletteInt64(color.Peach).Array())

		o.world.SetPostNode(o.postNode)
	}

	return o, nil
}

func (e *engine) Begin() {
	e.running = true

	e.loop()
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
	fpsCnt := 0
	previousT := time.Now()
	secondCnt := int64(0)
	renderElapsedTime := int64(0)
	renderCnt := int64(0)
	// avgRender := 0.0

	sceneGraph := e.world.NodeManager()

	if !sceneGraph.Begin() {
		panic("Not enough scenes to start engine. There must be 2 or more.")
	}

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
					sceneGraph.Update(msPerUpdate, float64(elapsedNano)*frameScaler)
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
		display.Pre() // Things like background clearing
		// **** Any rendering and timing must occur AFTER this point ****

		renderT := time.Now()

		sceneGraph.PreVisit()

		// Calc interpolation for nodes that need it.
		interpolation := float64(lag) / float64(nsPerUpdate)

		// Once the last scene has exited the stage we stop running.
		moreScenes := sceneGraph.Visit(interpolation)

		if !moreScenes {
			fmt.Println("Engine: no more nodes to visit. Exiting...")
			e.running = false
			continue
		}

		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		// Finish rendering
		// ~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--~--
		sceneGraph.PostVisit()

		if renderCnt >= renderMaxCnt {
			e.world.SetAvgRender(float64(renderElapsedTime) / float64(renderMaxCnt) / 1000.0)
			renderCnt = 0
			renderElapsedTime = 0
		} else {
			renderElapsedTime += (time.Now().Sub(renderT)).Microseconds()
			renderCnt++
		}

		secondCnt += elapsedNano
		if secondCnt >= second {
			if engProps.ShowTimingInfo {
				e.drawStats(e.world.Fps(), e.world.Ups(), e.world.AvgRender())
			}

			e.world.SetFps(fpsCnt)
			e.world.SetUps(upsCnt)
			upsCnt = 0
			fpsCnt = 0
			secondCnt = 0
		}

		// time.Sleep(time.Millisecond * 1)

		fpsCnt++
		previousT = currentT

		display.Swap()
	}
}

func (e *engine) End() {
	fmt.Println("Engine shutting down...")
	// Oh noooo! The world is coming to an end!
	e.world.End()

	e.windowDisplay.Shutdown()
}

func (e *engine) World() api.IWorld {
	return e.world
}

func (e *engine) drawStats(fps, ups int, avgRend float64) {
	// fmt.Printf("fps (%2d), ups (%2d), rend (%2.4f)\n", fps, ups, avgRend)
	if e.postNode != nil {
		w := e.world
		if w.Properties().Engine.ShowTimingInfo {
			gt2 := e.postNode.(api.IDynamicText)
			s := fmt.Sprintf("f:%d u:%d r:%2.3f", w.Fps(), w.Ups(), w.AvgRender())
			gt2.SetText(s)
		}
	}
}
