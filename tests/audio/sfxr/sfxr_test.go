package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
)

// go test -v -count=1 sfxr_test.go

func TestRunner(t *testing.T) {
	// sxfrLoad()
	// memoryBuf()
	waveStreamWriter()
}

func sxfrLoad() {
	audioJ := loadSfxr("../../../extras/sfxr/HighSproing_sfxr.json")
	values := audio.NewIntervalValues(audioJ)

	generator := audio.NewSfxrGenerator()
	generator.Init(values)
	generator.Generate(values)
	generator.CanBeDrained(true)
	samples := generator.Samples()
	fmt.Println("Sample size: ", len(*samples))

	// Buffers seem to require a precision of 2
	format := beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 2}
	buffer := beep.NewBuffer(format)
	buffer.Append(generator)
	generator = nil

	// fmt.Println("buff length: ", buffer.Len())
	sounds := buffer.Streamer(0, buffer.Len())
	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// done := make(chan bool)
	// speaker.Play(beep.Seq(sounds, beep.Callback(func() {
	// 	done <- true
	// })))
	// fmt.Println("Playing...")
	// <-done

	dataPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	waveF, err := os.Create(dataPath + "/" + "samples.wav")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer waveF.Close()
	wav.Encode(waveF, sounds, format)

	fmt.Println("Done.")

}

func waveStreamWriter() {
	values := audio.ConfigureTone(440, api.WaveSawtooth)
	values.SetSoundVol(0.1)

	generator := audio.NewSfxrGenerator()
	generator.Init(values)
	generator.Generate(values)
	generator.CanBeDrained(true)
	samples := generator.Samples()
	fmt.Println("Sample size: ", len(*samples))

	dataPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	waveF, err := os.Create(dataPath + "/" + "waveStrmTone.wav")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	format := beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 2}

	defer waveF.Close()
	wav.Encode(waveF, generator, format)

	fmt.Println("Done.")
}

func waveMemWriter() {
	values := audio.ConfigureTone(440, api.WaveSINE)

	generator := audio.NewSfxrGenerator()
	generator.Init(values)
	generator.Generate(values)
	samples := generator.Samples()
	fmt.Println("Sample size: ", len(*samples))

	format := beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 2}
	buffer := beep.NewBuffer(format)
	buffer.Append(generator)
	// buffer.Append(generator)
	// buffer.Append(generator)

	fmt.Println("buff length: ", buffer.Len())
	// sounds := buffer.Streamer(0, buffer.Len())

	dataPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	waveF, err := os.Create(dataPath + "/" + "waveTone.wav")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer waveF.Close()
	wav.Encode(waveF, generator, format)

	fmt.Println("Done.")
}

func memoryBuf() {

	rand.Seed(13163)
	values := audio.ConfigureTone(440, api.WaveSINE)
	// rand.Seed(13163)
	// values := audio.ConfigureExplosion(api.WaveNoise)
	// rand.Seed(1)
	// values := audio.ConfigurePickupCoin(api.WaveSINE, true)
	// rand.Seed(13163)
	// values := audio.ConfigureLaserShoot()
	// rand.Seed(6666)
	// values := audio.ConfigurePowerUp()
	// rand.Seed(10101)
	// values := audio.ConfigureHitHurt()
	// rand.Seed(10101)
	// values := audio.ConfigureJump()
	// rand.Seed(10101)
	// values := audio.ConfigureBlipSelect()

	// rand.Seed(10101) // Chimmy rampy
	// rand.Seed(666) // Thumping humm
	// rand.Seed(31) // Alarm buzz
	// rand.Seed(1) // Soft cricket
	// rand.Seed(2) // Soft explosion
	// rand.Seed(3) // High alarm
	// rand.Seed(7) // Whip + low vibrate
	// rand.Seed(8) // Zapper
	// rand.Seed(10) // 2tone warble alarm (long)
	// rand.Seed(11) // Mean phaser
	// rand.Seed(12) // Slide downward buzz alarm
	// rand.Seed(13) // Short static + low tone
	// rand.Seed(14) // Alarm clock
	// rand.Seed(15) // Sweep blip + lowF chime
	// rand.Seed(16) // Sci-Fi eery high wobble
	// rand.Seed(18) // Short noise + Alarm
	// rand.Seed(19) // Chirp + high ching
	// rand.Seed(21) // High + Long Low
	// values := audio.ConfigureRandom()

	generator := audio.NewSfxrGenerator()
	generator.Init(values)
	generator.Generate(values)
	generator.CanBeDrained(true)
	samples := generator.Samples()
	fmt.Println("Sample size: ", len(*samples))

	// Buffers seem to require a precision of 2
	format := beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 2}
	buffer := beep.NewBuffer(format)
	buffer.Append(generator)
	buffer.Append(generator)
	buffer.Append(generator)
	generator = nil

	fmt.Println("buff length: ", buffer.Len())
	sounds := buffer.Streamer(0, buffer.Len())

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(sounds, beep.Callback(func() {
		done <- true
	})))
	fmt.Println("Playing...")
	<-done

	fmt.Println("Done.")

}

func streamer() {
	values := audio.ConfigureTone(440, api.WaveSINE)

	generator := audio.NewSfxrGenerator()
	generator.Init(values)
	generator.Generate(values)

	samples := generator.Samples()
	fmt.Println("Sample size: ", len(*samples))

	format := beep.Format{SampleRate: 44100, NumChannels: 1, Precision: 1}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Take(format.SampleRate.N(1*time.Second), generator), beep.Callback(func() {
		done <- true
	})))
	fmt.Println("Playing...")
	<-done

	// dataPath, err := filepath.Abs(".")
	// if err != nil {
	// 	log.Fatalln("ERROR:", err)
	// }

	// waveF, err := os.Create(dataPath + "/" + "waveStreammTone.wav")
	// if err != nil {
	// 	log.Fatalln("ERROR:", err)
	// }

	// defer waveF.Close()
	// wav.Encode(waveF, generator, format)

	fmt.Println("Done.")
}

func loadSfxr(sxfrFile string) *audio.SfxrJSON {
	audioJ := &audio.SfxrJSON{}

	dataPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	eConfFile, err := os.Open(dataPath + "/" + sxfrFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer eConfFile.Close()

	bytes, err := ioutil.ReadAll(eConfFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	err = json.Unmarshal(bytes, audioJ)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	return audioJ
}
