package gui

import (
	"fmt"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/sound"
)

var (
	sliderChanged = false
)

type settingsWidget struct {
	title          string
	cTitle         string
	format         string
	controlChanged bool
	min, max       float32
	v              float32
	c              float32
}

type setValue func(v float64)
type getValue func() float64

func newSettingsWidget(title string, min, max float32, format string) *settingsWidget {
	o := new(settingsWidget)
	o.title = title
	o.cTitle = title
	o.min = min
	o.max = max
	o.format = format
	return o
}

func (s *settingsWidget) check(getGValue, getExtValue getValue, setGValue setValue, generator api.ISampleGenerator) {
	s.v = float32(getGValue())

	changed := imgui.SliderFloat(s.cTitle, &s.v, s.min, s.max)

	if changed {
		setGValue(float64(s.v))
		s.c = s.v
		s.controlChanged = true
	}

	if imgui.IsMouseReleased(0) && s.controlChanged {
		s.controlChanged = false
		setGValue(float64(s.c))
		if getExtValue != nil {
			s.cTitle = fmt.Sprintf(s.format, s.title, getExtValue())
		}
		generator.Init(sound.GValues)
		sound.Generate(sound.GValues, generator)
		sound.Play(generator)
	}
}

var envelopeAttack = newSettingsWidget("Attack", 0.0, 1.0, "%s (%0.3f)")
var envelopeSustain = newSettingsWidget("Sustain", 0.0, 1.0, "%s (%0.3f)")
var envelopePunch = newSettingsWidget("Sustain Punch", 0.0, 1.0, "%s (%0.3f)")
var envelopeDecay = newSettingsWidget("Decay", 0.0, 1.0, "%s (%0.3f)")

var baseFrequency = newSettingsWidget("Frequency", 0.042830679660512114, 1.009657448477109, "%s (%0.1f)") //15Hz -> 3280Hz
var frequencyLimit = newSettingsWidget("Min cutoff", 0.0, 1.0, "%s (%0.5f)")
var frequencyRamp = newSettingsWidget("Slide", -1.0, 1.0, "%s (%0.5f)")
var frequencyDeltaRamp = newSettingsWidget("Delta Slide", -1.0, 1.0, "%s (%0.5f)")

var vibratoStrength = newSettingsWidget("Depth", 0.0, 1.0, "%s (%0.3f)")
var vibratoSpeed = newSettingsWidget("Speed", 0.0, 1.0, "%s (%0.3f)")

var arpeggioMod = newSettingsWidget("Multiplier", -1.0, 1.0, "%s (%0.3f)")
var arpeggioSpeed = newSettingsWidget("Speed", 0.0, 1.0, "%s (%0.7f)##1")

var dutyCycle = newSettingsWidget("Duty Cycle", 0.0, 1.0, "%s (%0.4f)")
var dutyCycleRamp = newSettingsWidget("Sweep", -1.0, 1.0, "%s (%0.4f)##1")

var repeatSpeed = newSettingsWidget("Rate", 0.0, 0.96, "%s (%0.4f)")

var flangerPhaseOffset = newSettingsWidget("Offset", -1.0, 1.0, "%s (%0.4f)")
var flangerPhaseRamp = newSettingsWidget("Sweep", -1.0, 1.0, "%s (%0.4f)##2")

var lowPassFilterFrequency = newSettingsWidget("Cutoff Freq", 0.0, 1.0, "%s (%0.4f)##1")
var lowPassFilterFrequencyRamp = newSettingsWidget("Cutoff Sweep", -1.0, 1.0, "%s (%0.4f)##1")
var lowPassFilterFrequencyResonance = newSettingsWidget("Resonance", 0.035, 1.0, "%s (%0.4f)")

var highPassFilterFrequency = newSettingsWidget("Cutoff Freq", 0.0, 1.0, "%s (%0.4f)##2")
var highPassFilterFrequencyRamp = newSettingsWidget("Cutoff Sweep", -1.0, 1.0, "%s (%0.7f)##2")

