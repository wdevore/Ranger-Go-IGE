package gui

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/sound"
)

// BuildWaveformPanel ...
func BuildWaveformPanel(config *settings.ConfigJSON) {
	imgui.SetNextWindowPos(imgui.Vec2{X: 130, Y: 20.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 515, Y: 60})

	defaultColor := imgui.Vec4{X: 0.25, Y: 0.25, Z: 0.25, W: 1.0}

	imgui.Begin("Wave Form")

	if sound.GValues.WaveShape() == api.WaveSQUARE {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Square") {
		sound.GValues.SetWaveShape(api.WaveSQUARE)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveTriangle {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Triangle") {
		sound.GValues.SetWaveShape(api.WaveTriangle)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveSINE {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sine") {
		sound.GValues.SetWaveShape(api.WaveSINE)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveSawtooth {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sawtooth") {
		sound.GValues.SetWaveShape(api.WaveSawtooth)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveNoise {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("White Noise") {
		sound.GValues.SetWaveShape(api.WaveNoise)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveNoisePink {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Pink Noise") {
		sound.GValues.SetWaveShape(api.WaveNoisePink)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.GValues.WaveShape() == api.WaveNoiseBrownian {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Red Noise") {
		sound.GValues.SetWaveShape(api.WaveNoiseBrownian)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.End()
}
