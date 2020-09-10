package api

const (
	// --------------------------------
	// Wave shapes
	// --------------------------------

	// WaveSQUARE a square wave
	WaveSQUARE = 0
	// WaveSAWTOOTH a sawtooth wave
	WaveSAWTOOTH = 1
	// WaveSINE a sine wave
	WaveSINE = 2
	// WaveNoise is white noise
	WaveNoise = 3
	// WaveNoisePink is pink noise
	WaveNoisePink = 4
	// WaveNoiseBrownian is brownian/red noise
	WaveNoiseBrownian = 5

	// --------------------------------
	// Playback volume
	// --------------------------------

	// PlaybackMasterVolume is main volume level
	PlaybackMasterVolume = 1.0
	// PlaybackSoundVolume is ???
	PlaybackSoundVolume = 1.0

	// --------------------------------
	// Sampling
	// --------------------------------

	// StandardSampleRate is a standard rate
	StandardSampleRate = 44100
	// StandardSampleSize is a typical size
	StandardSampleSize = 8
	// StandardOverSampling is basic over sampling
	StandardOverSampling = 8

	// --------------------------------
	// ADSR
	// --------------------------------

	// EnvelopeAttack part of envelope
	EnvelopeAttack = 0
	// EnvelopeSustain part of envelope
	EnvelopeSustain = 1
	// EnvelopeDecay part of envelope
	// EnvelopeDecay part of envelope
	EnvelopeDecay = 2
	// EnvelopeRelease part of envelope
	EnvelopeRelease = 3
)

// ISampleGenerator is a sfxr sample generator
type ISampleGenerator interface {
	Generate(IGeneratorValues)
	CanBeDrained(bool)
	Samples() *[]float64

	Stream(samples [][2]float64) (n int, ok bool)
	Err() error
}
