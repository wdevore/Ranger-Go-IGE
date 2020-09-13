package audio

import "github.com/wdevore/Ranger-Go-IGE/api"

// SfxrJSON is configuration properties for the Generator
type SfxrJSON struct {
	Format                          string
	Category                        string
	Name                            string
	BaseFrequency                   float64
	FrequencyLimit                  float64
	FrequencyRamp                   float64
	FrequencyDeltaRamp              float64
	VibratoStrength                 float64
	VibratoSpeed                    float64
	VibratoDelay                    float64
	ArpeggioMod                     float64
	ArpeggioSpeed                   float64
	DutyCycle                       float64
	DutyCycleRamp                   float64
	RepeatSpeed                     float64
	FlangerPhaseOffset              float64
	FlangerPhaseRamp                float64
	LowPassFilterFrequency          float64
	LowPassFilterFrequencyRamp      float64
	LowPassFilterFrequencyResonance float64
	HighPassFilterFrequency         float64
	HighPassFilterFrequencyRamp     float64
	SoundVolume                     float64
	WaveShape                       int
	EnvelopeAttack                  float64
	EnvelopeSustain                 float64
	EnvelopePunch                   float64
	EnvelopeDecay                   float64
	SampleRate                      int
	SampleSize                      int
	Noise                           []float64
}

// CopyFrom transfer data from generator values
func (s *SfxrJSON) CopyFrom(v api.IGeneratorValues) {
	s.WaveShape = v.WaveShape()

	// Envelope
	s.EnvelopeAttack = v.Attack()
	s.EnvelopeSustain = v.Sustain()
	s.EnvelopePunch = v.Punch()
	s.EnvelopeDecay = v.Decay()

	// Tone
	s.BaseFrequency = v.BaseFreq()
	s.FrequencyLimit = v.FreqLimit()
	s.FrequencyRamp = v.FreqRamp()
	s.FrequencyDeltaRamp = v.FreqDramp()

	// Vibrato
	s.VibratoStrength = v.VibStrength()
	s.VibratoSpeed = v.VibSpeed()
	s.VibratoDelay = v.VibDelay()

	// Tonal change
	s.ArpeggioMod = v.ArpMod()
	s.ArpeggioSpeed = v.ArpSpeed()

	// Square wave duty (proportion of time signal is high vs. low)
	s.DutyCycle = v.Duty()
	s.DutyCycleRamp = v.DutyRamp()

	// Repeat
	s.RepeatSpeed = v.RepeatSpeed()

	// Flanger
	s.FlangerPhaseOffset = v.PhaOffset()
	s.FlangerPhaseRamp = v.PhaRamp()

	// Low-pass filter
	s.LowPassFilterFrequency = v.LpfFreq()
	s.LowPassFilterFrequencyRamp = v.LpfRamp()
	s.LowPassFilterFrequencyResonance = v.LpfResonance()

	// High-pass filter
	s.HighPassFilterFrequency = v.HpfFreq()
	s.HighPassFilterFrequencyRamp = v.HpfRamp()

	// Sample parameters
	s.SoundVolume = v.SoundVol()

	s.SampleRate = v.SampleRate()
	s.SampleSize = v.SampleSize()

	s.Noise = v.Noise()
}
