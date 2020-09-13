package audio

import (
	"math"
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

type generatorValues struct {
	baseValues

	// Tone
	baseFreq  float64 // Start frequency
	freqLimit float64 // Min frequency cutoff
	freqRamp  float64 // Slide (SIGNED)
	freqDramp float64 // Delta slide (SIGNED)

	// Vibrato
	vibStrength float64 // Vibrato depth
	vibSpeed    float64 // Vibrato speed
	vibDelay    float64

	// Tonal change
	arpMod   float64 // Change amount (SIGNED)
	arpSpeed float64 // Change speed

	// Square wave duty (proportion of time signal is high vs. low)
	duty     float64 // Square duty
	dutyRamp float64 // Duty sweep (SIGNED)

	// Repeat
	repeatSpeed float64 // Repeat speed

	// Flanger
	phaOffset float64 // Flanger offset (SIGNED)
	phaRamp   float64 // Flanger sweep (SIGNED)

	// Low-pass filter
	lpfFreq      float64 // Low-pass filter cutoff
	lpfRamp      float64 // Low-pass filter cutoff sweep (SIGNED)
	lpfResonance float64 // Low-pass filter resonance
	// High-pass filter
	hpfFreq float64 // High-pass filter cutoff
	hpfRamp float64 // High-pass filter cutoff sweep (SIGNED)

	// Sample parameters
	soundVol float64

	noise []float64
}

// NewIntervalValues create and load a new set of internal values
func NewIntervalValues(sxfrJ *SfxrJSON) api.IGeneratorValues {
	o := new(generatorValues)

	o.baseValues.waveShape = sxfrJ.WaveShape

	// Envelope
	o.baseValues.attack = sxfrJ.EnvelopeAttack
	o.baseValues.sustain = sxfrJ.EnvelopeSustain
	o.baseValues.punch = sxfrJ.EnvelopePunch
	o.baseValues.decay = sxfrJ.EnvelopeDecay

	// Tone
	o.baseFreq = sxfrJ.BaseFrequency       // Start frequency
	o.freqLimit = sxfrJ.FrequencyLimit     // Min frequency cutoff
	o.freqRamp = sxfrJ.FrequencyRamp       // Slide (SIGNED)
	o.freqDramp = sxfrJ.FrequencyDeltaRamp // Delta slide (SIGNED)

	// Vibrato
	o.vibStrength = sxfrJ.VibratoStrength // Vibrato depth
	o.vibSpeed = sxfrJ.VibratoSpeed       // Vibrato speed
	o.vibDelay = sxfrJ.VibratoDelay

	// Tonal change
	o.arpMod = sxfrJ.ArpeggioMod     // Change amount (SIGNED)
	o.arpSpeed = sxfrJ.ArpeggioSpeed // Change speed

	// Square wave duty (proportion of time signal is high vs. low)
	o.duty = sxfrJ.DutyCycle         // Square duty
	o.dutyRamp = sxfrJ.DutyCycleRamp // Duty sweep (SIGNED)

	// Repeat
	o.repeatSpeed = sxfrJ.RepeatSpeed // Repeat speed

	// Flanger
	o.phaOffset = sxfrJ.FlangerPhaseOffset // Flanger offset (SIGNED)
	o.phaRamp = sxfrJ.FlangerPhaseRamp     // Flanger sweep (SIGNED)

	// Low-pass filter
	o.lpfFreq = sxfrJ.LowPassFilterFrequency               // Low-pass filter cutoff
	o.lpfRamp = sxfrJ.LowPassFilterFrequencyRamp           // Low-pass filter cutoff sweep (SIGNED)
	o.lpfResonance = sxfrJ.LowPassFilterFrequencyResonance // Low-pass filter resonance

	// High-pass filter
	o.hpfFreq = sxfrJ.HighPassFilterFrequency     // High-pass filter cutoff
	o.hpfRamp = sxfrJ.HighPassFilterFrequencyRamp // High-pass filter cutoff sweep (SIGNED)

	// Sample parameters
	o.soundVol = sxfrJ.SoundVolume

	o.sampleRate = sxfrJ.SampleRate
	o.sampleSize = sxfrJ.SampleSize

	o.noise = sxfrJ.Noise
	return o
}

func (i *generatorValues) setToDefaults() {
	i.baseValues.waveShape = api.WaveSQUARE

	// Envelope
	i.baseValues.attack = 0.0
	i.baseValues.sustain = 0.3
	i.baseValues.punch = 0.0
	i.baseValues.decay = 0.4

	// Tone
	i.baseFreq = 0.3  // Start frequency
	i.freqLimit = 0.0 // Min frequency cutoff
	i.freqRamp = 0.0  // Slide (SIGNED)
	i.freqDramp = 0.0 // Delta slide (SIGNED)

	// Vibrato
	i.vibStrength = 0.0 // Vibrato depth
	i.vibSpeed = 0.0    // Vibrato speed
	i.vibDelay = 0.0

	// Tonal change
	i.arpMod = 0.0   // Change amount (SIGNED)
	i.arpSpeed = 0.0 // Change speed

	// Square wave duty (proportion of time signal is high vs. low)
	i.duty = 0.0     // Square duty
	i.dutyRamp = 0.0 // Duty sweep (SIGNED)

	// Repeat
	i.repeatSpeed = 0.0 // Repeat speed

	// Flanger
	i.phaOffset = 0.0 // Flanger offset (SIGNED)
	i.phaRamp = 0.0   // Flanger sweep (SIGNED)

	// Low-pass filter
	i.lpfFreq = 1.0      // Low-pass filter cutoff
	i.lpfRamp = 0.0      // Low-pass filter cutoff sweep (SIGNED)
	i.lpfResonance = 0.0 // Low-pass filter resonance
	// High-pass filter
	i.hpfFreq = 0.0 // High-pass filter cutoff
	i.hpfRamp = 0.0 // High-pass filter cutoff sweep (SIGNED)

	// Sample parameters
	i.soundVol = 1.0

	i.sampleRate = 44100
	i.sampleSize = 8
}

func (i *generatorValues) setForRepeat(sg *generator) {
	sg.setForRepeat(i)

	sg.settings.WaveShape = i.waveShape
	sg.elapsedSinceRepeat = 0.0

	sg.period = 100.0 / (i.baseFreq*i.baseFreq + 0.001)
	sg.periodMax = 100.0 / (i.freqLimit*i.freqLimit + 0.001)
	sg.enableFrequencyCutoff = (i.freqLimit > 0.0)
	sg.periodMult = 1 - math.Pow(i.freqRamp, 3.0)*0.01
	sg.periodMultSlide = -math.Pow(i.freqDramp, 3.0) * 0.000001

	sg.dutyCycle = 0.5 - i.duty*0.5
	sg.dutyCycleSlide = -i.dutyRamp * 0.00005

	if i.arpMod >= 0.0 {
		sg.arpeggioMultiplier = 1.0 - math.Pow(i.arpMod, 2.0)*0.9
	} else {
		sg.arpeggioMultiplier = 1.0 + math.Pow(i.arpMod, 2.0)*10.0
	}

	sg.arpeggioTime = int(math.Floor(math.Pow(1.0-i.arpSpeed, 2.0)*20000.0 + 32.0))
	if i.arpSpeed == 1.0 {
		sg.arpeggioTime = 0
	}

	// Vibrato
	sg.vibratoSpeed = math.Pow(i.vibSpeed, 2.0) * 0.01
	sg.vibratoAmplitude = i.vibStrength * 0.5

	// Repeat
	sg.repeatTime = int(math.Floor(math.Pow(1.0-i.repeatSpeed, 2.0)*20000.0)) + 32
	if i.repeatSpeed == 0.0 {
		sg.repeatTime = 0
	}

	sg.flangerOffset = math.Pow(i.phaOffset, 2.0) * 1020.0
	if i.phaOffset < 0.0 {
		sg.flangerOffset = -sg.flangerOffset
	}
	sg.flangerOffsetSlide = math.Pow(i.phaRamp, 2.0) * 1.0
	if i.phaRamp < 0.0 {
		sg.flangerOffsetSlide = -sg.flangerOffsetSlide
	}

	// Filter
	sg.fltw = math.Pow(i.lpfFreq, 3.0) * 0.1
	sg.enableLowPassFilter = (i.lpfFreq != 1.0)
	sg.fltwD = 1.0 + i.lpfRamp*0.0001
	sg.fltdmp = 5.0 / (1.0 + math.Pow(i.lpfResonance, 2.0)*20.0) * (0.01 + sg.fltw)
	if sg.fltdmp > 0.8 {
		sg.fltdmp = 0.8
	}

	sg.flthp = math.Pow(i.hpfFreq, 2.0) * 0.1
	sg.flthpD = 1.0 + i.hpfRamp*0.0003
}

func (i *generatorValues) Mutate() {
	if rand.Float64() > 0.5 {
		i.baseFreq += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.freqRamp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.freqDramp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.duty += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.dutyRamp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.vibStrength += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.vibSpeed += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.vibDelay += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.baseValues.attack += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.baseValues.sustain += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.baseValues.decay += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.baseValues.punch += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.lpfResonance += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.lpfFreq += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.lpfRamp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.hpfFreq += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.hpfRamp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.phaOffset += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.phaRamp += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.repeatSpeed += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.arpSpeed += frnd(0.1) - 0.05
	}
	if rand.Float64() > 0.5 {
		i.arpMod += frnd(0.1) - 0.05
	}
}

// ConfigurePickupCoin create internal values for sound type.
func ConfigurePickupCoin(waveShape int, withArp bool) api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = waveShape
	o.baseFreq = 0.4 + rand.Float64()*0.5
	o.baseValues.attack = 0.0
	o.baseValues.sustain = rand.Float64() * 0.1
	o.baseValues.decay = 0.1 + rand.Float64()*0.4
	o.baseValues.punch = 0.3 + rand.Float64()*0.3
	if withArp || rand.Float64() > 0.5 {
		o.arpSpeed = 0.5 + rand.Float64()*0.2
		o.arpMod = 0.2 + rand.Float64()*0.4
	}

	return o
}

