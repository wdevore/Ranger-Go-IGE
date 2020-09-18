package audio

import (
	"math"
	"math/rand"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

const (
	noiseBufferSize   = 32
	flangerBufferSize = 1024
)

type generator struct {
	elapsedSinceRepeat int

	period                float64
	periodMax             float64
	enableFrequencyCutoff bool
	periodMult            float64
	periodMultSlide       float64

	dutyCycle      float64
	dutyCycleSlide float64

	arpeggioMultiplier float64
	arpeggioTime       int

	waveShape     int
	prevWaveShape int
	sawtoothRise  bool

	// Filter
	fltw                float64
	enableLowPassFilter bool
	fltwD               float64
	fltdmp              float64
	flthp               float64
	flthpD              float64

	// Vibrato
	vibratoSpeed     float64
	vibratoAmplitude float64

	// Envelope
	envelopeLength []int
	envelopePunch  float64

	// Flanger
	flangerOffset      float64
	flangerOffsetSlide float64

	// Repeat
	repeatTime int

	bitsPerChannel int

	// Data: samples
	samples     []float64
	noiseBuffer []float64

	streamPosition int
	canBeDrained   bool

	values api.IGeneratorValues
}

// NewSfxrGenerator a new samples generator
func NewSfxrGenerator() api.ISampleGenerator {
	o := new(generator)

	return o
}

func (g *generator) Init(v api.IGeneratorValues) {
	g.values = v
	g.initForRepeat(v)
	g.setForRepeat(v)
}

func (g *generator) ReInit() {
	g.initForRepeat(g.values)
	g.setForRepeat(g.values)
}

func (g *generator) SetSawtoothRise(rise bool) {
	g.sawtoothRise = rise
}

func (g *generator) Samples() *[]float64 {
	return &g.samples
}

func (g *generator) initForRepeat(v api.IGeneratorValues) {
	g.elapsedSinceRepeat = 0

	g.period = 100 / (v.BaseFreq()*v.BaseFreq() + 0.001)
	g.periodMax = 100 / (v.FreqLimit()*v.FreqLimit() + 0.001)
	g.enableFrequencyCutoff = (v.FreqLimit() > 0)
	g.periodMult = 1 - math.Pow(v.FreqRamp(), 3)*0.01
	g.periodMultSlide = -math.Pow(v.FreqDramp(), 3) * 0.000001

	g.dutyCycle = 0.5 - v.Duty()*0.5
	g.dutyCycleSlide = -v.DutyRamp() * 0.00005

	if v.ArpMod() >= 0 {
		g.arpeggioMultiplier = 1 - math.Pow(v.ArpMod(), 2)*.9
	} else {
		g.arpeggioMultiplier = 1 + math.Pow(v.ArpMod(), 2)*10
	}
	g.arpeggioTime = int(math.Floor(math.Pow(1-v.ArpSpeed(), 2)*20000 + 32))
	if v.ArpSpeed() == 1 {
		g.arpeggioTime = 0
	}

	// Repeat
	g.repeatTime = int(math.Floor(math.Pow(1.0-v.RepeatSpeed(), 2.0)*20000.0)) + 32
	if v.RepeatSpeed() == 0.0 {
		g.repeatTime = 0
	}

	g.enableLowPassFilter = (v.LpfFreq() != 1.0)
}

func (g *generator) setForRepeat(v api.IGeneratorValues) {
	g.elapsedSinceRepeat = 0

	g.envelopeLength = []int{
		int(math.Floor(v.Attack() * v.Attack() * 100000.0)),
		int(math.Floor(v.Sustain() * v.Sustain() * 100000.0)),
		int(math.Floor(v.Decay() * v.Decay() * 100000.0)),
	}
	g.envelopePunch = v.Punch()

	g.prevWaveShape = g.waveShape
	g.waveShape = v.WaveShape()

	g.period = 100.0 / (v.BaseFreq()*v.BaseFreq() + 0.001)
	g.periodMax = 100.0 / (v.FreqLimit()*v.FreqLimit() + 0.001)
	g.enableFrequencyCutoff = (v.FreqLimit() > 0.0)
	g.periodMult = 1 - math.Pow(v.FreqRamp(), 3.0)*0.01
	g.periodMultSlide = -math.Pow(v.FreqDramp(), 3.0) * 0.000001

	g.dutyCycle = 0.5 - v.Duty()*0.5
	g.dutyCycleSlide = -v.DutyRamp() * 0.00005

	if v.ArpMod() >= 0.0 {
		g.arpeggioMultiplier = 1.0 - math.Pow(v.ArpMod(), 2.0)*0.9
	} else {
		g.arpeggioMultiplier = 1.0 + math.Pow(v.ArpMod(), 2.0)*10.0
	}

	g.arpeggioTime = int(math.Floor(math.Pow(1.0-v.ArpSpeed(), 2.0)*20000.0 + 32.0))
	if v.ArpSpeed() == 1.0 {
		g.arpeggioTime = 0
	}

	// Vibrato
	g.vibratoSpeed = math.Pow(v.VibSpeed(), 2.0) * 0.01
	g.vibratoAmplitude = v.VibStrength() * 0.5

	// Repeat
	g.repeatTime = int(math.Floor(math.Pow(1.0-v.RepeatSpeed(), 2.0)*20000.0)) + 32
	if v.RepeatSpeed() == 0.0 {
		g.repeatTime = 0
	}

	g.flangerOffset = math.Pow(v.PhaOffset(), 2.0) * 1020.0
	if v.PhaOffset() < 0.0 {
		g.flangerOffset = -g.flangerOffset
	}
	g.flangerOffsetSlide = math.Pow(v.PhaRamp(), 2.0) * 1.0
	if v.PhaRamp() < 0.0 {
		g.flangerOffsetSlide = -g.flangerOffsetSlide
	}

	// Filter
	g.fltw = math.Pow(v.LpfFreq(), 3.0) * 0.1
	g.enableLowPassFilter = (v.LpfFreq() != 1.0)
	g.fltwD = 1.0 + v.LpfRamp()*0.0001
	g.fltdmp = 5.0 / (1.0 + math.Pow(v.LpfResonance(), 2.0)*20.0) * (0.01 + g.fltw)
	if g.fltdmp > 0.8 {
		g.fltdmp = 0.8
	}

	g.flthp = math.Pow(v.HpfFreq(), 2.0) * 0.1
	g.flthpD = 1.0 + v.HpfRamp()*0.0003
}

// Generates samples in the range: [-1.0, 1.0].
func (g *generator) Generate(values api.IGeneratorValues) {
	g.bitsPerChannel = values.SampleSize()
	g.noiseBuffer = values.Noise()

	g.initForRepeat(values)
	// g.setForRepeat(values)

	g.samples = []float64{}

	if g.waveShape != g.prevWaveShape || g.noiseBuffer == nil {
		if g.waveShape == api.WaveNoise || g.waveShape == api.WaveNoisePink || g.waveShape == api.WaveNoiseBrownian {
			g.generateNoise(g.waveShape)
		} else {
			if g.noiseBuffer != nil {
				g.noiseBuffer = nil
			}
		}
	}

	envelopeStage := 0
	envelopeElapsed := 0

	vibratoPhase := 0.0

	phase := 0
	flangerIndex := 0

	flangerBuffer := make([]float64, flangerBufferSize)

	sampleSum := 0.0
	numSummed := 0
	summands := int(math.Floor(44100.0 / float64(values.SampleRate())))

	for t := 0; ; t++ {
		g.elapsedSinceRepeat++
		if g.repeatTime != 0 && g.elapsedSinceRepeat >= g.repeatTime {
			g.initForRepeat(values)
		}

		// -----------------------------
		// Arpeggio (single)
		// -----------------------------
		if g.arpeggioTime != 0 && t >= g.arpeggioTime {
			g.arpeggioTime = 0
			g.period *= g.arpeggioMultiplier
		}

		// -----------------------------
		// Frequency slide, and frequency slide slide!
		// -----------------------------
		g.periodMult += g.periodMultSlide
		g.period *= g.periodMult
		if g.period > g.periodMax {
			g.period = g.periodMax
			if g.enableFrequencyCutoff {
				break
			}
		}

		// -----------------------------
		// Vibrato
		// -----------------------------
		rfperiod := g.period
		if g.vibratoAmplitude > 0.0 {
			vibratoPhase += g.vibratoSpeed
			rfperiod = g.period * (1.0 + math.Sin(vibratoPhase)*g.vibratoAmplitude)
		}

		iPeriod := int(math.Floor(rfperiod))
		if iPeriod < api.StandardOverSampling {
			iPeriod = api.StandardOverSampling
		}

		// -----------------------------
		// Square/Sawtooth/Triangle wave duty cycle
		// -----------------------------
		g.dutyCycle += g.dutyCycleSlide
		g.dutyCycle = maths.Clamp(g.dutyCycle, 0.0, 0.5)

		// -----------------------------
		// Volume envelope
		// -----------------------------
		envelopeElapsed++
		if envelopeElapsed > g.envelopeLength[envelopeStage] {
			envelopeElapsed = 0
			envelopeStage++
			if envelopeStage > api.EnvelopeDecay {
				break // Hit Release stage
			}
		}

		envelopeVolume := 0.0
		if g.envelopeLength[envelopeStage] != 0 {
			envelopeVolume = float64(envelopeElapsed) / float64(g.envelopeLength[envelopeStage]) // Envelope Attack
		}

		switch envelopeStage {
		case api.EnvelopeAttack:
		case api.EnvelopeSustain:
			envelopeVolume = 1.0 + (1.0-envelopeVolume)*2.0*g.envelopePunch
		case api.EnvelopeDecay:
			envelopeVolume = 1.0 - envelopeVolume
		}

		// -----------------------------
		// Flanger step
		// -----------------------------
		g.flangerOffset += g.flangerOffsetSlide
		iPhase := int(math.Abs(math.Floor(g.flangerOffset)))
		if iPhase > flangerBufferSize-1 {
			iPhase = flangerBufferSize - 1
		}

		if g.flthpD != 0.0 {
			g.flthp *= g.flthpD
			g.flthp = maths.Clamp(g.flthp, 0.00001, 0.1)
		}

		// -----------------------------
		// Sampling
		// -----------------------------
		// The final sample after oversampling.
		sample := 0.0

		// Use Oversampling to calculate the final sample
		for si := 0; si < api.StandardOverSampling; si++ {
			phase++

			if phase >= iPeriod {
				phase %= iPeriod
				if g.waveShape == api.WaveNoise {
					for i := 0; i < noiseBufferSize; i++ {
						g.noiseBuffer[i] = rand.Float64()*2.0 - 1.0
					}
				}
			}

			subSample := 0.0

			// Base waveform
			fp := float64(phase) / float64(iPeriod)

			switch g.waveShape {
			case api.WaveSQUARE:
				if fp < g.dutyCycle {
					subSample = 0.5
				} else {
					subSample = -0.5
				}
			case api.WaveTriangle:
				if fp < g.dutyCycle {
					subSample = -1.0 + 2.0*fp/g.dutyCycle
				} else {
					subSample = 1.0 - 2.0*(fp-g.dutyCycle)/(1.0-g.dutyCycle)
				}
			case api.WaveSawtooth:
				if g.sawtoothRise {
					subSample = -1.0 + 1.0*fp/g.dutyCycle // Rising (default)
				} else {
					subSample = 1.0 - 2.0*(fp-g.dutyCycle)/(1.0-g.dutyCycle) // Falling
				}
			case api.WaveSINE:
				subSample = math.Sin(fp * 2.0 * math.Pi)
			case api.WaveNoise, api.WaveNoisePink, api.WaveNoiseBrownian:
				subSample = g.noiseBuffer[(phase * (noiseBufferSize / iPeriod))]
			default:
				panic("ERROR: Unknown wave type")
			}

			filterPass := 0.0
			fltdp := 0.0
			filterPassHigh := 0.0

			// -----------------------------
			// Low-pass filter
			// -----------------------------
			pp := filterPass
			g.fltw *= g.fltwD
			g.fltw = maths.Clamp(g.fltw, 0.0, 0.1)
			if g.enableLowPassFilter {
				fltdp += (subSample - filterPass) * g.fltw
				fltdp -= fltdp * g.fltdmp
			} else {
				filterPass = subSample
				fltdp = 0.0
			}
			filterPass += fltdp

			// -----------------------------
			// High-pass filter
			// -----------------------------
			filterPassHigh += filterPass - pp
			filterPassHigh -= filterPassHigh * g.flthp
			subSample = filterPassHigh

			// Flanger
			flangerBuffer[flangerIndex&(flangerBufferSize-1)] = subSample
			subSample += flangerBuffer[(flangerIndex-iPhase+flangerBufferSize)&(flangerBufferSize-1)]
			flangerIndex = (flangerIndex + 1) & (flangerBufferSize - 1)

			// final accumulation and envelope application
			sample += (subSample * envelopeVolume)
		}

		// Accumulate sub-samples appropriately for sample rate
		sampleSum += sample
		numSummed++
		if numSummed >= summands {
			numSummed = 0
			sample = sampleSum / float64(summands)
			sampleSum = 0.0
		} else {
			continue
		}

		// Reference O'Reilly's WebAudio book for Volume verses Gain.
		sample = sample / float64(api.StandardOverSampling) // * MASTER_VOLUME;
		sample *= values.SoundVol()

		// if sample < -1 {
		// 	sample = -1
		// }
		// if sample > 1 {
		// 	sample = 1
		// }
		// fmt.Println(sample)
		g.samples = append(g.samples, sample)
	}
}

func (g *generator) generateNoise(waveShape int) {
	g.noiseBuffer = make([]float64, noiseBufferSize)

	switch waveShape {
	case api.WaveNoise:
		// Noise between [-1.0, 1.0]
		for i := 0; i < noiseBufferSize; i++ {
			g.noiseBuffer[i] = rand.Float64()*2.0 - 1.0
		}
	case api.WaveNoisePink:
		b0 := 0.0
		b1 := 0.0
		b2 := 0.0
		b3 := 0.0
		b4 := 0.0
		b5 := 0.0
		b6 := 0.0
		for i := 0; i < noiseBufferSize; i++ {
			white := rand.Float64()*2.0 - 1.0
			b0 = 0.99886*b0 + white*0.0555179
			b1 = 0.99332*b1 + white*0.0750759
			b2 = 0.96900*b2 + white*0.1538520
			b3 = 0.86650*b3 + white*0.3104856
			b4 = 0.55000*b4 + white*0.5329522
			b5 = -0.7616*b5 - white*0.0168980
			g.noiseBuffer[i] = b0 + b1 + b2 + b3 + b4 + b5 + b6 + white*0.5362
			g.noiseBuffer[i] *= 0.11 // (roughly) compensate for gain
			g.noiseBuffer[i] = maths.Clamp(g.noiseBuffer[i], -1.0, 1.0)
			b6 = white * 0.115926
		}
	case api.WaveNoiseBrownian:
		lastOut := 0.0
		for i := 0; i < noiseBufferSize; i++ {
			white := rand.Float64()*2.0 - 1.0
			g.noiseBuffer[i] = rand.Float64()*2.0 - 1.0
			g.noiseBuffer[i] = (lastOut + (0.02 * white)) / 1.02
			lastOut = g.noiseBuffer[i]
			g.noiseBuffer[i] *= 3.5 // (roughly) compensate for gain
			g.noiseBuffer[i] = maths.Clamp(g.noiseBuffer[i], -1.0, 1.0)
		}
	}
}

func (g *generator) CanBeDrained(drained bool) {
	g.canBeDrained = drained
}

// --------------------------------------------------------------
// Streamer interface
// --------------------------------------------------------------
func (g *generator) Stream(samples [][2]float64) (n int, ok bool) {
	samL := len(g.samples)
	if g.streamPosition >= samL {
		g.streamPosition = 0
		return 0, !g.canBeDrained
	}

	for i := range samples {
		if g.streamPosition >= samL {
			break
		}
		sample := g.samples[g.streamPosition]
		samples[i][0] = sample
		samples[i][1] = sample
		g.streamPosition++
		n++
	}

	return n, true
}

func (g *generator) Err() error {
	return nil
}
