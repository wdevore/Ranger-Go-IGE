package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/misc"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

type draggableSquare struct {
	square    api.INode
	outSquare api.INode

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
	d.square, err = custom.NewStaticSquareNode("FilledSqr", true, true, world, parent)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	d.square.SetScale(100.0)
	d.square.SetPosition(90.0, 80.0)
	gsq := d.square.(*custom.StaticSquareNode)
	gsq.SetColor(color.NewPaletteInt64(color.LightPurple))

	d.outSquare, err = newCustomRectangleNode("OutlineSqr", true, false, world, parent)
	if err != nil {
		return err
	}
	d.outSquare.SetScale(100.0)
	d.outSquare.SetPosition(90.0, 80.0)
	gosq := d.outSquare.(*customRectangleNode)
	gosq.SetColor(color.NewPaletteInt64(color.White))

	return nil
}

func (d *draggableSquare) PointInside(p api.IPoint) bool {
	gsq := d.outSquare.(*customRectangleNode)
	inside := gsq.PointInside()
	return inside
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
		// This gets the local-space coords of the rectangle node.
		// Note: OpenGL's +Y axis is upward
		nodes.MapDeviceToNode(mx, my, d.square, d.localPosition)
		// fmt.Println("localPosition: ", d.localPosition)

		if d.drag.IsDragging() && d.PointInside(d.localPosition) {
			pos := d.square.Position()
			x := pos.X() + d.drag.Delta().X()
			y := pos.Y() + d.drag.Delta().Y()
			d.square.SetPosition(x, y)
			d.outSquare.SetPosition(x, y)
		}

	} else if event.GetType() == api.IOTypeMouseButtonDown || event.GetType() == api.IOTypeMouseButtonUp {
		mx, my := event.GetMousePosition()
		// On mouse events if state = 1 then dragging
		d.drag.SetButtonStateUsing(mx, my, event.GetButton(), event.GetState(), d.square)
	}

	return false
}