// ConfigureLaserShoot create internal values for sound type.
func ConfigureLaserShoot() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = madRand(3.0)

	if o.baseValues.waveShape == api.WaveSINE && rand.Float64() > 0.5 {
		o.baseValues.waveShape = madRand(1.0)
	}

	if madRand(2.0) == 0 {
		o.baseFreq = 0.3 + rand.Float64()*0.6
		o.freqLimit = rand.Float64() * 0.1
		o.freqRamp = -0.35 - rand.Float64()*0.3
	} else {
		o.baseFreq = 0.5 + rand.Float64()*0.5
		o.freqLimit = o.baseFreq - 0.2 - rand.Float64()*0.6
		if o.freqLimit < 0.2 {
			o.freqLimit = 0.2
		}
		o.freqRamp = -0.15 - rand.Float64()*0.2
	}

	if o.baseValues.waveShape == api.WaveTriangle {
		o.duty = 1.0
	}

	if rand.Float64() > 0.5 {
		o.duty = rand.Float64() * 0.5
		o.dutyRamp = rand.Float64() * 0.2
	} else {
		o.duty = 0.4 + rand.Float64()*0.5
		o.dutyRamp = -rand.Float64() * 0.7
	}

	o.baseValues.attack = 0.0
	o.baseValues.sustain = 0.1 + rand.Float64()*0.2
	o.baseValues.decay = rand.Float64() * 0.4
	if rand.Float64() > 0.5 {
		o.baseValues.punch = rand.Float64() * 0.3
	}
	if madRand(2.0) == 0 {
		o.phaOffset = rand.Float64() * 0.2
		o.phaRamp = -rand.Float64() * 0.2
	}
	o.hpfFreq = rand.Float64() * 0.3

	return o
}

