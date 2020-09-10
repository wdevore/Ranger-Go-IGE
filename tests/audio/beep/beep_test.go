package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// go test -v -count=1 beep_test.go

// Cleaned up original wave using:
// ffmpeg -i "RangerAudio.wav" -f wav -bitexact -acodec pcm_s16le -ar 22050 -ac 1 "ffmpeg.wav"

func TestRunner(t *testing.T) {
	noise2()
}

func ffmpeg() {
	f, err := os.Open("assets/ffmpeg.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	memoryBuffered(streamer, format)
}

func memoryBuffered(streamer beep.StreamSeekCloser, format beep.Format) {
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

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

func streamedDisc(streamer beep.StreamSeekCloser, format beep.Format) {

	defer streamer.Close()

	duration := streamer.Len() / format.SampleRate.N(time.Second)
	fmt.Println("Duration: ", duration)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	fmt.Println("Playing...")
	speaker.Play(streamer) // Async

	time.Sleep(time.Second * time.Duration(duration))

	fmt.Println("Done.")
}

func noise1() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(Noise{})
	select {}
}

func noise2() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Take(sr.N(2*time.Second), Noise{}), beep.Callback(func() {
		done <- true
	})))
	<-done
}

func noise3() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Take(sr.N(2*time.Second), Noise{}), beep.Callback(func() {
		done <- true
	})))
	<-done
}

func noise4() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(beep.Take(sr.N(2*time.Second), NoiseF()), beep.Callback(func() {
		done <- true
	})))
	<-done
}

func NoiseF() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = rand.Float64()*2 - 1
			samples[i][1] = rand.Float64()*2 - 1
		}
		return len(samples), true
	})
}

// ------------------------------------------------------
type Noise struct{}

func (no Noise) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = rand.Float64()*2 - 1
		samples[i][1] = rand.Float64()*2 - 1
	}
	return len(samples), true
}

func (no Noise) Err() error {
	return nil
}
