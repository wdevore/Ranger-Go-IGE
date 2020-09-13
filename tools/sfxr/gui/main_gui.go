package gui

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
)

var (
// textBuffer = ""
)

// DrawGui ...
func DrawGui(config *settings.ConfigJSON, generator api.ISampleGenerator) {
	BuildMenuBar(config, generator)
	DrawGeneratorsPanel(config, generator)
	BuildWaveformPanel(config)
	BuildSettingsPanel(config, generator)
}