// ConfigureExplosion create internal values for sound type.
func ConfigureExplosion(waveShape int) api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()
	o.baseValues.waveShape = waveShape

	if rand.Float64() > 0.5 {
		o.baseFreq = sqr(0.1 + rand.Float64()*0.4)
		o.freqRamp = -0.1 + rand.Float64()*0.4
	} else {
		o.baseFreq = sqr(0.2 + rand.Float64()*0.7)
		o.freqRamp = -0.2 - rand.Float64()*0.2
	}

	if madRand(4.0) == 0 {
		o.freqRamp = 0.0
	}

	if madRand(2.0) == 0 {
		o.repeatSpeed = 0.3 + rand.Float64()*0.5
	}

	o.baseValues.attack = 0.0
	o.baseValues.sustain = 0.1 + rand.Float64()*0.3
	o.baseValues.decay = rand.Float64() * 0.5
	o.baseValues.punch = 0.2 + rand.Float64()*0.6

	if rand.Float64() > 0.5 {
		o.phaOffset = -0.3 + rand.Float64()*0.9
		o.phaRamp = -rand.Float64() * 0.3
	}
	if rand.Float64() > 0.5 {
		o.vibStrength = rand.Float64() * 0.7
		o.vibSpeed = rand.Float64() * 0.6
	}

	if madRand(2.0) == 0 {
		o.arpSpeed = 0.6 + rand.Float64()*0.3
		o.arpMod = 0.8 - rand.Float64()*1.6
	}

	return o
}

