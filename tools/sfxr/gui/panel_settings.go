package gui

import (
	"github.com/inkyblackness/imgui-go/v2"
)

// BuildGeneratorsPanel ...
func BuildGeneratorsPanel() {
	imgui.SetNextWindowPos(imgui.Vec2{X: 130, Y: 85.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 445, Y: 760})

	imgui.Begin("Settings")

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Envelope") {

		f := float32(0.0)
		imgui.SliderFloat("Attack", &f, 0.0, 1.0)
		imgui.SliderFloat("Sustain Time", &f, 0.0, 1.0)
		imgui.SliderFloat("Sustain Punch", &f, 0.0, 1.0)
		imgui.SliderFloat("Decay Time", &f, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Frequency") {
		g := float32(0.0)
		imgui.SliderFloat("Start frequency", &g, 0.0, 1.0)
		imgui.SliderFloat("Min freq. cutoff", &g, 0.0, 1.0)
		imgui.SliderFloat("Slide", &g, 0.0, 1.0)
		imgui.SliderFloat("Delta Slide", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Vibrato") {
		g := float32(0.0)
		imgui.SliderFloat("Depth", &g, 0.0, 1.0)
		imgui.SliderFloat("Speed", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Arpeggiation") {
		g := float32(0.0)
		imgui.SliderFloat("Frequency mult", &g, 0.0, 1.0)
		imgui.SliderFloat("Change speed", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Duty Cycle") {
		g := float32(0.0)
		imgui.SliderFloat("Duty cycle", &g, 0.0, 1.0)
		imgui.SliderFloat("Sweep##1", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Retrigger") {
		g := float32(0.0)
		imgui.SliderFloat("Rate", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Flanger") {
		g := float32(0.0)
		imgui.SliderFloat("Offset", &g, 0.0, 1.0)
		imgui.SliderFloat("Sweep##2", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Low-Pass Filter") {
		g := float32(0.0)
		imgui.SliderFloat("Cutoff frequency#1", &g, 0.0, 1.0)
		imgui.SliderFloat("Cutoff sweep##1", &g, 0.0, 1.0)
		imgui.SliderFloat("Resonance", &g, 0.0, 1.0)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("High-Pass Filter") {
		g := float32(0.0)
		imgui.SliderFloat("Cutoff frequency#2", &g, 0.0, 1.0)
		imgui.SliderFloat("Cutoff sweep##2", &g, 0.0, 1.0)
	}

	imgui.End()
}
