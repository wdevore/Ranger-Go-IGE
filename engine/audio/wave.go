package audio

// !!!!! DEPRECATED in favor of https://github.com/faiface/beep

const (
	// Byte offsets                           // Size  Comment
	waveChunkID          = 0  // 4     "RIFF" = 0x52494646
	waveChunkSize        = 4  // 4     36+SubChunk2Size = 4+(8+SubChunk1Size)+(8+SubChunk2Size)
	waveFormat           = 8  // 4     "WAVE" = 0x57415645
	waveSubChunk1Index   = 12 // 4     "fmt " = 0x666d7420
	waveSubChunk1Size    = 16 // 4     16 = 0x10 for PCM
	waveAudioFormat      = 20 // 2     PCM = 0x01
	waveNumberOfChannels = 22 // 2     Mono = 0x01, Stereo = 0x02, etc.
	waveSampleRate       = 24 // 4     8000, 44100, etc
	waveByteRate         = 28 // 4     SampleRate * NumChannels * BitsPerSample / 8
	waveBlockAlign       = 32 // 2     NumChannels * BitsPerSample / 8
	waveBitsPerSample    = 34 // 2     8 bits = 0x08, 16 bits = 0x10, etc.
	waveSubChunk2Index   = 36 // 4     "data" = 0x64617461
	waveSubChunk2Size    = 40 // 4     data size = NumSamples * NumChannels * BitsPerSample / 8
	waveData             = 44 // N
	waveHeaderSize       = 44
	// 4 byte hex ascii strings.
	waveSRift = 0x52494646
	waveSWave = 0x57415645
	waveSFmt  = 0x666d7420
	waveSData = 0x64617461
)

type wave struct {
	// MIME http://en.wikipedia.org/wiki/Data_URI_scheme
	/// The MIME base64 encoded wave
	dataURI string

	subChunk1Size int

	/// Default rate is 8000
	sampleRate int
	/// Default to PCM = 1
	audioFormat int
	/// Default is 1 channel
	numberOfChannels int
	/// Default is 8
	bitsPerSample int

	/// {Informational} Simply indicates if clipping occurred and how many times.
	clipping int
}

func newWave() *wave {
	o := new(wave)

	o.subChunk1Size = 16
	o.sampleRate = 8000
	o.audioFormat = 1
	o.numberOfChannels = 1
	o.bitsPerSample = 8
	o.clipping = 0

	return o
}

func (w *wave) Create(data *[]int) {

}
