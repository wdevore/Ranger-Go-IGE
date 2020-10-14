package main

import (
	"errors"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
)

type mySquareNode struct {
	nodes.Node

	background api.INode

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

	atlas := world.GetAtlas(api.MonoAtlasName)

	if atlas == nil {
		return errors.New("Expected to find StaticMono Atlas")
	}

	b.SetAtlas(atlas)

	b.shapeID = atlas.GetShapeByName(api.CenteredFilledSquareShapeName)
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

// Draw renders shape
func (b *mySquareNode) Draw(model api.IMatrix4) {
	// Note: We don't need to call the Atlas's Use() method
	// because the node.Visit() will do that for us.
	atlas := b.Atlas()

	atlas.SetColor(b.color)
	atlas.Render(b.shapeID, model)
}
