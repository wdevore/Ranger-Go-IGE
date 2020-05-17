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
	AddRenderGraphic(graphic IRenderGraphic) int
	GetRenderGraphic(graphicID int) IRenderGraphic
	UseRenderGraphic(graphicID int) IRenderGraphic

	Atlas() IAtlas
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
