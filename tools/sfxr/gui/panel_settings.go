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
	title  string
	cTitle string

	controlChanged bool
	min, max       float32
	v              float32
	c              float32
}

type setValue func(v float64)
type getValue func() float64

func newSettingsWidget(title string, min, max float32) *settingsWidget {
	o := new(settingsWidget)
	o.title = title
	o.cTitle = title
	o.min = min
	o.max = max
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
			s.cTitle = fmt.Sprintf("%s (%0.3f)", s.title, getExtValue())
		}
		sound.Generate(sound.GValues, generator)
		sound.Play(generator)
	}
}

var baseFrequency = newSettingsWidget("Frequency", 0.042830679660512114, 1.009657448477109) //15Hz -> 3280Hz
var frequencyRamp = newSettingsWidget("Slide (8va/sec)", -.3, 0.3)
var envelopeSustain = newSettingsWidget("Sustain", 0.0, 3.0)

// BuildSettingsPanel ...
func BuildSettingsPanel(config *settings.ConfigJSON, generator api.ISampleGenerator) {
	imgui.SetNextWindowPos(imgui.Vec2{X: 130, Y: 85.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 500, Y: 770})

	imgui.Begin("Settings")

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Envelope") {

		f := float32(sound.SfxrJ.EnvelopeAttack)
		changed := imgui.SliderFloat("Attack (sec)", &f, 0.0, 3.0)
		if changed {
			sound.SfxrJ.EnvelopeAttack = float64(f)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		envelopeSustain.check(
			func() float64 { return sound.GValues.Sustain() },
			func() float64 { return 2.2 },
			func(v float64) { sound.GValues.SetSustain(v) },
			generator,
		)

		f = float32(sound.SfxrJ.EnvelopePunch)
		changed = imgui.SliderFloat("Sustain Punch (%)", &f, 0.0, 100.0)
		if changed {
			sound.SfxrJ.EnvelopePunch = float64(f)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		f = float32(sound.SfxrJ.EnvelopeDecay)
		changed = imgui.SliderFloat("Decay (sec)", &f, 0.0, 3.0)
		if changed {
			sound.SfxrJ.EnvelopeDecay = float64(f)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Frequency") {
		baseFrequency.check(
			func() float64 { return sound.GValues.BaseFreq() },
			func() float64 { return sound.GValues.ToEBaseFreq() },
			func(v float64) { sound.GValues.SetBaseFreq(v) },
			generator,
		)

		g := float32(sound.SfxrJ.FrequencyLimit)
		changed := imgui.SliderFloat("Min freq-cutoff (Hz)", &g, 3.0, 3600.0)
		if changed {
			sound.SfxrJ.FrequencyLimit = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		frequencyRamp.check(
			func() float64 { return sound.GValues.FreqRamp() },
			func() float64 { return sound.GValues.ToEFreqRamp() },
			func(v float64) { sound.GValues.SetFreqRamp(v) },
			generator,
		)

		g = float32(sound.SfxrJ.FrequencyDeltaRamp)
		changed = imgui.SliderFloat("Delta Slide (8va/sec^2)", &g, 0.09, -0.09)
		if changed {
			sound.SfxrJ.FrequencyDeltaRamp = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Vibrato") {
		g := float32(sound.SfxrJ.VibratoStrength)
		changed := imgui.SliderFloat("Depth (%)", &g, 0.0, 50.0)
		if changed {
			sound.SfxrJ.VibratoStrength = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.VibratoSpeed)
		changed = imgui.SliderFloat("Speed (Hz)", &g, 0.0, 70.0)
		if changed {
			sound.SfxrJ.VibratoSpeed = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.VibratoDelay)
		changed = imgui.SliderFloat("Delay", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.VibratoDelay = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Arpeggiation") {
		g := float32(sound.SfxrJ.ArpeggioMod)
		changed := imgui.SliderFloat("Frequency mult", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.ArpeggioMod = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.ArpeggioSpeed)
		changed = imgui.SliderFloat("Change speed (sec)", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.ArpeggioSpeed = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Duty Cycle") {
		g := float32(sound.SfxrJ.DutyCycle)
		changed := imgui.SliderFloat("Duty cycle", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.DutyCycle = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.DutyCycleRamp)
		changed = imgui.SliderFloat("Sweep##1", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.DutyCycleRamp = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Retrigger") {
		g := float32(sound.SfxrJ.RepeatSpeed)
		changed := imgui.SliderFloat("Rate", &g, 0.0, 1400.0)
		if changed {
			sound.SfxrJ.RepeatSpeed = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Flanger") {
		g := float32(sound.SfxrJ.FlangerPhaseOffset)
		changed := imgui.SliderFloat("Offset (msec)", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.FlangerPhaseOffset = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.FlangerPhaseRamp)
		changed = imgui.SliderFloat("Sweep (msec/sec)##2", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.FlangerPhaseRamp = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("Low-Pass Filter") {
		g := float32(sound.SfxrJ.LowPassFilterFrequency)
		changed := imgui.SliderFloat("Cutoff frequency (Hz)##1", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.LowPassFilterFrequency = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.LowPassFilterFrequencyRamp)
		changed = imgui.SliderFloat("Cutoff sweep (^sec)##1", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.LowPassFilterFrequencyRamp = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.LowPassFilterFrequencyResonance)
		changed = imgui.SliderFloat("Resonance (%)", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.LowPassFilterFrequencyResonance = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.SetNextItemOpen(true, imgui.ConditionAlways)
	if imgui.CollapsingHeader("High-Pass Filter") {
		g := float32(sound.SfxrJ.HighPassFilterFrequency)
		changed := imgui.SliderFloat("Cutoff frequency (Hz)##2", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.HighPassFilterFrequency = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}

		g = float32(sound.SfxrJ.HighPassFilterFrequencyRamp)
		changed = imgui.SliderFloat("Cutoff sweep (^sec)##2", &g, 0.0, 1.0)
		if changed {
			sound.SfxrJ.HighPassFilterFrequencyRamp = float64(g)
		}
		if imgui.IsMouseReleased(0) && changed {
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
	}

	imgui.End()
}
