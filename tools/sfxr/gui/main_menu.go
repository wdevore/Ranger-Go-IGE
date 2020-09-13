package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/audio"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/sound"
)

var (
	showOpenDialog     = false
	showAboutDialog    = false
	showSaveSfxrDialog = false

	hasOpenBeenShown     = false
	hasSaveSfxrBeenShown = false

	inputBuffer   = ""
	inputFilePath = ""
	absolutePath  = ""
	relativePath  = ""

	fileIsPresent bool

	aboutText = "Go edition of `Sfxr by DrPetter`\nPorted to Go By Will Cleveland"
)

// BuildMenuBar --
func BuildMenuBar(config *settings.ConfigJSON, generator api.ISampleGenerator) {
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

			if imgui.MenuItem("Save to Sfxr") {
				showSaveSfxrDialog = true
				absolutePath, _ = filepath.Abs(".")
				relativePath = config.RootPath
			}

			if imgui.MenuItem("Save to Wav") {
			}

			if imgui.MenuItem("Exit") {
				// Save application property settings.
				indentedJSON, _ := json.MarshalIndent(*config, "", "  ")

				err := ioutil.WriteFile("config.json", indentedJSON, 0644)
				if err != nil {
					panic(err)
				}

				os.Exit(0)
			}

			imgui.EndMenu()
		}

		// ------------------------------------------------------
		if imgui.BeginMenu("Go") {
			changed := imgui.Checkbox("AutoPlay", &config.Autoplay)
			if changed {
				if config.Autoplay {
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
			b = sound.SfxrJ.SampleRate == 44100
			changed := imgui.Checkbox("44k", &b)
			if changed {
				sound.SfxrJ.SampleRate = 44100
			}

			b = sound.SfxrJ.SampleRate == 22050
			changed = imgui.Checkbox("22k", &b)
			if changed {
				sound.SfxrJ.SampleRate = 22050
			}

			b = sound.SfxrJ.SampleRate == 11025
			changed = imgui.Checkbox("11k", &b)
			if changed {
				sound.SfxrJ.SampleRate = 11025
			}

			b = sound.SfxrJ.SampleRate == 5512
			changed = imgui.Checkbox("5.5k", &b)
			if changed {
				sound.SfxrJ.SampleRate = 5512
			}

			imgui.EndMenu()
		}

		if imgui.BeginMenu("SampleSize") {
			b = sound.SfxrJ.SampleSize == 8
			changed := imgui.Checkbox("8 bit", &b)
			if changed {
				sound.SfxrJ.SampleSize = 8
			}

			b = sound.SfxrJ.SampleSize == 16
			changed = imgui.Checkbox("16 bit", &b)
			if changed {
				sound.SfxrJ.SampleSize = 16
			}

			imgui.EndMenu()
		}

		if imgui.BeginMenu("Help") {

			if imgui.MenuItem("About") {
				showAboutDialog = true
			}

			imgui.EndMenu()
		}

		// ------------------------------------------------------

		imgui.EndMainMenuBar()
	}

	if showAboutDialog {
		imgui.SetNextWindowSize(imgui.Vec2{X: 250, Y: 150})
		imgui.StyleColorsLight()
		imgui.BeginV("About", &showAboutDialog, 0)

		imgui.PushItemWidth(250)
		imgui.InputTextMultiline("", &aboutText)
		imgui.PopItemWidth()

		// imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
		// if imgui.Button("Ok") {
		// 	showAboutDialog = false
		// }
		// imgui.PopStyleColor()

		imgui.End()
		imgui.StyleColorsDark()

	}

	if showOpenDialog {
		drawOpenDialog(config, generator)
	}

	if showSaveSfxrDialog {
		drawSaveSfxrDialog(config)
	}
}

