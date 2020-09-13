package api

// IGeneratorValues is an sfxr generator values
type IGeneratorValues interface {
	Mutate()

	Attack() float64
	SetAttack(float64)
	ToEAttack() float64 // Convert from Internal to External
	Sustain() float64
	SetSustain(float64)
	ToESustain() float64
	Punch() float64
	SetPunch(float64)
	Decay() float64
	SetDecay(float64)
	ToEDecay() float64

	SampleRate() int
	SetSampleRate(int)
	SampleSize() int
	SetSampleSize(int)

	BaseFreq() float64
	SetBaseFreq(float64)
	ToIBaseFreq(float64)  // Convert from External to Internal
	ToEBaseFreq() float64 // Convert from Internal to External
	FreqLimit() float64
	SetFreqLimit(float64)
	ToEFreqLimit() float64
	FreqRamp() float64
	SetFreqRamp(float64)
	ToIFreqRamp(float64) // Slide
	ToEFreqRamp() float64
	FreqDramp() float64
	SetFreqDramp(float64)
	ToEFreqDramp() float64

	VibStrength() float64
	SetVibStrength(float64)
	ToEVibStrength() float64
	VibSpeed() float64
	SetVibSpeed(float64)
	ToEVibSpeed() float64
	VibDelay() float64
	SetVibDelay(float64)

	ArpMod() float64
	SetArpMod(float64)
	ToEArpMod() float64
	ArpSpeed() float64
	SetArpSpeed(float64)
	ToEArpSpeed() float64

	Duty() float64
	SetDuty(float64)
	ToEDuty() float64
	DutyRamp() float64
	SetDutyRamp(float64)
	ToEDutyRamp() float64

	RepeatSpeed() float64
	SetRepeatSpeed(float64)
	ToERepeatSpeed() float64

	PhaOffset() float64
	SetPhaOffset(float64)
	ToEPhaOffset() float64
	PhaRamp() float64
	SetPhaRamp(float64)
	ToEPhaRamp() float64

	LpfFreq() float64
	SetLpfFreq(float64)
	ToELpfFreq() float64
	LpfRamp() float64
	SetLpfRamp(float64)
	ToELpfRamp() float64
	LpfResonance() float64
	SetLpfResonance(float64)
	ToELpfResonance() float64

	HpfFreq() float64
	SetHpfFreq(float64)
	ToEHpfFreq() float64
	HpfRamp() float64
	SetHpfRamp(float64)
	ToEHpfRamp() float64

	SoundVol() float64
	SetSoundVol(float64)

	WaveShape() int
	SetWaveShape(int)

	Noise() []float64
	SetNoise([]float64)
}
