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

	Projection() IMatrix4
	Viewspace() IMatrix4
	InvertedViewspace() IMatrix4

	RasterFont() IRasterFont
	AddAtlas(name string, atlas IAtlasX)
	GetAtlas(name string) IAtlasX

	// Debug Info
	Fps() int
	SetFps(int)

	Ups() int
	SetUps(int)

	AvgRender() float64
	SetAvgRender(float64)
}