// ConfigurePowerUp create internal values for sound type.
func ConfigurePowerUp() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	if rand.Float64() > 0.5 {
		o.baseValues.waveShape = api.WaveTriangle
		o.duty = 1.0
	} else {
		o.duty = rand.Float64() * 0.6
	}

	o.baseFreq = 0.2 + rand.Float64()*0.3

	if rand.Float64() > 0.5 {
		o.freqRamp = 0.1 + rand.Float64()*0.4
		o.repeatSpeed = 0.4 + rand.Float64()*0.4
	} else {
		o.freqRamp = 0.05 + rand.Float64()*0.2
		if rand.Float64() > 0.5 {
			o.vibStrength = rand.Float64() * 0.7
			o.vibSpeed = rand.Float64() * 0.6
		}
	}

	o.baseValues.attack = 0.0
	o.baseValues.sustain = rand.Float64() * 0.4
	o.baseValues.decay = 0.1 + rand.Float64()*0.4

	return o
}

// ConfigureHitHurt create internal values for sound type.
func ConfigureHitHurt() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = madRand(3.0)

	switch o.baseValues.waveShape {
	case api.WaveSINE:
		o.baseValues.waveShape = api.WaveNoise
	case api.WaveSQUARE:
		o.duty = rand.Float64() * 0.6
	case api.WaveTriangle, api.WaveSawtooth:
		o.duty = 1.0
	}

	o.baseFreq = 0.2 + rand.Float64()*0.6
	o.freqRamp = -0.3 - rand.Float64()*0.4

	o.baseValues.attack = 0.0
	o.baseValues.sustain = rand.Float64() * 0.1
	o.baseValues.decay = 0.1 + rand.Float64()*0.2

	if rand.Float64() > 0.5 {
		o.hpfFreq = rand.Float64() * 0.3
	}

	return o
}

