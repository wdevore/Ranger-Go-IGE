package engine

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
)

// World is the main component of ranger
type world struct {
	properties *configuration.Properties

	shader api.IShader
}

func newWorld(relativePath string) api.IWorld {
	o := new(world)
	o.properties = &configuration.Properties{}

	dataPath, err := filepath.Abs(relativePath)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	eConfFile, err := os.Open(dataPath + "/engine/configuration/config.json")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer eConfFile.Close()

	bytes, err := ioutil.ReadAll(eConfFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	err = json.Unmarshal(bytes, o.properties)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	shp := o.properties.Shaders

	if shp.UseDefault {
		shp := &o.properties.Shaders

		vertexShaderFile, err := os.Open(dataPath + "/engine/assets/shaders/" + shp.VertexShaderSrc)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		defer vertexShaderFile.Close()

		bytes, err := ioutil.ReadAll(vertexShaderFile)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		shp.VertexShaderCode = string(bytes)

		fragmentShaderFile, err := os.Open(dataPath + "/engine/assets/shaders/" + shp.FragmentShaderSrc)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		defer fragmentShaderFile.Close()

		bytes, err = ioutil.ReadAll(fragmentShaderFile)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}

		shp.FragmentShaderCode = string(bytes)
	}

	return o
}

func (w *world) Configure() error {
	shp := w.properties.Shaders

	w.shader = rendering.NewShaderFromCode(shp.VertexShaderCode, shp.FragmentShaderCode)

	err := w.shader.Compile()

	return err
}

func (w *world) Shader() api.IShader {
	return w.shader
}

func (w *world) Properties() *configuration.Properties {
	return w.properties
}

func (w *world) PropertiesOverride(configFile string) {
	eConfFile, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer eConfFile.Close()

	bytes, err := ioutil.ReadAll(eConfFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	// Merge on top other existing property values
	err = json.Unmarshal(bytes, w.properties)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}
