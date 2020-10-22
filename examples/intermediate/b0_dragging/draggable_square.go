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
}

func newDraggableSquare() *draggableSquare {
	o := new(draggableSquare)
	return o
}

func (d *draggableSquare) Build(world api.IWorld, parent api.INode) error {
	var err error

	d.drag = misc.NewDragState()
	d.localPosition = geometry.NewPoint()

	// ---------------------------------------------------------
	d.square, err = shapes.NewMonoSquareNode("Square", api.FILLOUTLINED, true, world, parent)
	if err != nil {
		return err
	}
	d.square.SetScale(100.0)
	d.square.SetPosition(90.0, 80.0)
	gsq := d.square.(*shapes.MonoSquareNode)
	gsq.SetFilledColor(color.NewPaletteInt64(color.LightPurple))
	gsq.SetOutlineColor(color.NewPaletteInt64(color.Black))

	d.dragDetect = newDraggable(true)

	return nil
}

func (d *draggableSquare) Position() api.IPoint {
	return d.square.Position()
}

func (d *draggableSquare) EventHandle(event api.IEvent) bool {
	if event.GetType() == api.IOTypeMouseMotion {
		mx, my := event.GetMousePosition()
		// fmt.Println("mx,my: ", mx, ", ", my)

		// Because the Layer and parent Scene have no transformation between
		// each other we could also pass "g" instead of "g.square".
		// Passing "g" would cause SetMotion...() to use g's parent which
		// is SplashScene verses rectangle node's parent which is GameLayer.
		// However, to be explicit I pass "g.square"
		d.drag.SetMotionStateUsing(mx, my, event.GetState(), d.square)

		gsq := d.square.(*shapes.MonoSquareNode)

		inside := d.dragDetect.PointInside()
		if inside {
			gsq.SetOutlineColor(color.NewPaletteInt64(color.White))
		} else {
			gsq.SetOutlineColor(color.NewPaletteInt64(color.Black))
		}

		d.dragDetect.Handle(d.square, event)

		if d.drag.IsDragging() && inside {
			pos := d.square.Position()
			x := pos.X() + d.drag.Delta().X()
			y := pos.Y() + d.drag.Delta().Y()
			d.square.SetPosition(x, y)
		}

	} else if event.GetType() == api.IOTypeMouseButtonDown || event.GetType() == api.IOTypeMouseButtonUp {
		mx, my := event.GetMousePosition()
		// On mouse events if state = 1 then dragging
		d.drag.SetButtonStateUsing(mx, my, event.GetButton(), event.GetState(), d.square)
	}

	return false
}
