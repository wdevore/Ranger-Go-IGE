package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
)

var (
	showOpenDialog = false
	filePath       = ""
	absolutePath   = ""
	relativePath   = ""
)

// BuildMenuBar --
func BuildMenuBar(config *settings.ConfigJSON) {
	// ---------------------------------------------------------
	// Build the application GUI
	// ---------------------------------------------------------

	if imgui.BeginMainMenuBar() {
		if imgui.BeginMenu("File") {
			if imgui.MenuItem("New") {
			}

			if imgui.MenuItem("Open Sfxr") {
				// Show open dialog
				showOpenDialog = true
				absolutePath, _ = filepath.Abs(".")
				relativePath = config.RootPath
			}

			if imgui.MenuItem("Save") {
			}

			if imgui.MenuItem("Save As Sfxr") {
			}

			if imgui.MenuItem("Save As Wav") {
			}

			if imgui.MenuItem("Exit") {
				// Save application property settings.
				// config.Save()
				// environment.IssueCmd("killSim")
				os.Exit(0)
			}

			imgui.EndMenu()
		}

		// ------------------------------------------------------
		if imgui.BeginMenu("Go") {
			changed := imgui.Checkbox("AutoPlay", &autoPlay)
			if changed {
				if autoPlay {
					fmt.Println("AutoPlay enabled")
				} else {
					fmt.Println("AutoPlay disabled")
				}
			}

			if imgui.MenuItem("Play") {
			}

			if imgui.MenuItem("Copy") {
			}

			imgui.EndMenu()
		}

		b := false

		if imgui.BeginMenu("SampleRate(Hz)") {
			b = sampleRate == 0
			changed := imgui.Checkbox("44k", &b)
			if changed {
				sampleRate = 0
			}

			b = sampleRate == 1
			changed = imgui.Checkbox("22k", &b)
			if changed {
				sampleRate = 1
			}

			b = sampleRate == 2
			changed = imgui.Checkbox("11k", &b)
			if changed {
				sampleRate = 2
			}

			b = sampleRate == 3
			changed = imgui.Checkbox("6k", &b)
			if changed {
				sampleRate = 3
			}

			imgui.EndMenu()
		}

		if imgui.BeginMenu("SampleSize") {
			b = sampleSize == 0
			changed := imgui.Checkbox("8 bit", &b)
			if changed {
				sampleSize = 0
			}

			b = sampleSize == 1
			changed = imgui.Checkbox("16 bit", &b)
			if changed {
				sampleSize = 1
			}

			imgui.EndMenu()
		}

		// ------------------------------------------------------

		imgui.EndMainMenuBar()
	}

	if showOpenDialog {
		// Pass a pointer to our bool variable (the window will have a closing button that will clear the bool when clicked)
		imgui.SetNextWindowSize(imgui.Vec2{X: 500, Y: 100})
		imgui.BeginV("Open Sfxr", &showOpenDialog, 0)

		imgui.PushItemWidth(450)
		imgui.InputText("Path", &relativePath)
		filePath = config.DefaultFile

		if config.LastOpenedFile != "" {
			filePath = config.LastOpenedFile
		}
		imgui.InputText("File", &filePath)

		if imgui.Button("Open") {
			showOpenDialog = false

			file := relativePath + "/" + filePath + "." + config.SfxrExtention
			sfxrFile, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer sfxrFile.Close()

			config.LastOpenedFile = filePath

			bytes, err := ioutil.ReadAll(sfxrFile)
			if err != nil {
				panic(err)
			}

			SfxrJ = audio.SfxrJSON{}

			err = json.Unmarshal(bytes, &SfxrJ)
			if err != nil {
				panic(err)
			}

			fmt.Println("Opened: ", file)
		}

		imgui.SameLine()
		if imgui.Button("Cancel") {
			showOpenDialog = false
		}

		imgui.PopItemWidth()

		imgui.End()
	}
}
