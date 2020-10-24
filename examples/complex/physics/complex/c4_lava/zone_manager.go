package main

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

// zoneManager handles zones
// The ZM coordinates between zones and any animations created by them.
// When a zone is entered all other zones' animations must stop
type zoneManager struct {
	parent api.INode

	zones         []*ZoneCircle
	enteredZoneID int

	// Zooming
	zoom api.INode

	zoomScale float64

	animationActive bool
}

// newZoneManager creates a zone manager
func newZoneManager(parent api.INode) *zoneManager {
	o := new(zoneManager)
	o.parent = parent
	return o
}

func (z *zoneManager) Build(world api.IWorld) {
	var err error

	z.zoom, err = extras.NewZoomNode("ZoomNode", world, z.parent)
	if err != nil {
		panic(err)
	}
	gz := z.zoom.(*extras.ZoomNode)
	gz.SetStepSize(0.05)

	zone := NewZoneCircle("RightCircleZone", objectRightZone, z)
	// TODO the zone positions should be able to target anchors/markers
	// defined ON the land
	zone.Build(20.0, 22.0, z.parent.Position(), z.parent.World(), z.zoom)
	z.zones = append(z.zones, zone)
	zone.SetTweenRange(1.0, 2.0)
	zone.SetTweenDuration(1000.0)
	zone.RequestNotification(z)
	zone.SetPosition(30.0, -10.0)

	zone = NewZoneCircle("LeftCircleZone", objectLeftZone, z)
	zone.Build(13.0, 15.0, z.parent.Position(), z.parent.World(), z.zoom)
	z.zones = append(z.zones, zone)
	zone.SetTweenRange(1.0, 2.0)
	zone.SetTweenDuration(1000.0)
	zone.RequestNotification(z)
	zone.SetPosition(-30.0, -10.0)
	// gr.SetPosition(0.0, 15.0)
}

// GetZoom returns zoom INode
func (z *zoneManager) GetZoom() api.INode {
	return z.zoom
}

// UpdateCheck updates zone tweens
func (z *zoneManager) UpdateCheck(point api.IPoint, msPerUpdate float64) {
	isFinished := true

	for _, zone := range z.zones {
		zone.UpdateCheck(point)

		// Animate only the zone that was enter. The other zone's
		// animation is frozen/stopped.
		if z.enteredZoneID == zone.ID() {
			z.zoomScale, isFinished = zone.TweenUpdate(msPerUpdate)
			if !isFinished {
				gz := z.zoom.(*extras.ZoomNode)
				gz.ScaleTo(float32(z.zoomScale))
			}
		}
	}
}

func (z *zoneManager) AnimationActive() bool {
	return z.animationActive
}

func (z *zoneManager) SetAnimationActive(active bool) {
	z.animationActive = active
}

func (z *zoneManager) ZoomScale() float32 {
	gz := z.zoom.(*extras.ZoomNode)

	return gz.ZoomScale()
}

// ----------------------------------------------------------
// IZoneListener implementation
// ----------------------------------------------------------

// Notify receives messages from IZone objects
func (z *zoneManager) Notify(state, id int) {
	if state != api.CrossStateEntered {
		return
	}

	z.enteredZoneID = id

	// fmt.Println("ZM notified: ", z.enteredZoneID)

	// Find zone that matches "id"
	for _, zone := range z.zones {
		if z.enteredZoneID == zone.ID() {
			gz := z.zoom.(*extras.ZoomNode)
			gz.SetFocalPoint(zone.Position().X(), zone.Position().Y())
			break
		}
	}
}
