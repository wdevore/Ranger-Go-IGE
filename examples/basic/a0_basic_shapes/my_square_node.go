package main

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type mySquareNode struct {
	nodes.Node

	background api.INode

	atlas   api.IAtlasX
	shapeID int

	color []float32
}

func newMySquareNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(mySquareNode)

	o.Initialize(name)
	o.SetParent(parent)

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

func (b *mySquareNode) build(world api.IWorld) error {
	b.Node.Build(world)

	b.atlas = world.GetAtlas(api.MonoAtlasName)

	b.shapeID = b.atlas.GetShapeByName(api.CenteredFilledSquareShapeName)
	if b.shapeID < 0 {
		return errors.New("mySquareNode: Could not find square shape")
	}

	return nil
}

func (b *mySquareNode) setColor(color api.IPalette) {
	b.color = color.Array()
}

func (b *mySquareNode) setAlpha(alpha float32) {
	b.color[3] = alpha
}

func (b *mySquareNode) Atlas() api.IAtlasX {
	return b.atlas
}

// Draw renders shape
func (b *mySquareNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	b.atlas.SetColor(b.color)
	b.atlas.Render(b.shapeID, model)
}
