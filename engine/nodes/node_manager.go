package nodes

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// The node manager is basically the SceneGraph
type nodeManager struct {
	world api.IWorld

	clearBackground bool

	// Stack of node
	nodStack *nodeStack

	transStack api.ITransformStack

	timingTargets api.INodeList
	eventTargets  api.INodeList

	projection *display.Projection
	viewport   *display.Viewport

	viewSpace    api.IAffineTransform
	invViewSpace api.IAffineTransform
}

// NewNodeManager constructs a manager for node.
// It manages the lifecycle and events
func NewNodeManager(world api.IWorld) api.INodeManager {
	o := new(nodeManager)
	o.world = world

	// It is very rare that the manager would clear the background
	// because almost all nodes will handle clearing/painting their
	// own backgrounds.
	o.clearBackground = false

	o.nodStack = newNodeStack()
	o.transStack = newTransformStack()

	o.viewSpace = maths.NewTransform()
	o.invViewSpace = maths.NewTransform()

	o.timingTargets = NewNodeList()
	o.eventTargets = NewNodeList()

	return o
}

func (n *nodeManager) Configure() {
	// Setup view/projection matrix composition

	n.configureProjections(n.world)

	// Setup initial VP portion of MVP
	m4 := maths.NewMatrix4()
	m4.SetFromAffine(n.viewSpace)
	// (i.e. m4 = projection * view)
	m4.PostMultiply(n.projection.Matrix())

	// Initialize will make m4 the current matrix ready to be
	// placed on the stack on the first save() call.
	n.transStack.Initialize(m4)
}

func (n *nodeManager) ClearEnabled(clear bool) {
	n.clearBackground = clear
}

func (n *nodeManager) PreVisit() {
}

func (n *nodeManager) Visit(interpolation float64) bool {
	if n.nodStack.isEmpty() {
		return false
	}

	// fmt.Println("NodeManager: visiting ", m.stack.runningNode)

	if n.nodStack.hasNextNode() {
		n.setNextNode()
	}

	n.transStack.Save()

	runningScene := n.nodStack.runningNode.(api.IScene)

	action := runningScene.TransitionAction()

	if action == api.SceneReplaceTake {
		repl := runningScene.GetReplacement()
		// fmt.Println("NodeManager: SceneReplaceTake with ", repl)
		if repl != nil {
			n.nodStack.replace(repl)
			// Immediately switch to the new replacement node
			if n.nodStack.hasNextNode() {
				n.setNextNode()
			}
		} else {
			n.exitNodes(n.nodStack.runningNode)
			n.nodStack.pop()
		}
	}

	// Visit the running node
	Visit(n.nodStack.runningNode, n.transStack, interpolation)

	n.transStack.Restore()

	return true // continue to draw.
}

func (n *nodeManager) PostVisit() {
}

func (n *nodeManager) PopNode() api.INode {
	return n.nodStack.pop()
}

func (n *nodeManager) PushNode(node api.INode) {
	n.nodStack.nextNode = node
	n.nodStack.push(node)
}

func (n *nodeManager) ReplaceNode(node api.INode) {
	n.nodStack.replace(node)
}

// --------------------------------------------------------------------------
// Timing
// --------------------------------------------------------------------------

func (n *nodeManager) Update(msPerUpdate, secPerUpdate float64) {
	for _, target := range n.timingTargets.Items() {
		target.Update(msPerUpdate, secPerUpdate)
	}
}

func (n *nodeManager) RegisterTarget(target api.INode) {
	n.timingTargets.Add(target)
}

func (n *nodeManager) UnRegisterTarget(target api.INode) {
	n.timingTargets.Remove(target)
}

// --------------------------------------------------------------------------
// IO events
// --------------------------------------------------------------------------

func (n *nodeManager) RegisterEventTarget(target api.INode) {
	n.eventTargets.Add(target)
}

