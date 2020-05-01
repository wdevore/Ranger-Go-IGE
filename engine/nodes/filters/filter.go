package filters

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// #############################################################################
// Warning: Filters can interfer with things like dragging which
// rely on parent transform properties.
// If you need hiearchial dragging then skip using Filter and just
// manually manage the child transform properties relative to their parent.
// #############################################################################

// Filter is the base property of Filter nodes
type Filter struct {
	// The node's immediate parent translation components
	components api.IAffineTransform

	// What to exclude from the parent
	excludeTranslation bool
	excludeRotation    bool
	excludeScale       bool
}

func (f *Filter) initializeFilter() {
	f.components = maths.NewTransform()
}

// InheritAll causes the filter to pass all of the parent's transform
// properties: Translate, Rotation and Scale.
func (f *Filter) InheritAll() {
	f.excludeTranslation = false
	f.excludeRotation = false
	f.excludeScale = false
}

// InheritOnlyRotation causes the filter to pass only the parent's rotational
// property.
func (f *Filter) InheritOnlyRotation() {
	f.excludeTranslation = true
	f.excludeRotation = false
	f.excludeScale = true
}

// InheritOnlyScale causes the filter to pass only the parent's scale
// property.
func (f *Filter) InheritOnlyScale() {
	f.excludeTranslation = true
	f.excludeRotation = true
	f.excludeScale = false
}

// InheritOnlyTranslation causes the filter to pass only the parent's translation
// property.
func (f *Filter) InheritOnlyTranslation() {
	f.excludeTranslation = false
	f.excludeRotation = true
	f.excludeScale = true
}

// InheritRotationAndTranslation causes the filter to pass the parent's translation
// and rotational properties.
func (f *Filter) InheritRotationAndTranslation() {
	f.excludeTranslation = false
	f.excludeRotation = false
	f.excludeScale = true
}
