package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// QTreeNode is a node for rendering quadtrees
type QTreeNode struct {
	nodes.Node

	color []float32

	shapeID int

	model api.IMatrix4

	tree api.IQuadTree
}

// NewQTreeNode constructs a node
func NewQTreeNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(QTreeNode)

	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *QTreeNode) build(world api.IWorld) error {
	r.Node.Build(world)

	r.model = maths.NewMatrix4()

	r.color = color.NewPaletteInt64(color.White).Array()

	atl := world.GetAtlas(api.MonoAtlasName)
	r.SetAtlas(atl)

	atlas := atl.(api.IStaticAtlasX)

	r.shapeID = atlas.GetShapeByName(api.UnCenteredOutlinedSquareShapeName + "_Q")

	if r.shapeID < 0 {
		// Add shape
		vertices, indices, mode := generators.GenerateUnitRectangleVectorShape(false, false)
		r.shapeID = atlas.AddShape(api.UnCenteredOutlinedSquareShapeName, vertices, indices, mode)
	}

	return nil
}

// SetTree assigns a quadtree to this visual node.
func (r *QTreeNode) SetTree(tree api.IQuadTree) {
	r.tree = tree
}

// SetColor sets square color
func (r *QTreeNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *QTreeNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *QTreeNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := r.Atlas()

	atlas.SetColor(r.color)

	r.model.Set(model)

	// Traverse iterates the entire tree.
	r.tree.Traverse(func(bounds api.IRectangle) {
		r.model.SetTranslate3Comp(bounds.Left(), bounds.Bottom(), 0.0)
		r.model.ScaleByComp(bounds.Width(), bounds.Height(), 1.0)
		atlas.Render(r.shapeID, r.model)
	})
}