func (n *nodeManager) UnRegisterEventTarget(target api.INode) {
	n.eventTargets.Remove(target)
}

func (n *nodeManager) RouteEvents(event api.IEvent) {
	if n.eventTargets == nil {
		return
	}

	for _, target := range n.eventTargets.Items() {
		handled := target.Handle(event)

		if handled {
			break
		}
	}
}

func (n *nodeManager) setNextNode() {
	if n.nodStack.hasRunningNode() {
		n.exitNodes(n.nodStack.runningNode)
	}

	n.nodStack.runningNode = n.nodStack.nextNode
	n.nodStack.clearNextNode()

	// fmt.Println("NodeManager: new running node ", m.stack.runningNode)

	n.enterNodes(n.nodStack.runningNode)
}

// End cleans up NodeManager by clearing the stack and calling all Exits
func (n *nodeManager) End() {
	// Dump the stack

	pn := n.PopNode()

	for pn != nil {
		n.exitNodes(pn)
		pn = n.PopNode()
	}

	n.eventTargets = nil
}

// -----------------------------------------------------
// Scene lifecycles
// -----------------------------------------------------

func (n *nodeManager) enterNodes(node api.INode) {
	// fmt.Println("NodeManager: enter-node ", node)
	node.EnterNode(n)

	children := node.Children()
	for _, child := range children {
		n.enterNodes(child)
	}
}

func (n *nodeManager) exitNodes(node api.INode) {
	// fmt.Println("NodeManager: exit-node ", node)
	node.ExitNode(n)

	children := node.Children()
	for _, child := range children {
		n.exitNodes(child)
	}
}

func (n *nodeManager) Debug() {
}

func (n nodeManager) String() string {
	return fmt.Sprintf("%s", n.nodStack)
}

// DeleteAt removes an item from the slice
func DeleteAt(i int, slice []api.INode) {
	// Remove the element at index i from slice.
	copy(slice[i:], slice[i+1:]) // Shift a[i+1:] left one index.
	slice[len(slice)-1] = nil    // Erase last element (write zero value).
	slice = slice[:len(slice)-1] // Truncate slice.
}

// FindFirstElement finds the first item in the slice
func FindFirstElement(node api.INode, slice []api.INode) int {
	for idx, item := range slice {
		if item.ID() == node.ID() {
			return idx
		}
	}

	return -1
}

func (n *nodeManager) configureProjections(world api.IWorld) {
	wp := world.Properties().Window

	n.viewport = display.NewViewport()

	n.viewport.SetDimensions(0, 0, wp.DeviceRes.Width, wp.DeviceRes.Height)
	n.viewport.Apply()

	// Calc the aspect ratio between the physical (aka device) dimensions and the
	// the virtual (aka user's design choice) dimensions.

	deviceRatio := float64(wp.DeviceRes.Width) / float64(wp.DeviceRes.Height)
	virtualRatio := float64(wp.VirtualRes.Width) / float64(wp.VirtualRes.Height)

	xRatioCorrection := float64(wp.DeviceRes.Width) / float64(wp.VirtualRes.Width)
	yRatioCorrection := float64(wp.DeviceRes.Height) / float64(wp.VirtualRes.Height)

	var ratioCorrection float64

	if virtualRatio < deviceRatio {
		ratioCorrection = yRatioCorrection
	} else {
		ratioCorrection = xRatioCorrection
	}

	n.projection = display.NewCamera()

	if world.Properties().Camera.Centered {
		n.projection.SetCenteredProjection()
	} else {
		n.projection.SetProjection(
			float32(ratioCorrection),
			0.0, 0.0,
			float32(wp.DeviceRes.Height), float32(wp.DeviceRes.Width))
	}

	center := maths.NewTransform()
	center.Scale(float32(xRatioCorrection), float32(yRatioCorrection))

	n.viewSpace.SetByTransform(center)

	n.invViewSpace.SetByTransform(center)
	n.invViewSpace.Invert()
}