// ConfigureJump create internal values for sound type.
func ConfigureJump() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = api.WaveSQUARE

	o.duty = rand.Float64() * 0.6

	o.baseFreq = 0.3 + rand.Float64()*0.3
	o.freqRamp = 0.1 + rand.Float64()*0.2

	o.baseValues.attack = 0.0
	o.baseValues.sustain = 0.1 + rand.Float64()*0.3
	o.baseValues.decay = 0.1 + rand.Float64()*0.2

	if rand.Float64() > 0.5 {
		o.hpfFreq = rand.Float64() * 0.3
	}
	if rand.Float64() > 0.5 {
		o.lpfFreq = 1 - rand.Float64()*0.6
	}

	return o
}

// ConfigureBlipSelect create internal values for sound type.
func ConfigureBlipSelect() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = madRand(1.0)
	if o.baseValues.waveShape == api.WaveSQUARE {
		o.duty = rand.Float64() * 0.6
	} else {
		o.duty = 1.0
	}
	o.baseFreq = 0.2 + rand.Float64()*0.4

	o.baseValues.attack = 0.0
	o.baseValues.sustain = 0.1 + rand.Float64()*0.1
	o.baseValues.decay = rand.Float64() * 0.2

	o.hpfFreq = 0.1

	return o
}

// ConfigureSynth creates a random synthetic sound
func ConfigureSynth() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.baseValues.waveShape = madRand(1.0)
	if rand.Float64() > 0.5 {
		o.baseFreq = 0.2477
	} else {
		o.baseFreq = 0.1737
	}

	if madRand(4.0) > 3.0 {
		o.baseValues.attack = frnd(0.5)
	} else {
		o.baseValues.attack = 0
	}

	o.baseValues.sustain = frnd(1)
	o.baseValues.punch = frnd(1)
	o.baseValues.decay = frnd(0.9) + 0.1

	aMod := []float64{0, 0, 0, 0, -0.3162, 0.7454, 0.7454}
	o.arpMod = aMod[madRand(6)]
	o.arpSpeed = frnd(0.5) + 0.4

	o.duty = frnd(1)
	if madRand(2) == 2 {
		o.dutyRamp = frnd(1)
	} else {
		o.dutyRamp = 0
	}

	if madRand(1) == 0 {
		o.lpfFreq = 1
	} else {
		o.lpfFreq = frnd(1) * frnd(1)
	}

	o.lpfRamp = rndr(-1, 1)
	o.lpfResonance = frnd(1)

	if madRand(3) == 3 {
		o.hpfFreq = frnd(1)
	} else {
		o.hpfFreq = 0
	}
	if madRand(3) == 3 {
		o.hpfRamp = frnd(1)
	} else {
		o.hpfRamp = 0
	}

	return o
}

// ConfigureRandom create internal values for sound type.
func ConfigureRandom() api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	if rand.Float64() > 0.5 {
		o.baseFreq = cube(frnd(2.0)-1) + 0.5
	} else {
		o.baseFreq = sqr(frnd(1.0))
	}
	o.freqLimit = 0.0
	o.freqRamp = math.Pow(frnd(2.0)-1, 5)
	if o.baseFreq > 0.7 && o.freqRamp > 0.2 {
		o.freqRamp = -o.freqRamp
	}
	if o.baseFreq < 0.2 && o.freqRamp < -0.05 {
		o.freqRamp = -o.freqRamp
	}
	o.freqDramp = math.Pow(frnd(2.0)-1, 3)

	o.duty = frnd(2.0) - 1
	o.dutyRamp = math.Pow(frnd(2.0)-1, 3)

	o.vibStrength = math.Pow(frnd(2.0)-1, 3)
	o.vibSpeed = rndr(-1.0, 1.0)

	o.baseValues.attack = cube(rndr(-1.0, 1.0))
	o.baseValues.sustain = sqr(rndr(-1.0, 1.0))
	o.baseValues.decay = rndr(-1.0, 1.0)
	o.baseValues.punch = math.Pow(frnd(0.8), 2)
	if o.baseValues.attack+o.baseValues.sustain+o.baseValues.decay < 0.2 {
		o.baseValues.sustain += 0.2 + frnd(0.3)
		o.baseValues.decay += 0.2 + frnd(0.3)
	}

	o.lpfResonance = rndr(-1.0, 1.0)
	o.lpfFreq = 1 - math.Pow(frnd(1.0), 3)
	o.lpfRamp = math.Pow(frnd(2.0)-1, 3)
	if o.lpfFreq < 0.1 && o.lpfRamp < -0.05 {
		o.lpfRamp = -o.lpfRamp
	}

	o.hpfFreq = math.Pow(frnd(1.0), 5)
	o.hpfRamp = math.Pow(frnd(2.0)-1, 5)

	o.phaOffset = math.Pow(frnd(2.0)-1, 3)
	o.phaRamp = math.Pow(frnd(2.0)-1, 3)

	o.repeatSpeed = frnd(2.0) - 1

	o.arpSpeed = frnd(2.0) - 1
	o.arpMod = frnd(2.0) - 1

	return o
}

