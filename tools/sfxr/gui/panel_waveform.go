package gui

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v2"
)

var waveForm int

// BuildWaveformPanel ...
func BuildWaveformPanel() {
	imgui.SetNextWindowPos(imgui.Vec2{X: 130, Y: 20.0})

	// b := false
	defaultColor := imgui.Vec4{X: 0.25, Y: 0.25, Z: 0.25, W: 1.0}

	imgui.Begin("Wave Form")

	if waveForm == 0 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Square") {
		waveForm = 0
		if autoPlay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if waveForm == 1 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sawtooth") {
		waveForm = 1
		if autoPlay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if waveForm == 2 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sine") {
		waveForm = 2
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if waveForm == 3 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("White Noise") {
		waveForm = 3
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if waveForm == 4 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Pink Noise") {
		waveForm = 4
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if waveForm == 5 {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Red Noise") {
		waveForm = 5
	}
	imgui.PopStyleColor()

	imgui.End()
}
