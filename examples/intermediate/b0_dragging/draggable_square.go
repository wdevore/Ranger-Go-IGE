package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras/misc"
	"github.com/wdevore/Ranger-Go-IGE/extras/shapes"
)

type draggableSquare struct {
	square     api.INode
	dragDetect *draggable

	// Dragging
	drag          api.IDragging
	mx, my        int32
	localPosition api.IPoint
	pointInside   bool
	clickedInside bool
}

func newDraggableSquare() *draggableSquare {
	o := new(draggableSquare)
	return o
}

func (d *draggableSquare) Build(scale float32, world api.IWorld, parent api.INode) error {
	var err error

	d.drag = misc.NewDragState()
	d.localPosition = geometry.NewPoint()

	// ---------------------------------------------------------
	d.square, err = shapes.NewMonoSquareNode("Square", api.FILLOUTLINED, true, world, parent)
	if err != nil {
		return err
	}
	d.square.SetScale(scale)
	d.square.SetPosition(90.0, 80.0)
	gsq := d.square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightPurple))
	gsq.SetOutlineColor(color.NewPaletteInt64(color.Black))
	gsq.SetFilledAlpha(0.5)

	d.dragDetect = newDraggable(true)

	return nil
}

func (d *draggableSquare) Position() api.IPoint {
	return d.square.Position()
}

func (d *draggableSquare) BaseNode() api.INode {
	return d.square
}

func (d *draggableSquare) EventHandle(event api.IEvent) bool {
	mx, my := event.GetMousePosition()

	if event.GetType() == api.IOTypeMouseMotion {
		gsq := d.square.(*shapes.MonoSquareNode)

		inside := d.dragDetect.Handle(mx, my, d.square, event)

		if inside {
			gsq.SetOutlineColor(color.NewPaletteInt64(color.White))
		} else {
			gsq.SetOutlineColor(color.NewPaletteInt64(color.Black))
		}
	}

	if event.GetType() == api.IOTypeMouseButtonDown {
		d.clickedInside = d.dragDetect.Handle(mx, my, d.square, event)
		d.drag.SetButtonStateUsing(mx, my, event.GetButton(), event.GetState(), d.square)
	}

	if d.clickedInside {
		// Because the Layer and parent Scene have no transformation between
		// each other we could also pass "d" instead of "d.square".
		// Passing "d" would cause SetMotion...() to use d's parent which
		// is SplashScene verses rectangle node's parent which is GameLayer.
		// However, to be explicit I pass "d.square"
		d.drag.SetMotionStateUsing(mx, my, event.GetState(), d.square)

		if d.drag.IsDragging() {
			pos := d.square.Position()
			x := pos.X() + d.drag.Delta().X()
			y := pos.Y() + d.drag.Delta().Y()
			d.square.SetPosition(x, y)
		}
	}

	if event.GetType() == api.IOTypeMouseButtonUp {
		// On mouse events if state = 1 then dragging
		d.drag.SetButtonStateUsing(mx, my, event.GetButton(), event.GetState(), d.square)
	}

	return false
}
