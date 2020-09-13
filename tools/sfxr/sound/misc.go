package sound

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
)

var (
	// BFormat is a default beep format
	BFormat = beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 2}
	// GValues is generator values
	GValues api.IGeneratorValues
	// SfxrJ holds the settings for an effect
	SfxrJ audio.SfxrJSON
)

// Generate runs the effects generator
func Generate(values api.IGeneratorValues, generator api.ISampleGenerator) {
	BFormat.SampleRate = beep.SampleRate(SfxrJ.SampleRate)

	generator.Generate(values)
	generator.CanBeDrained(true)
}

// UpdateSfxrData updates SfxrJ values
func UpdateSfxrData(values api.IGeneratorValues) {
	SfxrJ.EnvelopeAttack = values.Attack()
	SfxrJ.EnvelopeSustain = values.Sustain()
	SfxrJ.EnvelopePunch = values.Punch()
	SfxrJ.EnvelopeDecay = values.Decay()

	SfxrJ.BaseFrequency = values.BaseFreq()
	SfxrJ.FrequencyLimit = values.FreqLimit()
	SfxrJ.FrequencyRamp = values.FreqRamp()
	SfxrJ.FrequencyDeltaRamp = values.FreqDramp()

	SfxrJ.VibratoStrength = values.VibStrength()
	SfxrJ.VibratoSpeed = values.VibSpeed()
	SfxrJ.VibratoDelay = values.VibDelay()

	SfxrJ.ArpeggioMod = values.ArpMod()
	SfxrJ.ArpeggioSpeed = values.ArpSpeed()

	SfxrJ.DutyCycle = values.Duty()
	SfxrJ.DutyCycleRamp = values.DutyRamp()

	SfxrJ.RepeatSpeed = values.RepeatSpeed()

	SfxrJ.FlangerPhaseOffset = values.PhaOffset()
	SfxrJ.FlangerPhaseRamp = values.PhaRamp()

	SfxrJ.LowPassFilterFrequency = values.LpfFreq()
	SfxrJ.LowPassFilterFrequencyRamp = values.LpfRamp()
	SfxrJ.LowPassFilterFrequencyResonance = values.LpfResonance()

	SfxrJ.HighPassFilterFrequency = values.HpfFreq()
	SfxrJ.HighPassFilterFrequencyRamp = values.HpfRamp()
}

// Play plays a sound using the speaker
func Play(generator api.ISampleGenerator) {
	buffer := beep.NewBuffer(BFormat)
	buffer.Append(generator)

	// fmt.Println("buf len: ", buffer.Len())
	sound := buffer.Streamer(0, buffer.Len())

	speaker.Play(sound)

	// done := make(chan bool)
	// speaker.Play(beep.Seq(sounds, beep.Callback(func() {
	// 	done <- true
	// })))
	// fmt.Println("Playing...")
	// <-done
	// fmt.Println("Done.")
}
