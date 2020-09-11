package main

import (
	"fmt"
	"time"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/gui"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
)

// Platform covers mouse/keyboard/gamepad inputs, cursor shape, timing, windowing.
type Platform interface {
	// ShouldStop is regularly called as the abort condition for the program loop.
	ShouldStop() bool
	// ProcessEvents is called once per render loop to dispatch any pending events.
	ProcessEvents()
	// DisplaySize returns the dimension of the display.
	DisplaySize() [2]float32
	// FramebufferSize returns the dimension of the framebuffer.
	FramebufferSize() [2]float32
	// NewFrame marks the begin of a render pass. It must update the imgui IO state according to user input (mouse, keyboard, ...)
	NewFrame()
	// PostRender marks the completion of one render pass. Typically this causes the display buffer to be swapped.
	PostRender()
	// ClipboardText returns the current text of the clipboard, if available.
	ClipboardText() (string, error)
	// SetClipboardText sets the text as the current text of the clipboard.
	SetClipboardText(text string)
}

type clipboard struct {
	platform Platform
}

func (board clipboard) Text() (string, error) {
	return board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	board.platform.SetClipboardText(text)
}

// Renderer covers rendering imgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [3]float32)
	// Render draws the provided imgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)
}

const (
	millisPerSecond = 1000
	sleepDuration   = time.Millisecond * 33
)

// run implements the main program loop of the app. It returns when the platform signals to stop.
func run(p Platform, r Renderer, config *settings.ConfigJSON) {
	imgui.CurrentIO().SetClipboard(clipboard{platform: p})

	clearColor := [3]float32{0.95, 0.90, 0.85}

	// ch := make(chan string)

	// // Start simulation thread. It will idle by default.
	// go simulator.Run(ch)
	// imgui.StyleColorsClassic()

	// -------------------------------------------------------------
	// Now start main GUI loop
	// -------------------------------------------------------------
	for !p.ShouldStop() {
		p.ProcessEvents()
		p.NewFrame()

		imgui.NewFrame()

		// ---------------------------------------------------------
		// Draw Graph
		// ---------------------------------------------------------
		// ---------------------------------------------------------

		gui.BuildGui(config)

		imgui.EndFrame()
		// ---------------------------------------------------------

		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		r.PreRender(clearColor)

		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())

		p.PostRender()

		// sleep to avoid 100% CPU usage
		<-time.After(sleepDuration)
	}

	fmt.Println("Exiting application")
	// gui.Shutdown(environment)
}
