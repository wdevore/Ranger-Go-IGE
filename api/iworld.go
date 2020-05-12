package api

import (
	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
)

// IWorld represents properties of the game world
type IWorld interface {
	Configure() error

	Properties() *configuration.Properties
	PropertiesOverride(configFiel string)

	Shader() IShader
	// Shader locations
	ModelLoc() int32
	ColorLoc() int32

	Atlas() IAtlas
	VecObj() IVectorObject

	RasterFont() IRasterFont
}