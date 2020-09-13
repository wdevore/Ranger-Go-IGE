package gui

import (
	"math/rand"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/sound"
)

var (
	hasVolumeBeenSet = false
)

// DrawGeneratorsPanel ...
func DrawGeneratorsPanel(config *settings.ConfigJSON, generator api.ISampleGenerator) {
	imgui.SetNextWindowPos(imgui.Vec2{X: 10, Y: 20.0})
	imgui.SetNextWindowSize(imgui.Vec2{X: 115, Y: 420})

	imgui.Begin("Generators")

	if imgui.Button("Pickup/Coin") {
		// Generate a new set of values and feed to generator.
		sound.GValues = audio.ConfigurePickupCoin(true)
		generator.Init(sound.GValues)

		// Update GUI
		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Laser/Shoot") {
		sound.GValues = audio.ConfigureLaserShoot()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Explosion") {
		sound.GValues = audio.ConfigureExplosion()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("PowerUp") {
		sound.GValues = audio.ConfigurePowerUp()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Hit/Hurt") {
		sound.GValues = audio.ConfigureHitHurt()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Blip/Select") {
		sound.GValues = audio.ConfigureBlipSelect()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Synth") {
		sound.GValues = audio.ConfigureSynth()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Random") {
		sound.GValues = audio.ConfigureRandom()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Tone") {
		tone := 10.0 + rand.Float64()*4186.0

		sound.GValues = audio.ConfigureTone(tone, sound.SfxrJ.WaveShape)
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	if imgui.Button("Mutate") {
		sound.GValues.Mutate()
		generator.Init(sound.GValues)

		sound.UpdateSfxrData(sound.GValues)

		sound.Generate(sound.GValues, generator)

		sound.Play(generator)
	}

	imgui.SetNextItemOpen(true, imgui.ConditionOnce)
	if imgui.CollapsingHeader("Control") {
		imgui.IndentV(20)
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.9, Y: 0.5, Z: 0.0, W: 1.0})
		if imgui.Button(" Play ") {
			sound.GValues.SetSoundVol(sound.SfxrJ.SoundVolume)
			generator.Init(sound.GValues)
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
		imgui.PopStyleColor()

		imgui.Unindent()
		imgui.Separator()

		h := float32(sound.SfxrJ.SoundVolume)
		if !hasVolumeBeenSet {
			// On start of app set a default volume
			sound.SfxrJ.SoundVolume = api.PlaybackSoundVolume
			h = float32(sound.SfxrJ.SoundVolume)
			hasVolumeBeenSet = true
		}
		imgui.IndentV(10)
		if imgui.VSliderFloatV("Volume", imgui.Vec2{X: 30, Y: 100}, &h, 0.0, 10.0, "%.1f", 1.5) {
			sound.SfxrJ.SoundVolume = float64(h)
			sound.GValues.SetSoundVol(sound.SfxrJ.SoundVolume)
		}
	}

	imgui.End()
}
