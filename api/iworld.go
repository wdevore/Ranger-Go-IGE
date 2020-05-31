package api

import (
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
)

// IWorld represents properties of the game world
type IWorld interface {
	Configure() error

	Properties() *configuration.Properties
	PropertiesOverride(configFiel string)

	GenGraphicID() int
	AddRenderGraphic(graphic IRenderGraphic, graphicID int)
	GetRenderGraphic(graphicID int) IRenderGraphic
	UseRenderGraphic(graphicID int) IRenderGraphic

	Viewspace() IMatrix4
	InvertedViewspace() IMatrix4

	Atlas() IAtlas      // Static atlas
	DynoAtlas() IAtlas  // DynamicMulti atlas
	PixelAtlas() IAtlas // Dynamic atlas
	PostProcess()

	Shader() IShader

	RasterFont() IRasterFont

	// Debug Info
	Fps() int
	SetFps(int)

	Ups() int
	SetUps(int)

	AvgRender() float64
	SetAvgRender(float64)
}
