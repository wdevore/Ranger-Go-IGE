package main

import (
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/misc"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes/custom"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// ZoneCircle is a basic vector circle shape.
type ZoneCircle struct {
	innerCircle api.INode
	outerCircle api.INode

	innerColor   api.IPalette
	outerColor   api.IPalette
	enteredColor api.IPalette

	zone      api.IZone // CircleZone
	zoneState int
	zoneID    int

	tweenEnabled      bool
	tweenZoomIn       *gween.Tween
	tweenZoomOut      *gween.Tween
	isFinished        bool
	tweenCurrentValue float32

	zoomTo   float64
	zoomFrom float64
	duration float64

	// Typically a zone manager would be the subscriber
	subscribers []api.IZoneListener

	zoneMan *zoneManager
}

// NewZoneCircle constructs a circle zone
func NewZoneCircle(name string, id int, zoneMan *zoneManager) *ZoneCircle {
	o := new(ZoneCircle)
	o.zoneMan = zoneMan
	o.zoneID = id
	return o
}

// Build configures the node
func (z *ZoneCircle) Build(innerRadius, outerRadius float32, position api.IPoint, world api.IWorld, parent api.INode) {
	z.subscribers = []api.IZoneListener{}

	z.innerColor = color.NewPaletteInt64(color.PanSkin)
	z.outerColor = color.NewPaletteInt64(color.Silver)
	z.enteredColor = color.NewPaletteInt64(color.LightPurple)

	var err error
	z.innerCircle, err = custom.NewStaticCircleNode("InnerCircle", false, world, parent)
	if err != nil {
		panic(err)
	}
	z.innerCircle.SetVisible(false)
	gol2 := z.innerCircle.(*custom.StaticCircleNode)
	gol2.SetColor(z.innerColor)

	z.outerCircle, err = custom.NewStaticCircleNode("OuterCircle", false, world, parent)
	if err != nil {
		panic(err)
	}
	z.outerCircle.SetVisible(false)
	gol2 = z.outerCircle.(*custom.StaticCircleNode)
	gol2.SetColor(z.outerColor)

	z.zone = misc.NewCircleZone()

	z.SetRadi(innerRadius, outerRadius)
	z.SetPosition(position.X(), position.Y())

	z.isFinished = true
}

// ID returns the zone's unique id
func (z *ZoneCircle) ID() int {
	return z.zoneID
}

// RequestNotification asks for notification when an event on the zone
// happens. ZoneManager is a subscriber.
func (z *ZoneCircle) RequestNotification(listener api.IZoneListener) {
	z.subscribers = append(z.subscribers, listener)
}

// SetTweenRange sets the from and to values
func (z *ZoneCircle) SetTweenRange(from, to float64) {
	z.zoomTo = to
	z.zoomFrom = from
}

// SetTweenDuration sets the animation time duration
func (z *ZoneCircle) SetTweenDuration(duration float64) {
	z.duration = duration
}

// SetPosition sets position of zone
func (z *ZoneCircle) SetPosition(x, y float32) {
	z.innerCircle.SetPosition(x, y)
	z.outerCircle.SetPosition(x, y)

	zo := z.zone.(*misc.CircleZone)
	zo.SetPosition(x, y)
}

// Position returns the zone's center
func (z *ZoneCircle) Position() api.IPoint {
	zo := z.zone.(*misc.CircleZone)
	return zo.Position()
}

// SetRadi sets circle's inner and outer radi
func (z *ZoneCircle) SetRadi(innerRadius, outerRadius float32) {
	zo := z.zone.(*misc.CircleZone)
	zo.SetRadi(float64(innerRadius/2), float64(outerRadius/2))

	cr := z.innerCircle.(*custom.StaticCircleNode)
	cr.SetScale(innerRadius)

	cr = z.outerCircle.(*custom.StaticCircleNode)
	cr.SetScale(outerRadius)
}

// SetSegments sets how many segments on the circle (default = 12)
func (z *ZoneCircle) SetSegments(segments int) {
	// z.segments = segments
}

// SetInnerColor sets circle's inner color (default = LightGray)
func (z *ZoneCircle) SetInnerColor(color api.IPalette) {
	z.innerColor = color
}

// SetOuterColor sets circle's outer color (default = Silver)
func (z *ZoneCircle) SetOuterColor(color api.IPalette) {
	z.outerColor = color
}

// UpdateCheck forces the zone to update based on a given point
func (z *ZoneCircle) UpdateCheck(point api.IPoint) (state, id int) {
	newState, stateChanged := z.zone.Update(point)
	id = z.zoneID
	if stateChanged {
		z.zoneState = newState

		// Send message to listeners. The "id" is a self identifier.
		// Most likely the ZoneManager
		for _, listener := range z.subscribers {
			listener.Notify(z.zoneState, id)
		}

		z.createTween(z.zoneState, id)
	}

	return newState, id
}

// TweenUpdate updates any tweens if enabled
func (z *ZoneCircle) TweenUpdate(msPerUpdate float64) (float64, bool) {
	if z.tweenEnabled {
		switch z.zoneState {
		case api.CrossStateEntered:
			z.tweenCurrentValue, z.isFinished = z.tweenZoomIn.Update(float32(msPerUpdate))
			if z.isFinished {
				z.tweenEnabled = false
				z.zoneMan.SetAnimationActive(false)
			}
		case api.CrossStateExited:
			z.tweenCurrentValue, z.isFinished = z.tweenZoomOut.Update(float32(msPerUpdate))
			if z.isFinished {
				z.tweenEnabled = false
				z.zoneMan.SetAnimationActive(false)
			}
		}
	}

	return float64(z.tweenCurrentValue), z.isFinished
}

// ----------------------------------------------------------
// IZoneListener implementation
// ----------------------------------------------------------

func (z *ZoneCircle) createTween(state, id int) {
	// Whenever an tween needs to be created we are animating from a
	// "begin" value to an "end" value--regardless of zooming In or Out.

	// The "to" should also we the final zoom value.
	// The "from" should be whatever the initial "to" value is OR
	//   the current value of the ZoomNode if an animation was in progress
	//   during the Enter event (aka tweenEnabled).

	// fmt.Println("Scale: ", z.zoneMan.ZoomScale(), ", current: ", z.tweenCurrentValue, ", from: ", z.zoomFrom, ", to: ", z.zoomTo)
	switch state {
	case api.CrossStateEntered:
		if !z.isFinished {
			z.tweenZoomIn = gween.New(float32(z.zoneMan.ZoomScale()), float32(z.zoomTo), float32(z.duration), ease.InOutQuad)
		} else {
			if z.zoneMan.AnimationActive() {
				z.tweenZoomIn = gween.New(float32(z.zoneMan.ZoomScale()), float32(z.zoomTo), float32(z.duration), ease.InOutQuad)
			} else {
				z.tweenZoomIn = gween.New(float32(z.zoomFrom), float32(z.zoomTo), float32(z.duration), ease.InOutQuad)
			}
		}
		z.innerColor = color.NewPaletteInt64(color.Lime)
		z.tweenEnabled = true
		z.zoneMan.SetAnimationActive(true)
	case api.CrossStateExited:
		if !z.isFinished {
			z.tweenZoomOut = gween.New(float32(z.zoneMan.ZoomScale()), float32(z.zoomFrom), float32(z.duration), ease.InOutQuad)
		} else {
			if z.zoneMan.AnimationActive() {
				z.tweenZoomOut = gween.New(float32(z.zoneMan.ZoomScale()), float32(z.zoomFrom), float32(z.duration), ease.InOutQuad)
			} else {
				z.tweenZoomOut = gween.New(float32(z.zoomTo), float32(z.zoomFrom), float32(z.duration), ease.InOutQuad)
			}
		}
		z.innerColor = color.NewPaletteInt64(color.LightGray)
		z.tweenEnabled = true
		z.zoneMan.SetAnimationActive(true)
	}
}
