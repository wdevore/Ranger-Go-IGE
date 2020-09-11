package gui

import (
	"github.com/inkyblackness/imgui-go/v2"
)

// BuildSettingsPanel ...
func BuildSettingsPanel() {
	imgui.SetNextWindowPos(imgui.Vec2{X: 10, Y: 20.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 115, Y: 420})

	imgui.Begin("Generators")

	if imgui.Button("Pickup/Coin") {
	}

	if imgui.Button("Laser/Shoot") {
	}

	if imgui.Button("Explosion") {
	}

	if imgui.Button("PowerUp") {
	}
	if imgui.Button("Hit/Hurt") {
	}
	if imgui.Button("Blip/Select") {
	}
	if imgui.Button("Synth") {
	}
	if imgui.Button("Random") {
	}
	if imgui.Button("Tone") {
	}

	if imgui.Button("Mutate") {
	}

	imgui.SetNextItemOpen(true, imgui.ConditionOnce)
	if imgui.CollapsingHeader("Control") {
		imgui.IndentV(20)
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.9, Y: 0.5, Z: 0.0, W: 1.0})
		if imgui.Button(" Play ") {
		}
		imgui.PopStyleColor()

		imgui.Unindent()
		imgui.Separator()

		h := float32(1.0)
		imgui.IndentV(10)
		imgui.VSliderFloatV("Volume", imgui.Vec2{X: 30, Y: 100}, &h, 0.0, 10.0, "%.1f", 1.5)
	}

	imgui.End()

}
