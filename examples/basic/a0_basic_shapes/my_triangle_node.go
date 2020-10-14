package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type myTriangleNode struct {
	nodes.Node

	background api.INode

	atlas   api.IAtlasX
	shapeID int

	color []float32
}

func newMyTriangleNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(myTriangleNode)

	o.Initialize(name)
	o.SetParent(parent)

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *myTriangleNode) build(world api.IWorld) error {
	b.Node.Build(world)

	b.atlas = world.GetAtlas(monoAtlasName)

	b.shapeID = b.atlas.GetShapeByName(triangleName)

	return nil
}

func (b *myTriangleNode) setColor(color api.IPalette) {
	b.color = color.Array()
}

func (b *myTriangleNode) setAlpha(alpha float32) {
	b.color[3] = alpha
}

// Draw renders shape
func (b *myTriangleNode) Draw(model api.IMatrix4) {
	b.atlas.Use()
	b.atlas.SetColor(b.color)
	b.atlas.Render(b.shapeID, model)
	b.atlas.UnUse()
}
