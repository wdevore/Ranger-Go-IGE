package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/inkyblackness/imgui-go/v2"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/platforms"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/renderers"
	"github.com/wdevore/Ranger-Go-IGE/tools/sfxr/settings"
)

func main() {
	fmt.Println("Welcome to Deuron8 Go edition")

	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	// -------------------------------------------------------------
	// GLFW Window
	// -------------------------------------------------------------
	platform, err := platforms.NewGLFW("Go-Sfxr", 600, 900, io, platforms.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	// -------------------------------------------------------------
	// Renderer used by GLFW
	// -------------------------------------------------------------
	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	config := settings.ConfigJSON{}
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)

	if err != nil {
		panic(err)
	}

	// Finally run main gui application
	run(platform, renderer, &config)
}
