package audio

type baseValues struct {
	// Wave shape/type
	waveShape int

	// Envelope
	attack  float64 // Attack time (secconds)
	sustain float64 // Sustain time (secconds)
	punch   float64 // Sustain punch (proportion)
	decay   float64 // Decay time (seconds)

	sampleRate int // Hz
	sampleSize int // bits per channel
}