// ConfigureTone create internal values for sound type.
// A4 = 440Hz
func ConfigureTone(tone float64, waveShape int) api.IGeneratorValues {
	o := new(generatorValues)

	o.setToDefaults()

	o.soundVol = api.PlaybackSoundVolume
	o.sampleRate = api.StandardSampleRate
	o.sampleSize = api.StandardSampleSize

	o.waveShape = waveShape

	// sqrt((440Hz / (oversampling = 8) / 441) - 0.001)
	// o.baseFreq = 0.35173363968773563 // 440 Hz
	o.ToIBaseFreq(tone)
	// o.freqRamp = 0.27
	o.attack = 0.0

	// seconds = p^2 * 100000 / 44100
	o.sustain = math.Sqrt(1.0 / 100000 * 44100)
	// o.sustain = 0.664078309 // 1 sec
	// o.sustain = 0.939148551 // 2 sec

	o.decay = 0.0
	o.punch = 0.0

	return o
}

func madRand(max float64) int {
	return int(math.Floor(rand.Float64() * (max + 1.0)))
}

func sqr(v float64) float64 {
	return v * v
}

func rndr(from, to float64) float64 {
	return rand.Float64()*(to-from) + from
}

func cube(v float64) float64 {
	return v * v * v
}
func log(x, b float64) float64 {
	return math.Log(x) / math.Log(b)
}
func flurp(x float64) float64 {
	return x / (1.0 - x)
}
func frnd(rang float64) float64 {
	return rand.Float64() * rang
}

