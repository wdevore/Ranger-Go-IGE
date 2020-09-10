package audio

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
