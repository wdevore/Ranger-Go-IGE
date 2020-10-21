package engine

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
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

	projLoc int32
	viewLoc int32

	defaultBackgroundEnabled bool
	backgroundAtlas          api.IAtlasX

	// -----------------------------------------
	// Debug
	// -----------------------------------------
	stepEnabled bool
	infoNode    api.INode
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

	// -----------------------------------------------------------
	// Display and OpenGL
	// -----------------------------------------------------------
	o.windowDisplay = display.NewDisplay(o)

	// Initializes GLFW and GL.
	err = o.windowDisplay.Initialize(o.world)

	if err != nil {
		return nil, errors.New("Engine.Construct Display error: " + err.Error())
	}

	// _-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-
	// ---- Anything GL wise can be called after this point. -----
	// -_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_-_

	wpc := o.world.Properties().Window.ClearColor
	o.windowDisplay.SetClearColor(wpc.R, wpc.G, wpc.B, wpc.A)

	err = o.world.Configure()

	if err != nil {
		return nil, errors.New("Engine.Construct World Configure error: " + err.Error())
	}

	// The NodeManager can't configure itself until both the World and OpenGL has
	// been initialized which happens above during display construction.
	o.world.NodeManager().Configure(o.world)

	// -----------------------------------------------------------
	// Shaders
	// -----------------------------------------------------------
	err = o.configureUniforms()

	if err != nil {
		return nil, errors.New("Engine.Construct Uniforms Configure error: " + err.Error())
	}

	err = o.world.Begin()
	if err != nil {
		return nil, err
	}

	err = o.configureBackgroundForgrounds()
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (e *engine) configureBackgroundForgrounds() error {
	// -----------------------------------------------------------
	// Timing Info.
	// -----------------------------------------------------------
	worldProps := e.world.Properties()
	var err error

	if worldProps.Engine.ShowTimingInfo {
		e.infoNode, err = shapes.NewDynamicPixelTextNode("FPS", e.world, e.world.Overlay())
		if err != nil {
			return err
		}
		e.infoNode.SetScale(1.5)
		gt := e.infoNode.(*shapes.DynamicPixelPixelTextNode)
		gt.SetText("FPS...")
		gt.SetColor(color.NewPaletteInt64(color.Peach).Array())
		gt.SetPixelSize(2.0)

		dvr := e.world.Properties().Window.DeviceRes
		e.infoNode.SetPosition(-float32(dvr.Width/2)+10.0, -float32(dvr.Height/2)+10.0)
	}

	// -----------------------------------------------------------
	// Should we attach a default background?
	// -----------------------------------------------------------
	switch worldProps.Window.ClearStyle {
	case "SingleColor":
		// The StaticMono Atlas needs to exist BEFORE trying to create
		// static nodes.
		e.backgroundAtlas = e.world.GetAtlas(api.MonoAtlasName)
		if e.backgroundAtlas == nil {
			e.backgroundAtlas = atlas.NewStaticMonoAtlas(e.world)
			e.world.AddAtlas(api.MonoAtlasName, e.backgroundAtlas)
		}

		square, err := shapes.NewMonoSquareNode("Background", api.FILLED, true, e.world, e.world.Underlay())
		if err != nil {
			return err
		}
		dvr := worldProps.Window.DeviceRes
		square.SetScaleComps(float32(dvr.Width), float32(dvr.Height))

		sq := square.(*shapes.MonoSquareNode)
		bgCol := worldProps.Window.BackgroundColor
		sq.SetFilledColor(color.NewPaletteFromFloats(bgCol.R, bgCol.G, bgCol.B, bgCol.A))

		e.defaultBackgroundEnabled = true
	case "Checkerboard":
	}

	return nil
}

// Begin is called after Construct() and as the last thing the game
// does to start the game.
func (e *engine) Begin() error {
	e.running = true

	sceneGraph := e.world.NodeManager()

	var err error

	err = sceneGraph.Begin()
	if err != nil {
		return errors.New("not enough scenes to start engine. There must be 2 or more")
	}

	// If a default background was requested via the config.json then
	// we need to make sure that the associated atlas has been "burnt"
	// prior to starting the loop.
	if e.defaultBackgroundEnabled && e.backgroundAtlas != nil {
		if !e.backgroundAtlas.Burnt() {
			err = e.backgroundAtlas.Burn()
			if err != nil {
				return err
			}
		}
	}

	// nodes.PrintTree(e.world.Root())
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
	fpsCnt := 0
	previousT := time.Now()
	secondCnt := int64(0)
	renderElapsedTime := int64(0)
	renderCnt := int64(0)
	// avgRender := 0.0

	sceneGraph := e.world.NodeManager()

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

func (e *engine) configureUniforms() error {
	// -------------------------------------------------------
	// Default Shader
	programID := e.world.Shader().Program()
	e.world.Shader().Use()

	e.projLoc = gl.GetUniformLocation(programID, gl.Str("projection\x00"))
	if e.projLoc < 0 {
		return errors.New("NodeManager: couldn't find 'projection' uniform variable")
	}

	e.viewLoc = gl.GetUniformLocation(programID, gl.Str("view\x00"))
	if e.viewLoc < 0 {
		return errors.New("NodeManager: couldn't find 'view' uniform variable")
	}

	pm := e.world.Projection().Matrix()
	gl.UniformMatrix4fv(e.projLoc, 1, false, &pm[0])

	vm := e.world.Viewspace().Matrix()
	gl.UniformMatrix4fv(e.viewLoc, 1, false, &vm[0])
	// -------------------------------------------------------

	// -------------------------------------------------------
	// Texture Shader
	programID = e.world.TextureShader().Program()
	e.world.TextureShader().Use()

	projLoc := gl.GetUniformLocation(programID, gl.Str("projection\x00"))
	if projLoc < 0 {
		return errors.New("NodeManager: couldn't find 'projection' uniform variable")
	}

	viewLoc := gl.GetUniformLocation(programID, gl.Str("view\x00"))
	if viewLoc < 0 {
		return errors.New("NodeManager: couldn't find 'view' uniform variable")
	}

	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &vm[0])
	// -------------------------------------------------------

	return nil
}

// TODO : This info should be on the Overlay Node.
func (e *engine) drawStats(fps, ups int, avgRend float64) {
	// fmt.Printf("fps (%2d), ups (%2d), rend (%2.4f)\n", fps, ups, avgRend)
	if e.infoNode != nil {
		w := e.world
		if w.Properties().Engine.ShowTimingInfo {
			// gt2, ok := e.postNode.(api.IDynamicText)
			gt2, ok := e.infoNode.(*shapes.DynamicPixelPixelTextNode)
			if !ok {
				panic("drawStats failed text type assertion")
			}
			s := fmt.Sprintf("f:%d u:%d r:%2.3f", w.Fps(), w.Ups(), w.AvgRender())
			gt2.SetText(s)
		}
	}
}
