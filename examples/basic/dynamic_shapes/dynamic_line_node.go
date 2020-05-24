package main

import (
	"unsafe"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicLineNode is a generic node
type DynamicLineNode struct {
	nodes.Node

	color []float32

	atlasName          string
	shape              api.IAtlasShape
	elementIndexOffset int
	backingArrayIdx    int

	p1Index int
	p2Index int

	// VBO update var
	vboOffset  int
	countBytes int
}

// NewDynamicLineNode constructs a generic shape node
func newDynamicLineNode(name string, world api.IWorld, parent api.INode) api.INode {
	o := new(DynamicLineNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)
	o.Build(world)
	return o
}

// Build configures the node
func (r *DynamicLineNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.Red).Array()

	r.shape = world.DynoAtlas().GetNextShape("Line")
	r.backingArrayIdx = r.shape.BackingArrayIdx()

	r.p1Index = 0
	r.p2Index = 1
	r.SetVBOOffset(0)
	// Each line is xyz + xyz = 6 components
	r.SetCountBytes(rendering.VboElementCountPerLine)
	r.SetElementOffset(0) // relative to EBO

	return nil
}

// SetColor sets line color
func (r *DynamicLineNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetPoint1 sets point 1 end position
func (r *DynamicLineNode) SetPoint1(x, y float32) {
	// r.shape.A
	r.shape.SetVertex2D(x, y, r.p1Index)
	r.SetDirty(true)
}

// SetPoint2 sets point 2 end position
func (r *DynamicLineNode) SetPoint2(x, y float32) {
	r.shape.SetVertex2D(x, y, r.p2Index)
	r.SetDirty(true)
}

// SetElementOffset sets EBOs offset requirement--the value is scaled
// to byte count using sizeof(int32)
func (r *DynamicLineNode) SetElementOffset(offset int) {
	r.elementIndexOffset = offset * int(unsafe.Sizeof(int32(0)))
	r.shape.SetElementOffset(r.elementIndexOffset)
}

// SetVBOOffset sets
func (r *DynamicLineNode) SetVBOOffset(offset int) {
	r.vboOffset = offset * int(unsafe.Sizeof(float32(0)))
}

// SetCountBytes sets
func (r *DynamicLineNode) SetCountBytes(count int) {
	r.countBytes = count * int(unsafe.Sizeof(float32(0)))
}

// Draw renders shape
func (r *DynamicLineNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.DynamicRenderGraphic)

	if r.IsDirty() {
		// Update buffer
		bufVertices := r.shape.Vertices(r.backingArrayIdx)

		renG.UpdatePreScaledUsing(r.vboOffset, r.countBytes, bufVertices)

		r.SetDirty(false)
	}

	renG.SetColor(r.color)

	renG.RenderElements(r.shape, r.shape.ElementCount(), r.elementIndexOffset, model)
}
