package gui

import (
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
)

var (
	textBuffer = ""
	sampleRate int
	sampleSize int
	autoPlay   bool

	// SfxrJ holds the settings for an effect
	SfxrJ audio.SfxrJSON
)

// BuildGui ...
func BuildGui(config *settings.ConfigJSON) {
	BuildMenuBar(config)
	BuildGeneratorsPanel()
	BuildWaveformPanel()
	BuildSettingsPanel()

}
