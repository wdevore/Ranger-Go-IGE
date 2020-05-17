package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
)

// World is the main component of ranger
type world struct {
	properties   *configuration.Properties
	relativePath string

	shader api.IShader

	staticAtlas  api.IAtlas
	dynamicAtlas api.IAtlas

	rasterFont api.IRasterFont

	renderIdx  int
	renderRepo map[int]api.IRenderGraphic

	activeRenGID int
	activeRenG   api.IRenderGraphic

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
	w.renderRepo = make(map[int]api.IRenderGraphic)
	w.renderIdx = api.GlobalRenderGraphic // Default

	shp := w.properties.Shaders

	// ---------------------------------------
	// Compile shader
	// ---------------------------------------
	w.shader = rendering.NewShaderFromCode(shp.VertexShaderCode, shp.FragmentShaderCode)

	err := w.shader.Compile()

	if err != nil {
		fmt.Println("RenderGraphic error: ")
		panic(err)
	}

	// Activate shader so we can query it.
	w.shader.Use()

	// Create a graphic that will store Static shapes
	// pass functor for populating
	w.staticAtlas = rendering.NewStaticAtlas()
	renG := rendering.NewRenderGraphic(true, w.staticAtlas, w.shader)
	w.AddRenderGraphic(renG)

	w.dynamicAtlas = rendering.NewDynamicAtlas()
	renG = rendering.NewRenderGraphic(false, w.dynamicAtlas, w.shader)
	w.AddRenderGraphic(renG)

	// Force UseRenderGraphic to UnUse/Use for the first node visited
	w.activeRenGID = -1

	fmt.Println("Loading Raster font...")
	w.rasterFont = fonts.NewRasterFont()
	err = w.rasterFont.Initialize("raster_font.data", w.relativePath)

	return err
}

func (w *world) GenGraphicID() int {
	id := w.renderIdx
	w.renderIdx++
	return id
}

func (w *world) AddRenderGraphic(graphic api.IRenderGraphic) int {
	w.activeRenG = graphic
	id := w.GenGraphicID()
	w.renderRepo[id] = graphic
	return id
}

func (w *world) GetRenderGraphic(graphicID int) api.IRenderGraphic {
	return w.renderRepo[graphicID]
}

func (w *world) UseRenderGraphic(graphicID int) api.IRenderGraphic {
	if w.activeRenGID != graphicID {
		// Deactivate current graphic
		w.activeRenG.UnUse()

		// Activate new graphic
		fmt.Println("graphicID: ", graphicID)
		w.activeRenGID = graphicID
		w.activeRenG = w.renderRepo[graphicID]
		w.activeRenG.Use()
	}

	return w.activeRenG
}

func (w *world) Atlas() api.IAtlas {
	// return w.activeRenG.Atlas()
	return w.staticAtlas
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