func drawSaveSfxrDialog(config *settings.ConfigJSON) {
	imgui.SetNextWindowSize(imgui.Vec2{X: 500, Y: 100})

	imgui.StyleColorsLight()
	imgui.BeginV("Save to Sfxr", &showSaveSfxrDialog, 0)

	imgui.PushItemWidth(450)
	imgui.InputText("Path", &relativePath)

	if !hasSaveSfxrBeenShown {
		if config.DefaultFile != "" {
			inputFilePath = config.DefaultFile
		}
		hasSaveSfxrBeenShown = true
	}

	if !fileIsPresent {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.5, Y: 0.5, Z: 0.5, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 0.5, Z: 0.0, W: 1.0})
	}
	imgui.InputText("File", &inputFilePath)
	imgui.PopStyleColor()

	// Does the file actually exist
	file := relativePath + "/" + inputFilePath + "." + config.SfxrExtention

	if _, err := os.Stat(file); err == nil {
		fileIsPresent = true
	} else {
		fileIsPresent = false
	}

	if imgui.Button("Save") {
		fmt.Println("Saveing... ", file)
		showSaveSfxrDialog = false

		file := relativePath + "/" + inputFilePath + "." + config.SfxrExtention

		// Transfer GValues to SfxrJ
		sound.SfxrJ.CopyFrom(sound.GValues)

		indentedJSON, _ := json.MarshalIndent(sound.SfxrJ, "", "  ")

		err := ioutil.WriteFile(file, indentedJSON, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Saved: ", file)
	}

	imgui.SameLine()

	// imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.5, Y: 0.5, Z: 0.5, W: 1.0})
	if imgui.Button("Cancel") {
		showSaveSfxrDialog = false
	}
	// imgui.PopStyleColor()

	imgui.PopItemWidth()

	imgui.End()
	imgui.StyleColorsDark()
}

func drawOpenDialog(config *settings.ConfigJSON, generator api.ISampleGenerator) {
	// Pass a pointer to our bool variable (the window will have a closing button that will clear the bool when clicked)
	imgui.SetNextWindowSize(imgui.Vec2{X: 500, Y: 100})
	imgui.StyleColorsLight()

	imgui.BeginV("Open Sfxr", &showOpenDialog, 0)

	imgui.PushItemWidth(450)

	// ----------------------------------------------------
	if imgui.InputText("Path", &relativePath) {
	}

	if !hasOpenBeenShown {
		if config.LastOpenedFile != "" {
			inputFilePath = config.LastOpenedFile
		}
		hasOpenBeenShown = true
	}

	// ----------------------------------------------------
	if fileIsPresent {
		imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{X: 0.5, Y: 0.5, Z: 0.5, W: 1.0})
	} else {
		imgui.PushStyleColor(imgui.StyleColorText, imgui.Vec4{X: 1.0, Y: 0.0, Z: 0.0, W: 1.0})
	}

	if imgui.InputText("File", &inputFilePath) {
	}
	imgui.PopStyleColor()

	// Does the file actually exist
	file := relativePath + "/" + inputFilePath + "." + config.SfxrExtention

	if _, err := os.Stat(file); err == nil {
		fileIsPresent = true
	} else {
		fileIsPresent = false
	}

	// ----------------------------------------------------
	if fileIsPresent {
		if imgui.Button("Open") {
			showOpenDialog = false

			file := relativePath + "/" + inputFilePath + "." + config.SfxrExtention
			sfxrFile, err := os.Open(file)
			if err != nil {
				panic(err)
			}
			defer sfxrFile.Close()

			config.LastOpenedFile = inputFilePath

			bytes, err := ioutil.ReadAll(sfxrFile)
			if err != nil {
				panic(err)
			}

			sound.SfxrJ = audio.SfxrJSON{}

			err = json.Unmarshal(bytes, &sound.SfxrJ)
			if err != nil {
				panic(err)
			}

			config.LastOpenedFile = inputFilePath
			fmt.Println("Opened: ", file)

			// Transfer to IGeneratorValues
			sound.GValues = audio.NewIntervalValues(&sound.SfxrJ)
			generator.Init(sound.GValues)

			sound.UpdateSfxrData(sound.GValues)
			sound.Generate(sound.GValues, generator)
			sound.Play(generator)
		}
		imgui.SameLine()
	}

	if imgui.Button("Cancel") {
		showOpenDialog = false
	}

	imgui.PopItemWidth()

	imgui.End()
	imgui.StyleColorsDark()

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