func (i *generatorValues) Attack() float64           { return i.baseValues.attack }
func (i *generatorValues) SetAttach(v float64)       { i.baseValues.attack = v }
func (i *generatorValues) Sustain() float64          { return i.baseValues.sustain }
func (i *generatorValues) SetSustain(v float64)      { i.baseValues.sustain = v }
func (i *generatorValues) Punch() float64            { return i.baseValues.punch }
func (i *generatorValues) SetPunch(v float64)        { i.baseValues.punch = v }
func (i *generatorValues) Decay() float64            { return i.baseValues.decay }
func (i *generatorValues) SetDecay(v float64)        { i.baseValues.decay = v }
func (i *generatorValues) SampleRate() int           { return i.sampleRate }
func (i *generatorValues) SetSampleRate(v int)       { i.sampleRate = v }
func (i *generatorValues) SampleSize() int           { return i.sampleSize }
func (i *generatorValues) SetSampleSize(v int)       { i.sampleSize = v }
func (i *generatorValues) BaseFreq() float64         { return i.baseFreq }
func (i *generatorValues) SetBaseFreq(v float64)     { i.baseFreq = v }
func (i *generatorValues) FreqLimit() float64        { return i.freqLimit }
func (i *generatorValues) SetFreqLimit(v float64)    { i.freqLimit = v }
func (i *generatorValues) FreqRamp() float64         { return i.freqRamp }
func (i *generatorValues) SetFreqRamp(v float64)     { i.freqRamp = v }
func (i *generatorValues) FreqDramp() float64        { return i.freqDramp }
func (i *generatorValues) SetFreqDramp(v float64)    { i.freqDramp = v }
func (i *generatorValues) VibStrength() float64      { return i.vibStrength }
func (i *generatorValues) SetVibStrength(v float64)  { i.vibStrength = v }
func (i *generatorValues) VibSpeed() float64         { return i.vibSpeed }
func (i *generatorValues) SetVibSpeed(v float64)     { i.vibSpeed = v }
func (i *generatorValues) VibDelay() float64         { return i.vibDelay }
func (i *generatorValues) SetVibDelay(v float64)     { i.vibDelay = v }
func (i *generatorValues) ArpMod() float64           { return i.arpMod }
func (i *generatorValues) SetArpMod(v float64)       { i.arpMod = v }
func (i *generatorValues) ArpSpeed() float64         { return i.arpSpeed }
func (i *generatorValues) SetArpSpeed(v float64)     { i.arpSpeed = v }
func (i *generatorValues) Duty() float64             { return i.duty }
func (i *generatorValues) SetDuty(v float64)         { i.duty = v }
func (i *generatorValues) DutyRamp() float64         { return i.dutyRamp }
func (i *generatorValues) SetDutyRamp(v float64)     { i.dutyRamp = v }
func (i *generatorValues) RepeatSpeed() float64      { return i.repeatSpeed }
func (i *generatorValues) SetRepeatSpeed(v float64)  { i.repeatSpeed = v }
func (i *generatorValues) PhaOffset() float64        { return i.phaOffset }
func (i *generatorValues) SetPhaOffset(v float64)    { i.phaOffset = v }
func (i *generatorValues) PhaRamp() float64          { return i.phaRamp }
func (i *generatorValues) SetPhaRamp(v float64)      { i.phaRamp = v }
func (i *generatorValues) LpfFreq() float64          { return i.lpfFreq }
func (i *generatorValues) SetLpfFreq(v float64)      { i.lpfFreq = v }
func (i *generatorValues) LpfRamp() float64          { return i.lpfRamp }
func (i *generatorValues) SetLpfRamp(v float64)      { i.lpfRamp = v }
func (i *generatorValues) LpfResonance() float64     { return i.lpfResonance }
func (i *generatorValues) SetLpfResonance(v float64) { i.lpfResonance = v }
func (i *generatorValues) HpfFreq() float64          { return i.hpfFreq }
func (i *generatorValues) SetHpfFreq(v float64)      { i.hpfFreq = v }
func (i *generatorValues) HpfRamp() float64          { return i.hpfRamp }
func (i *generatorValues) SetHpfRamp(v float64)      { i.hpfRamp = v }
func (i *generatorValues) SoundVol() float64         { return i.soundVol }
func (i *generatorValues) SetSoundVol(v float64)     { i.soundVol = v }
func (i *generatorValues) WaveShape() int            { return i.waveShape }
func (i *generatorValues) SetWaveShape(v int)        { i.waveShape = v }
func (i *generatorValues) Noise() []float64          { return i.noise }
func (i *generatorValues) SetNoise(v []float64)      { i.noise = v }

func (g *generatorValues) ToIBaseFreq(e float64) { g.baseFreq = math.Sqrt((e / 8.0 / 441.0) - 0.001) }
func (g *generatorValues) ToEBaseFreq() float64 {
	return api.StandardOverSampling * 441.0 * (sqr(g.baseFreq) + 0.001)
}

func (g *generatorValues) ToIFreqRamp(e float64) { g.baseFreq = math.Sqrt((e / 8.0 / 441.0) - 0.001) }
func (g *generatorValues) ToEFreqRamp() float64 {
	return 44100.0 * log(1.0-cube(g.freqRamp)/100.0, 0.5)
}
