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

	// b := false
	defaultColor := imgui.Vec4{X: 0.25, Y: 0.25, Z: 0.25, W: 1.0}

	imgui.Begin("Wave Form")

	if sound.SfxrJ.WaveShape == api.WaveSQUARE {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Square") {
		sound.SfxrJ.WaveShape = api.WaveSQUARE
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveTriangle {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Triangle") {
		sound.SfxrJ.WaveShape = api.WaveTriangle
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveSINE {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sine") {
		sound.SfxrJ.WaveShape = api.WaveSINE
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveSawtooth {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Sawtooth") {
		sound.SfxrJ.WaveShape = api.WaveSawtooth
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveNoise {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("White Noise") {
		sound.SfxrJ.WaveShape = api.WaveNoise
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveNoisePink {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Pink Noise") {
		sound.SfxrJ.WaveShape = api.WaveNoisePink
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.SameLine()
	if sound.SfxrJ.WaveShape == api.WaveNoiseBrownian {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorButton, defaultColor)
	}
	if imgui.Button("Red Noise") {
		sound.SfxrJ.WaveShape = api.WaveNoiseBrownian
		sound.GValues.SetWaveShape(sound.SfxrJ.WaveShape)
		if config.Autoplay {
			fmt.Println("Auto playing")
		}
	}
	imgui.PopStyleColor()

	imgui.End()
}
