package main

import (
	"github.com/wdevore/Ranger-Go-IGE/engine"
)

func main() {
	engine := engine.Construct("../../..")

	// Override some of the world properties for this example
	engine.World().PropertiesOverride("config.json")

	engine.Start()
}
