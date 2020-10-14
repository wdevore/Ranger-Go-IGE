package main

import (
	"errors"

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

	b.atlas = world.GetAtlas(api.MonoAtlasName)

	b.shapeID = b.atlas.GetShapeByName(api.FilledTriangleShapeName)

	if b.shapeID < 0 {
		return errors.New("myTriangleNode: Could not find triangle shape")
	}

	return nil
}

func (b *myTriangleNode) setColor(color api.IPalette) {
	b.color = color.Array()
}

func (b *myTriangleNode) setAlpha(alpha float32) {
	b.color[3] = alpha
}

func (b *myTriangleNode) Atlas() api.IAtlasX {
	return b.atlas
}

// Draw renders shape
func (b *myTriangleNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	b.atlas.SetColor(b.color)
	b.atlas.Render(b.shapeID, model)
}