// BuildSettingsPanel ...
func BuildSettingsPanel(config *settings.ConfigJSON, generator api.ISampleGenerator) {
	imgui.SetNextWindowPos(imgui.Vec2{X: 130, Y: 85.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 510, Y: 770})

	imgui.Begin("Settings")

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Envelope") {
		envelopeAttack.check(
			func() float64 { return sound.GValues.Attack() },
			func() float64 { return sound.GValues.ToEAttack() },
			func(v float64) { sound.GValues.SetAttack(v) },
			generator,
		)

		envelopeSustain.check(
			func() float64 { return sound.GValues.Sustain() },
			func() float64 { return sound.GValues.ToESustain() },
			func(v float64) { sound.GValues.SetSustain(v) },
			generator,
		)

		envelopePunch.check(
			func() float64 { return sound.GValues.Punch() },
			func() float64 { return sound.GValues.Punch() },
			func(v float64) { sound.GValues.SetPunch(v) },
			generator,
		)

		envelopeDecay.check(
			func() float64 { return sound.GValues.Decay() },
			func() float64 { return sound.GValues.ToEDecay() },
			func(v float64) { sound.GValues.SetDecay(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Frequency") {
		baseFrequency.check(
			func() float64 { return sound.GValues.BaseFreq() },
			func() float64 { return sound.GValues.ToEBaseFreq() },
			func(v float64) { sound.GValues.SetBaseFreq(v) },
			generator,
		)

		frequencyLimit.check(
			func() float64 { return sound.GValues.FreqLimit() },
			func() float64 { return sound.GValues.ToEFreqLimit() },
			func(v float64) { sound.GValues.SetFreqLimit(v) },
			generator,
		)

		frequencyRamp.check(
			func() float64 { return sound.GValues.FreqRamp() },
			func() float64 { return sound.GValues.ToEFreqRamp() },
			func(v float64) { sound.GValues.SetFreqRamp(v) },
			generator,
		)

		frequencyDeltaRamp.check(
			func() float64 { return sound.GValues.FreqDramp() },
			func() float64 { return sound.GValues.ToEFreqDramp() },
			func(v float64) { sound.GValues.SetFreqDramp(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Vibrato") {
		vibratoStrength.check(
			func() float64 { return sound.GValues.VibStrength() },
			func() float64 { return sound.GValues.ToEVibStrength() },
			func(v float64) { sound.GValues.SetVibStrength(v) },
			generator,
		)

		vibratoSpeed.check(
			func() float64 { return sound.GValues.VibSpeed() },
			func() float64 { return sound.GValues.ToEVibSpeed() },
			func(v float64) { sound.GValues.SetVibSpeed(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Arpeggiation") {
		arpeggioMod.check(
			func() float64 { return sound.GValues.ArpMod() },
			func() float64 { return sound.GValues.ToEArpMod() },
			func(v float64) { sound.GValues.SetArpMod(v) },
			generator,
		)

		arpeggioSpeed.check(
			func() float64 { return sound.GValues.ArpSpeed() },
			func() float64 { return sound.GValues.ToEArpSpeed() },
			func(v float64) { sound.GValues.SetArpSpeed(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Duty Cycle") {
		switch sound.GValues.WaveShape() {
		case api.WaveSQUARE, api.WaveSawtooth, api.WaveTriangle:
			dutyCycle.check(
				func() float64 { return sound.GValues.Duty() },
				func() float64 { return sound.GValues.ToEDuty() },
				func(v float64) { sound.GValues.SetDuty(v) },
				generator,
			)
			dutyCycleRamp.check(
				func() float64 { return sound.GValues.DutyRamp() },
				func() float64 { return sound.GValues.ToEDutyRamp() },
				func(v float64) { sound.GValues.SetDutyRamp(v) },
				generator,
			)
		default:
			imgui.PushStyleColor(imgui.StyleColorSliderGrab, imgui.Vec4{X: 0.75, Y: 0.75, Z: 0.75, W: 1.0})
			dutyCycle.check(
				func() float64 { return 1.0 },
				nil,
				func(v float64) { sound.GValues.SetDuty(0.0) },
				generator,
			)
			dutyCycleRamp.check(
				func() float64 { return 1.0 },
				nil,
				func(v float64) { sound.GValues.SetDutyRamp(0.0) },
				generator,
			)
			imgui.PopStyleColor()
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Retrigger") {
		repeatSpeed.check(
			func() float64 { return sound.GValues.RepeatSpeed() },
			func() float64 { return sound.GValues.ToERepeatSpeed() },
			func(v float64) { sound.GValues.SetRepeatSpeed(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Flanger") {
		flangerPhaseOffset.check(
			func() float64 { return sound.GValues.PhaOffset() },
			func() float64 { return sound.GValues.ToEPhaOffset() },
			func(v float64) { sound.GValues.SetPhaOffset(v) },
			generator,
		)

		flangerPhaseRamp.check(
			func() float64 { return sound.GValues.PhaRamp() },
			func() float64 { return sound.GValues.ToEPhaRamp() },
			func(v float64) { sound.GValues.SetPhaRamp(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Low-Pass Filter") {
		lowPassFilterFrequency.check(
			func() float64 { return sound.GValues.LpfFreq() },
			func() float64 { return sound.GValues.ToELpfFreq() },
			func(v float64) { sound.GValues.SetLpfFreq(v) },
			generator,
		)

		lowPassFilterFrequencyRamp.check(
			func() float64 { return sound.GValues.LpfRamp() },
			func() float64 { return sound.GValues.ToELpfRamp() },
			func(v float64) { sound.GValues.SetLpfRamp(v) },
			generator,
		)

		lowPassFilterFrequencyResonance.check(
			func() float64 { return sound.GValues.LpfResonance() },
			func() float64 { return sound.GValues.ToELpfResonance() },
			func(v float64) { sound.GValues.SetLpfResonance(v) },
			generator,
		)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("High-Pass Filter") {
		highPassFilterFrequency.check(
			func() float64 { return sound.GValues.HpfFreq() },
			func() float64 { return sound.GValues.ToEHpfFreq() },
			func(v float64) { sound.GValues.SetHpfFreq(v) },
			generator,
		)

		highPassFilterFrequencyRamp.check(
			func() float64 { return sound.GValues.HpfRamp() },
			func() float64 { return sound.GValues.ToEHpfRamp() },
			func(v float64) { sound.GValues.SetHpfRamp(v) },
			generator,
		)
	}

	imgui.End()
}
