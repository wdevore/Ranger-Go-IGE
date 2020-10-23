package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// QTreeNode is a node for rendering quadtrees
type QTreeNode struct {
	nodes.Node

	color                  []float32
	minX, minY, maxX, maxY float32

	shape  api.IAtlasShape
	filled bool

	model api.IMatrix4

	tree api.IQuadTree
}

// NewQTreeNode constructs a node
func NewQTreeNode(minX, minY, maxX, maxY float32, name string, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(QTreeNode)
	o.Initialize(name)
	o.SetParent(parent)

	o.filled = filled
	o.minX = minX
	o.minY = minY
	o.maxX = maxX
	o.maxY = maxY

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

	if r.filled {
		name := "FilledRectangle"
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := "OutlineRectangle"
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate(r.minX, r.minY, r.maxX, r.maxY)

	return nil
}

func (r *QTreeNode) populate(minX, minY, maxX, maxY float32) {
	var vertices []float32

	vertices = []float32{
		minX, maxY, 0.0,
		minX, minY, 0.0,
		maxX, minY, 0.0,
		maxX, maxY, 0.0,
	}

	r.shape.SetVertices(vertices)

	var indices []uint32

	if r.filled {
		indices = []uint32{
			0, 1, 2,
			0, 2, 3,
		}
	} else {
		indices = []uint32{
			0, 1, 2, 3,
		}
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
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
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)

	renG.SetColor(r.color)

	r.model.Set(model)

	// Traverse iterates the entire tree.
	r.tree.Traverse(func(bounds api.IRectangle) {
		r.model.SetTranslate3Comp(bounds.Left(), bounds.Bottom(), 0.0)
		r.model.ScaleByComp(bounds.Width(), bounds.Height(), 1.0)
		renG.Render(r.shape, r.model)
	})
}
