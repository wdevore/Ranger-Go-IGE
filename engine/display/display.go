package display

import (
	"errors"
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/wdevore/Ranger-Go-IGE/api"
)

// GlfwDisplay glfw window
type GlfwDisplay struct {
	window        *glfw.Window
	quitTriggered bool
}

// New creates a new diplay
func New() *GlfwDisplay {
	o := new(GlfwDisplay)
	return o
}

// Closed checks the window's close status
func (g *GlfwDisplay) Closed() bool {
	return g.window.ShouldClose()
}

// Poll checks for quit or polls events
func (g *GlfwDisplay) Poll() {
	if g.quitTriggered {
		g.window.SetShouldClose(true)
	} else {
		glfw.PollEvents()
	}
}

// Swap is synced to the vertical which means it is waits based on the monitor refresh rate.
// The Clear is also locked to the sync, so if we don't swap the display just waits/locks thus the
// engine appears frozen.
func (g *GlfwDisplay) Swap() {
	g.window.SwapBuffers()
}

// Initialize init GLFW and GL
func (g *GlfwDisplay) Initialize(world api.IWorld) error {
	err := g.initGLFW(world)

	if err != nil {
		return err
	}

	err = g.initGL(world)

	if err != nil {
		return err
	}

	return nil
}

func (g *GlfwDisplay) initGLFW(world api.IWorld) error {
	// Init will call glfw.Terminate if it fails.
	println("Initializing GLFW...")
	err := glfw.Init()
	if err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	wp := world.Properties().Window
	// Create a GLFWwindow object that we can use for GLFW's functions
	g.window, err = glfw.CreateWindow(
		wp.DeviceRes.Width, wp.DeviceRes.Height,
		wp.Title,
		nil, nil)

	if err != nil {
		glfw.Terminate()
		return errors.New("Failed to create GLFW window")
	}

	g.window.SetSizeCallback(g.framebufferSizeCallback)
	// w.window.SetUserPointer()

	g.window.SetPos(wp.Position.X, wp.Position.Y)

	g.window.MakeContextCurrent()

	if world.Properties().Engine.ShowMonitorInfo {
		println("---------------------------- Monitor Info ---------------------------------------")
		monitor := glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()
		fmt.Printf("Monitor refresh rate: %d Hz\n", mode.RefreshRate)
		fmt.Printf("Monitor colors: RGB(%d, %d, %d)\n", mode.RedBits, mode.GreenBits, mode.BlueBits)
		fmt.Printf("Monitor size: %d x %d\n", mode.Width, mode.Height)
		pWidth, pHeight := monitor.GetPhysicalSize()

		fmt.Printf("Physical size: %d x %d\n", pWidth, pHeight)

		fbWidth, fbHeight := g.window.GetFramebufferSize()
		fmt.Printf("Framebuffer size: %d x %d\n", fbWidth, fbHeight)
		println("-------------------------------------------------------------------")
	}

	g.window.SetKeyCallback(g.keyCallback)

	if wp.LockToVSync {
		println("Locking to VSync")
		glfw.SwapInterval(1)
	}

	return nil
}

func (g *GlfwDisplay) initGL(world api.IWorld) error {
	println("Initializing OpenGL...")

	err := gl.Init()

	if err != nil {
		return nil
	}

	ep := world.Properties().Engine

	if ep.ShowGLInfo {
		println("---------------------------- GL Info ---------------------------------------")
		fmt.Printf("Requesting OpenGL minimum of: %d.%d\n", ep.GLMajorVersion, ep.GLMinorVersion)

		version := gl.GoStr(gl.GetString(gl.VERSION))
		fmt.Printf("GL Version obtained: %s\n", version)

		vender := gl.GoStr(gl.GetString(gl.VENDOR))
		fmt.Printf("GL vender: %s\n", vender)

		renderer := gl.GoStr(gl.GetString(gl.RENDERER))
		fmt.Printf("GL renderer: %s\n", renderer)

		glslVersion := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
		fmt.Printf("GLSL version: %s\n", glslVersion)

		var nrAttributes int32
		gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, &nrAttributes)
		fmt.Printf("Max # of vertex attributes supported: %d\n", nrAttributes)
		println("-------------------------------------------------------------------")
	}

	return nil
}

func (g *GlfwDisplay) keyCallback(glfwW *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// println("key pressed")
	if key == glfw.KeyQ && action == glfw.Press {
		g.quitTriggered = true
	}
}

func (g *GlfwDisplay) framebufferSizeCallback(glfwW *glfw.Window, width int, height int) {
	fmt.Printf("Framebuffer re-size: %d x %d", width, height)
}
