package extras

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicTextureNode is a dynamic texture ready node.
type DynamicTextureNode struct {
	nodes.Node

	textureRenderer api.ITextureRenderer
	textureAtlas    api.ITextureAtlas

	index int

	text  string
	color []float32
}

// NewDynamicTextureNode constructs a generic shape node
func NewDynamicTextureNode(name string, textureAtlas api.ITextureAtlas, textureRenderer api.ITextureRenderer, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicTextureNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	o.textureRenderer = textureRenderer
	o.textureAtlas = textureAtlas

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (d *DynamicTextureNode) build(world api.IWorld) error {
	d.Node.Build(world)

	d.color = color.NewPaletteInt64(color.White).Array()

	return nil
}

// Draw renders shape
func (d *DynamicTextureNode) Draw(model api.IMatrix4) {
	d.textureRenderer.Use()
	d.textureRenderer.SetColor(d.color)

	d.textureRenderer.SelectCoordsByIndex(d.index)

	d.textureRenderer.Draw(model)
}

// SetIndex sets the active sub texture index
func (d *DynamicTextureNode) SetIndex(index int) {
	d.index = index
}

// SetColor sets mixin color
func (d *DynamicTextureNode) SetColor(color []float32) {
	d.color = color
}
