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
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
	"github.com/wdevore/Ranger-Go-IGE/engine/textures"
)

// World is the main component of ranger
type world struct {
	properties   *configuration.Properties
	relativePath string

	defaultShader api.IShader
	textureShader api.IShader

	staticAtlas   api.IAtlas
	dynamicAtlas  api.IAtlas
	pixelAtlas    api.IAtlas
	textureAtlas  api.IAtlas
	texture2Atlas api.IAtlas

	rasterFont api.IRasterFont

	renderIdx  int
	renderRepo map[int]api.IRenderGraphic

	activeRenGID int
	activeRenG   api.IRenderGraphic

	viewSpace    api.IMatrix4
	invViewSpace api.IMatrix4

	textureMan api.ITextureManager

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

	shaders := &o.properties.Shaders

	if shaders.UseDefault {
		code, err := o.loadShaderSource(dataPath, shaders.VertexShaderSrc)
		if err != nil {
			panic(err)
		}
		shaders.VertexShaderCode = *code

		code, err = o.loadShaderSource(dataPath, shaders.FragmentShaderSrc)
		if err != nil {
			panic(err)
		}
		shaders.FragmentShaderCode = *code
	}

	if shaders.UseTextureShader {
		// Texture shaders
		code, err := o.loadShaderSource(dataPath, shaders.TextureVertexShaderSrc)
		if err != nil {
			panic(err)
		}
		shaders.TextureVertexShaderCode = *code

		code, err = o.loadShaderSource(dataPath, shaders.TextureFragmentShaderSrc)
		if err != nil {
			panic(err)
		}
		shaders.TextureFragmentShaderCode = *code
	}

	o.textureMan = textures.NewTextureManager()

	return o
}

func (w *world) loadShaderSource(dataPath, shaderSrc string) (code *string, err error) {
	file, err := os.Open(dataPath + "/engine/assets/shaders/" + shaderSrc)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	bytes = append(bytes, 0)
	sCode := string(bytes)

	return &sCode, nil
}

func (w *world) Configure() error {
	w.renderRepo = make(map[int]api.IRenderGraphic)
	w.renderIdx = api.StaticRenderGraphic // Default

	w.viewSpace = maths.NewMatrix4()
	w.invViewSpace = maths.NewMatrix4()

	shaderP := w.properties.Shaders

	// ---------------------------------------
	// Compile shader
	// ---------------------------------------
	w.defaultShader = rendering.NewShaderFromCode(shaderP.VertexShaderCode, shaderP.FragmentShaderCode)

	err := w.defaultShader.Compile()

	if err != nil {
		fmt.Println("Default Shader compile error: ")
		panic(err)
	}

	// Activate shader so we can query it.
	w.defaultShader.Use()

	// Create a graphic that will store Static shapes
	w.staticAtlas = rendering.NewShapeAtlas()
	renG := rendering.NewRenderGraphic(w.staticAtlas, w.defaultShader)
	w.AddRenderGraphic(renG, api.StaticRenderGraphic)

	w.pixelAtlas = rendering.NewShapeAtlas()
	renG = rendering.NewRenderGraphic(w.pixelAtlas, w.defaultShader)
	w.AddRenderGraphic(renG, api.DynamicPixelBufRenderGraphic)

	w.dynamicAtlas = rendering.NewShapeAtlas()
	renG = rendering.NewRenderGraphic(w.dynamicAtlas, w.defaultShader)
	w.AddRenderGraphic(renG, api.DynamicRenderGraphic)

	// --------------------------------------------------------------
	// A Texture renderer
	w.textureShader = rendering.NewShaderFromCode(shaderP.TextureVertexShaderCode, shaderP.TextureFragmentShaderCode)

	err = w.textureShader.Compile()

	if err != nil {
		fmt.Println("Texture Shader compile error: ")
		panic(err)
	}
	w.textureAtlas = rendering.NewShapeAtlas()
	renG = rendering.NewTextureRenderGraphic(w.textureAtlas, w.textureShader)
	w.AddRenderGraphic(renG, api.TextureRenderGraphic)
	// ---------------------------------------------------------------

	// Force UseRenderGraphic to UnUse/Use for the first node visited
	w.activeRenGID = -1

	fmt.Println("Loading Raster font...")
	w.rasterFont = fonts.NewRasterFont()
	err = w.rasterFont.Initialize("raster_font.data", w.relativePath)

	return err
}

func (w *world) Viewspace() api.IMatrix4 {
	return w.viewSpace
}

func (w *world) InvertedViewspace() api.IMatrix4 {
	return w.invViewSpace
}

func (w *world) GenGraphicID() int {
	id := w.renderIdx
	w.renderIdx++
	return id
}

func (w *world) AddRenderGraphic(graphic api.IRenderGraphic, graphicID int) {
	w.activeRenG = graphic
	// id := w.GenGraphicID()
	w.renderRepo[graphicID] = graphic
}

func (w *world) GetRenderGraphic(graphicID int) api.IRenderGraphic {
	return w.renderRepo[graphicID]
}

func (w *world) UseRenderGraphic(graphicID int) api.IRenderGraphic {
	if w.activeRenGID != graphicID {
		w.SwitchRenderGraphic(graphicID)
	}

	return w.activeRenG
}

func (w *world) SwitchRenderGraphic(graphicID int) api.IRenderGraphic {
	// Deactivate current graphic
	w.activeRenG.UnUse()

	// Activate new graphic
	// fmt.Println("graphicID: ", graphicID)
	w.activeRenGID = graphicID
	w.activeRenG = w.renderRepo[graphicID]
	w.activeRenG.Use()

	return w.activeRenG
}

func (w *world) PostProcess() {
	// All atlases are copied to OpenGL via VBO/TBO and EBO bindings

	if w.staticAtlas.HasShapes() {
		renG := w.renderRepo[api.StaticRenderGraphic]
		renG.Construct(api.MeshStatic, w.staticAtlas)
	}

	if w.pixelAtlas.HasShapes() {
		renG := w.renderRepo[api.DynamicPixelBufRenderGraphic]
		renG.Construct(api.MeshDynamic, w.pixelAtlas)
	}

	if w.dynamicAtlas.HasShapes() {
		renG := w.renderRepo[api.DynamicRenderGraphic]
		renG.Construct(api.MeshDynamicMulti, w.dynamicAtlas)
	}
}

func (w *world) Atlas() api.IAtlas {
	return w.staticAtlas
}

func (w *world) DynoAtlas() api.IAtlas {
	return w.dynamicAtlas
}

func (w *world) PixelAtlas() api.IAtlas {
	return w.pixelAtlas
}

func (w *world) TextureAtlas() api.IAtlas {
	return w.textureAtlas
}

func (w *world) Shader() api.IShader {
	return w.defaultShader
}

func (w *world) TextureShader() api.IShader {
	return w.textureShader
}

func (w *world) TextureManager() api.ITextureManager {
	return w.textureMan
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
