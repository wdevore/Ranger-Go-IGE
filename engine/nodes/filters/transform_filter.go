package filters

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// TransformFilter by default passes Rotation and Translation components,
// which is the most common use case.
type TransformFilter struct {
	nodes.Node

	// This makes the node an IFilter type.
	// However, you MUST make sure your filter implements all methods
	// correctly. node.Visit() makes the assumption that if the node
	// can't be cast to an IFilter that it must be an INode.
	Filter
}

// NewTransformFilter constructs a default transform filter. Default
// is to inherit both Rotation and Translation.
func NewTransformFilter(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(TransformFilter)
	o.Initialize(name)
	o.SetParent(parent)
	o.initializeFilter()
	o.InheritRotationAndTranslation()
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (t *TransformFilter) Build(world api.IWorld) error {
	return t.Node.Build(world)
}

// Visit is special in that it they provide their own implementation
func (t *TransformFilter) Visit(transStack api.ITransformStack, interpolation float64) {
	if !t.IsVisible() {
		return
	}

	transStack.Save()

	children := t.Children()

	for _, child := range children {
		transStack.Save()

		if t.HasParent() {
			parent := t.Parent()

			// Re-introduce only the parent's components as defined by exclusion flags.
			parent.CalcFilteredTransform(t.excludeTranslation,
				t.excludeRotation,
				t.excludeScale, t.components)

			// Combine using pre-multiply
			// "parent.InverseTransform" removes the immediate parent's transform effects
			maths.Multiply(t.components, parent.InverseTransform(), t.AffineTransform())

			// Merge them with the current context.
			transStack.ApplyAffine(t.AffineTransform())
		} else {
			fmt.Println("TransformFilter: node ", t, " has NO parent")
			return
		}

		// Now visit the child with the modified context
		nodes.Visit(child, transStack, interpolation)

		transStack.Restore()
	}

	transStack.Restore()
}
