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
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/fonts"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

// World is the main component of ranger
type world struct {
	// -----------------------------------------
	// Scene graph is a node manager
	// -----------------------------------------
	sceneGraph api.INodeManager
	root       api.INode
	underlay   api.INode
	scenes     api.INode
	overlay    api.INode

	properties   *configuration.Properties
	relativePath string

	atlases map[string]api.IAtlasX

	rasterFont api.IRasterFont

	projection   api.IMatrix4
	viewSpace    api.IMatrix4
	invViewSpace api.IMatrix4

	// Debug Info
	fps        int
	ups        int
	renderTime float64
}

func newWorld(relativePath string) api.IWorld {
	o := new(world)

	o.sceneGraph = nodes.NewNodeManager()

	o.properties = &configuration.Properties{}
	o.relativePath = relativePath

	o.atlases = make(map[string]api.IAtlasX)

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

func (w *world) NodeManager() api.INodeManager {
	return w.sceneGraph
}

// Begin is called by the Engine during Construct(). The Under/Over lays may
// be populated afterwards.
func (w *world) Begin() error {
	var err error

	// The NodeManager needs to build a baseline INode structure for
	// the runtime environment. Structure:
	//
	//                            Root
	//         /-----------------/  | \-----------------\
	//         |                    |                   |
	//      Underlay              Scenes             Overlay
	//                           /     \
	//                   In:Scene       Out:Scene
	//
	// From here the NM's job is to Add/Remove Scenes from the Scenes-Node

	// Create Root first and above all (pun intended) do it NOW! ;-)
	w.root, err = extras.NewGroupNode("Root", w, nil)
	if err != nil {
		return err
	}

	w.underlay, err = extras.NewGroupNode("Underlay", w, w.root)
	if err != nil {
		return err
	}

	w.scenes, err = extras.NewGroupNode("Scenes", w, w.root)
	if err != nil {
		return err
	}

	w.overlay, err = extras.NewGroupNode("Overlay", w, w.root)
	if err != nil {
		return err
	}

	w.sceneGraph.SetRoot(w.root)

	return nil
}

func (w *world) End() {
	// So THIS is where the world actually comes to an END!
	w.sceneGraph.End()
}

func (w *world) RelativePath() string {
	return w.relativePath
}

func (w *world) Root() api.INode {
	return w.root
}
func (w *world) Underlay() api.INode {
	return w.underlay
}
func (w *world) Scenes() api.INode {
	return w.scenes
}
func (w *world) Overlay() api.INode {
	return w.overlay
}

func (w *world) Push(scene api.INode) {
	_, ok := scene.(api.IScene)
	if !ok {
		panic("Scene being pushed doesn't implementing IScene interface.")
	}

	w.sceneGraph.PushNode(scene)
}

func (w *world) RouteEvents(event api.IEvent) {
	w.NodeManager().RouteEvents(event)
}

func (w *world) Configure() error {

	w.viewSpace = maths.NewMatrix4()
	w.invViewSpace = maths.NewMatrix4()

	fmt.Println("Loading Raster font...")
	w.rasterFont = fonts.NewRasterFont()
	err := w.rasterFont.Initialize("raster_font.data", w.relativePath)

	if err != nil {
		return nil
	}

	// ------------------------------------------------------------
	// Projection space
	// ------------------------------------------------------------
	camera := w.Properties().Camera
	wp := w.Properties().Window

	projection := display.NewCamera()
	projection.SetProjection(
		0.0, 0.0, // bottom,left
		float32(wp.DeviceRes.Height), float32(wp.DeviceRes.Width), // top,right
		camera.Depth.Near, camera.Depth.Far)

	w.projection = projection.Matrix()
	return nil
}

func (w *world) Projection() api.IMatrix4 {
	return w.projection
}

func (w *world) Viewspace() api.IMatrix4 {
	return w.viewSpace
}

func (w *world) InvertedViewspace() api.IMatrix4 {
	return w.invViewSpace
}

func (w *world) AddAtlas(name string, atlas api.IAtlasX) {
	w.atlases[name] = atlas
}

func (w *world) GetAtlas(name string) api.IAtlasX {
	return w.atlases[name]
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
