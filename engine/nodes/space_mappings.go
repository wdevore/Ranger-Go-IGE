package nodes

import (
	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

var tViewPoint = geometry.NewPoint()
var comp = maths.NewTransform()
var out = maths.NewTransform()

// MapDeviceToView maps mouse-space device coordinates to view-space
func MapDeviceToView(world api.IWorld, dvx, dvy int32, viewPoint api.IPoint) {
	viewPoint.SetByComp(float32(dvx), float32(dvy))
	viewPoint.MulPoint(world.InvertedViewspace())
}

// MapDeviceToNode maps mouse-space device coordinates to local node-space
// Optional scaling before returning from function. Extremely rare to use
// localPoint.SetByComp(localPoint.X() * node.Scale(), localPoint.Y() * node.Scale())
func MapDeviceToNode(dvx, dvy int32, node api.INode, localPoint api.IPoint) {
	// Mapping from device to node requires transforms from two "directions"
	// 1st is upwards transform and the 2nd is downwards transform.

	// downwards from device-space to view-space
	MapDeviceToView(node.World(), dvx, dvy, tViewPoint)

	// OpenGL's +Y axis is upwards so we either flip the Y axis here
	// or flip the mouse's +Y axis in cursorPosCallback(...)
	// dvr := node.World().Properties().Window.DeviceRes
	// MapDeviceToView(node.World(), dvx, int32(dvr.Height)-dvy, tViewPoint)

	// Upwards from node to world-space (aka view-space)
	wtn := WorldToNodeTransform(node, nil)

	// Now map view-space point to local-space of node
	wtn.TransformCompToPoint(tViewPoint.X(), tViewPoint.Y(), localPoint)
}

// MapNodeToNode maps node's local origin (0,0) to another node's space
// Supplying a psuedo-root can reduce multiplications, thus use it if possible.
func MapNodeToNode(sourceNode, destinationNode api.INode, nodePoint api.IPoint, psuedoRoot api.INode) {
	ntw := NodeToWorldTransform(sourceNode, psuedoRoot)
	ntw.TransformCompToPoint(0.0, 0.0, nodePoint)
	// nodePoint is now in world-space

	wtn := WorldToNodeTransform(destinationNode, psuedoRoot) // nodePoint
	wtn.TransformCompToPoint(nodePoint.X(), nodePoint.Y(), nodePoint)
	// nodePoint is now in the destination node's space

	// OR via view-space (not recommended as it is extra operations and rounding)
	// MapNodeToDevice(sourceNode.World(), sourceNode, nodePoint)
	// MapDeviceToNode(int32(nodePoint.X()), int32(nodePoint.Y()), destinationNode, nodePoint)
}

// MapNodeToWorld maps node's local origin to world-space
func MapNodeToWorld(node api.INode, point api.IPoint) {
	ntw := NodeToWorldTransform(node, nil)
	ntw.TransformCompToPoint(0.0, 0.0, point)
}

// MapWorldToNode maps a world-space coord to a specific node
func MapWorldToNode(node api.INode, worldPoint, nodePoint api.IPoint) {
	wtn := WorldToNodeTransform(node, nil) // nodePoint
	wtn.TransformCompToPoint(worldPoint.X(), worldPoint.Y(), nodePoint)
}

// MapNodeToDevice maps node local origin to device-space (aka mouse-space)
func MapNodeToDevice(world api.IWorld, node api.INode, devicePoint api.IPoint) {
	ntw := NodeToWorldTransform(node, nil)
	ntw.TransformCompToPoint(0.0, 0.0, devicePoint)
	devicePoint.SetByComp(devicePoint.X(), devicePoint.Y())
	devicePoint.MulPoint(world.Viewspace())

	// world.ViewSpace().TransformCompToPoint(viewPoint.X(), viewPoint.Y(), viewPoint)
}

// WorldToNodeTransform maps a world-space coordinate to local-space of node
func WorldToNodeTransform(node api.INode, psuedoRoot api.INode) api.IAffineTransform {
	wtn := NodeToWorldTransform(node, psuedoRoot)
	wtn.Invert()
	return wtn
}

// NodeToWorldTransform maps a local-space coordinate to world-space
func NodeToWorldTransform(node api.INode, psuedoRoot api.INode) api.IAffineTransform {
	// A transform to accumulate the parent transforms.
	comp.SetByTransform(node.CalcTransform())

	// Iterate "upwards" starting with the child towards the parents
	// starting with this child's parent.
	p := node.Parent()

	for p != nil {
		parentT := p.CalcTransform()

		// Because we are iterating upwards we need to pre-multiply each child.
		// Ex: [child] x [parent]
		// ----------------------------------------------------------
		//           [comp] x [parentT]
		//               |
		//               | out
		//               v
		//             [comp] x [parentT]
		//                 |
		//                 | out
		//                 v
		//               [comp] x [parentT...]
		//
		// This is a pre-multiply order
		// [child] x [parent of child] x [parent of parent of child]...
		//
		// In other words the child is mutiplied "into" the parent.
		maths.Multiply(comp, parentT, out)
		comp.SetByTransform(out)

		if p == psuedoRoot {
			// fmt.Println("SpaceMappings: hit psuedoRoot")
			break
		}

		// next parent upwards
		p = p.Parent()
	}

	return comp
}
