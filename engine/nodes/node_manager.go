package nodes

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.5-core/gl"

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

	// Used during preVisit
	preNode api.INode

	projection *display.Projection
	viewport   *display.Viewport

	projLoc int32
	viewLoc int32
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

	o.timingTargets = NewNodeList()
	o.eventTargets = NewNodeList()

	return o
}

func (n *nodeManager) Configure() error {
	// Setup view/projection matrix composition

	n.configureProjections(n.world)

	programID := n.world.Shader().Program()

	n.projLoc = gl.GetUniformLocation(programID, gl.Str("projection\x00"))
	if n.projLoc < 0 {
		return errors.New("SplashScene: couldn't find 'projection' uniform variable")
	}

	n.viewLoc = gl.GetUniformLocation(programID, gl.Str("view\x00"))
	if n.viewLoc < 0 {
		return errors.New("SplashScene: couldn't find 'view' uniform variable")
	}

	pm := n.projection.Matrix().Matrix()
	gl.UniformMatrix4fv(n.projLoc, 1, false, &pm[0])

	// vm := n.viewSpace.Matrix()
	vm := n.world.Viewspace().Matrix()
	gl.UniformMatrix4fv(n.viewLoc, 1, false, &vm[0])

	identity := maths.NewMatrix4()

	// Initialize will make an Identity matrix the current matrix ready to be
	// placed on the stack on the first save() call.
	n.transStack.Initialize(identity)

	return nil
}

func (n *nodeManager) configureProjections(world api.IWorld) {
	wp := world.Properties().Window

	// ------------------------------------------------------------
	// Viewport device-space
	// ------------------------------------------------------------
	n.viewport = display.NewViewport()

	n.viewport.SetDimensions(0, 0, wp.DeviceRes.Width, wp.DeviceRes.Height)
	n.viewport.Apply()

	// ------------------------------------------------------------
	// Projection space
	// ------------------------------------------------------------
	n.projection = display.NewCamera()

	camera := world.Properties().Camera
	n.projection.SetProjection(
		0.0, 0.0, // bottom,left
		float32(wp.DeviceRes.Height), float32(wp.DeviceRes.Width), //top,right
		camera.Depth.Near, camera.Depth.Far)

	// ------------------------------------------------------------
	// View-space
	// ------------------------------------------------------------
	offsetX := float32(0.0)
	offsetY := float32(0.0)
	if camera.Centered {
		offsetX = float32(wp.DeviceRes.Width) / 2.0
		offsetY = float32(wp.DeviceRes.Height) / 2.0
		// fmt.Println("NodeManager.configureProjections: center offset: ", offsetX, ",", offsetY)
	}

	// Note: OpenGL's +Y axis is upwards relative to window device.
	world.Viewspace().SetTranslate3Comp(offsetX, offsetY, 1.0)

	// Rarely would you perform a Scale or Rotation on the view-space.
	// But you could if you need to.
	// world.Viewspace().ScaleByComp(1.0, -1.0, 1.0)

	invVSP := world.InvertedViewspace()
	invVSP.Set(world.Viewspace())
	invVSP.Invert()
}

func (n *nodeManager) ClearEnabled(clear bool) {
	n.clearBackground = clear
}

func (n *nodeManager) SetPreNode(node api.INode) {
	n.preNode = node
}

func (n *nodeManager) PreVisit() {
	// Custom node activities, for example, clear background
	// using custom nodes.
	if n.preNode != nil {
		nodeRender, isRenderType := n.preNode.(api.IRender)
		if isRenderType {
			nodeRender.Draw(nil)
		} else {
			log.Fatalf("NodeManager.PreVisit: oops, PreNode '%s' doesn't implement IRender.Draw method", n.preNode)
		}
	}
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
