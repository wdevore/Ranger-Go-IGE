package api

import (
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
)

// IWorld represents properties of the game world
type IWorld interface {
	Configure() error
	Begin() error
	End()
	RelativePath() string

	NodeManager() INodeManager
	Properties() *configuration.Properties
	PropertiesOverride(configFiel string)

	Root() INode
	Underlay() INode
	Scenes() INode
	Overlay() INode

	Push(scene INode)

	RouteEvents(event IEvent)

	// Deprecate
	GenGraphicID() int
	AddRenderGraphic(graphic IRenderGraphic, graphicID int)
	GetRenderGraphic(graphicID int) IRenderGraphic
	UseRenderGraphic(graphicID int) IRenderGraphic
	SwitchRenderGraphic(graphicID int) IRenderGraphic
	SetRenderGraphic(graphicID int)

	Projection() IMatrix4
	Viewspace() IMatrix4
	InvertedViewspace() IMatrix4

	AddAtlas(name string, atlas IAtlasX)
	GetAtlas(name string) IAtlasX

	// Deprecate
	Atlas() IAtlas        // Static atlas
	DynoAtlas() IAtlas    // DynamicMulti atlas
	PixelAtlas() IAtlas   // Dynamic atlas
	TextureAtlas() IAtlas // Texture atlas

	PostProcess()

	Shader() IShader

	TextureShader() IShader
	TextureManager() ITextureManager

	RasterFont() IRasterFont

	// Debug Info
	Fps() int
	SetFps(int)

	Ups() int
	SetUps(int)

	AvgRender() float64
	SetAvgRender(float64)
}
