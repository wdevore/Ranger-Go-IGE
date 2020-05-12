package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
)

// World is the main component of ranger
type world struct {
	properties   *configuration.Properties
	relativePath string

	shader   api.IShader
	modelLoc int32
	colorLoc int32

	vecObj api.IVectorObject
	atlas  api.IAtlas

	rasterFont api.IRasterFont

	// Debug Info
	fps        int
	ups        int
	renderTime float64
}

func newWorld(relativePath string) api.IWorld {
	o := new(world)
	o.properties = &configuration.Properties{}
	o.relativePath = relativePath

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

	if err != nil {
		return err
	}

	// Construct and populate the vector shape atlas
	w.vecObj = rendering.NewVectorObject()
	w.vecObj.Construct()

	w.atlas = atlas.NewAtlas()
	w.atlas.Initialize(w.vecObj)

	w.vecObj.Bind()

	w.shader.Use()

	programID := w.shader.Program()

	w.modelLoc = gl.GetUniformLocation(programID, gl.Str("model\x00"))
	if w.modelLoc < 0 {
		return errors.New("World: couldn't find 'model' uniform variable")
	}

	w.colorLoc = gl.GetUniformLocation(programID, gl.Str("fragColor\x00"))
	if w.colorLoc < 0 {
		return errors.New("World: couldn't find 'fragColor' uniform variable")
	}

	fmt.Println("Loading Raster font...")
	w.rasterFont = fonts.NewRasterFont()
	err = w.rasterFont.Initialize("raster_font.data", w.relativePath)

	return err
}

func (w *world) Shader() api.IShader {
	return w.shader
}

func (w *world) ModelLoc() int32 {
	return w.modelLoc
}

func (w *world) ColorLoc() int32 {
	return w.colorLoc
}

func (w *world) Atlas() api.IAtlas {
	return w.atlas
}

func (w *world) VecObj() api.IVectorObject {
	return w.vecObj
}

func (w *world) RasterFont() api.IRasterFont {
	return w.rasterFont
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

func (w *world) Fps() int {
	return w.fps
}

func (w *world) SetFps(v int) {
	w.fps = v
}

func (w *world) Ups() int {
	return w.ups
}

func (w *world) SetUps(v int) {
	w.ups = v
}

func (w *world) AvgRender() float64 {
	return w.renderTime
}

func (w *world) SetAvgRender(v float64) {
	w.renderTime = v
}
