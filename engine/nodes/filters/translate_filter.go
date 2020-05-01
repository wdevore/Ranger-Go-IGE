package filters

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

// TranslateFilter will exclude or "block" rotations and scales from
// propagating to the children (default), but "passes" translations.
type TranslateFilter struct {
	nodes.Node
	world api.IWorld

	// This makes the node an IFilter type
	Filter
}

// NewTranslateFilter constructs a Translate filter.
func NewTranslateFilter(name string, parent api.INode) api.INode {
	o := new(TranslateFilter)
	o.Initialize(name)
	o.SetParent(parent)
	o.initializeFilter()
	return o
}

// Build configures the node
func (t *TranslateFilter) Build(world api.IWorld) {
	t.Node.Build(world)
}

// Visit is special in that it they provide their own implementation.
// Because this is a Translate filter we "filter out" everything
// but the translation component from the immediate parent.
func (t *TranslateFilter) Visit(transStack api.ITransformStack, interpolation float64) {
	if !t.IsVisible() {
		return
	}

	transStack.Save()

	children := t.Children()

	for _, child := range children {
		transStack.Save()

		if t.HasParent() {
			parent := t.Parent()

			// This removes the immediate parent's transform effects
			transStack.Apply(parent.InverseTransform())

			// Re-introduce only the parent's translation component by
			// excluding Rotation and Scale
			parent.CalcFilteredTransform(false, true, true, t.components)

			// And update context to reflect the exclusion.
			transStack.Apply(t.components)
		} else {
			fmt.Println("TranslateFilter: node ", t, " has NO parent")
			return
		}

		// Now visit the child with the modified context
		nodes.Visit(child, transStack, interpolation)

		transStack.Restore()
	}

	transStack.Restore()
}
