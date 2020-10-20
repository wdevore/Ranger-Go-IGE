package nodes

import (
	"errors"
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/display"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
)

// The node manager is basically the SceneGraph
type nodeManager struct {
	clearBackground bool

	// Stack of nodes
	stack *nodeStack

	transStack api.ITransformStack

	timingTargets api.INodeList
	eventTargets  api.INodeList

	root   api.INode
	scenes api.INode

	nextScene    api.INode
	currentScene api.INode

	projection *display.Projection
	viewport   *display.Viewport

	preM4  api.IMatrix4
	postM4 api.IMatrix4
}

// NewNodeManager constructs a manager for node.
// It manages the lifecycle and events
func NewNodeManager() api.INodeManager {
	o := new(nodeManager)

	// It is very rare that the manager would clear the background
	// because almost all nodes will handle clearing/painting their
	// own backgrounds.
	o.clearBackground = false

	o.stack = newNodeStack()
	o.transStack = newTransformStack()

	o.timingTargets = NewNodeList()
	o.eventTargets = NewNodeList()

	o.preM4 = maths.NewMatrix4()
	o.postM4 = maths.NewMatrix4()
	return o
}

func (n *nodeManager) Configure(world api.IWorld) error {
	// Setup view/projection matrix composition

	n.configureSpaces(world)

	identity := maths.NewMatrix4()

	// Initialize will make an Identity matrix the current matrix ready to be
	// placed on the stack on the first save() call.
	n.transStack.Initialize(identity)

	dvr := world.Properties().Window.DeviceRes
	n.postM4.SetTranslate3Comp(float32(-dvr.Width/2+10.0), float32(-dvr.Height/2)+10.0, 0.0)

	return nil
}

func (n *nodeManager) configureSpaces(world api.IWorld) {
	wp := world.Properties().Window

	// ------------------------------------------------------------
	// Viewport device-space
	// ------------------------------------------------------------
	n.viewport = display.NewViewport()

	n.viewport.SetDimensions(0, 0, wp.DeviceRes.Width, wp.DeviceRes.Height)
	n.viewport.Apply()

	camera := world.Properties().Camera

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
	vpM := world.Viewspace()
	vpM.SetTranslate3Comp(offsetX, offsetY, 1.0)

	// Rarely would you perform a Scale or Rotation on the view-space.
	// But you could if you need to.
	scale := world.Properties().Window.ViewScale
	vpM.ScaleByComp(float32(scale), float32(scale), 1.0)

	invVSP := world.InvertedViewspace()
	invVSP.Set(vpM)
	invVSP.Invert()
}

func (n *nodeManager) ClearEnabled(clear bool) {
	n.clearBackground = clear
}

func (n *nodeManager) SetRoot(root api.INode) {
	n.root = root
}

func (n *nodeManager) Begin() error {
	if n.stack.isEmpty() || n.stack.count() < 2 {
		return errors.New("not enough scenes to start engine. There must be 2 or more")
	}

	// The currentScene is the Incoming scene so add it first
	n.currentScene = n.stack.pop()
	n.enterScene(n.currentScene)

	// We need to set next-scene to the Top incase the game
	// starts with only two scenes.
	n.nextScene = n.stack.top()

	n.scenes = n.root.GetChildByName("Scenes")

	return nil
}

func (n *nodeManager) Visit(interpolation float64) bool {
	n.transStack.Save()

	var visitState bool

	// Up to two scene nodes can run at a time: Outgoing and Incoming.
	visitState = n.continueVisit(interpolation)

	n.transStack.Restore()

	return visitState // continue to draw.
}

func (n *nodeManager) continueVisit(interpolation float64) bool {
	// --------------------------------------------------------
	// Current scene
	// --------------------------------------------------------
	oScene, _ := n.currentScene.(api.IScene)

	oCurrentState, _ := oScene.State()

	switch oCurrentState {
	case api.SceneOffStage:
		// The current scene is off stage which means we need to tell it
		// to begin transitioning onto the stage.
		// ShowState("NM C: ", n.currentScene, " Notify SceneTransitionStartIn")
		n.scenes.InsertAndShift(n.currentScene, 2)
		PrintTree(n.root)

		n.setSceneState(n.currentScene, api.SceneTransitionStartIn)
	case api.SceneTransitionStartOut:
		// The current scene wants to transition off the stage.
		// Notify it that it can do so.
		n.setSceneState(n.currentScene, api.SceneTransitionStartOut)
		// ShowState("NM C: ", n.currentScene, " Notify SceneTransitionStartOut")

		// At the same time we need to tell the next scene (if there is one) that it can
		// start transitioning onto the stage.
		if n.stack.isEmpty() {
			// fmt.Println("---- Stack empty ------")
			n.nextScene = nil
		} else {
			n.nextScene = n.stack.pop()
			// fmt.Println("NM : Popped next scene: ", n.nextScene)
			n.enterScene(n.nextScene)

			// ShowState("NM NS: ", n.nextScene, " Notify SceneTransitionStartIn")
			n.setSceneState(n.nextScene, api.SceneTransitionStartIn)
			n.scenes.InsertAndShift(n.nextScene, 2)
		}
		PrintTree(n.root)
	case api.SceneExitedStage:
		// The current scene has finished leaving the stage.
		// ShowState("NM NS: ", n.currentScene, "")
		// TODO replace "pooled" with "cleanup/dispose"
		pooled := n.exitScene(n.currentScene) // Let it cleanup and exit.

		n.scenes.RemoveLast()

		if pooled {
			// fmt.Println("Returning node to pool: ", n.currentScene)
		}

		// Promote next-scene to current-scene
		// fmt.Println("NM NS: overlay ", n.currentScene, " with ", n.nextScene)
		n.currentScene = n.nextScene
		n.nextScene = nil // This isn't actually needed but it is good form.
		PrintTree(n.root)
	}

	// if n.currentScene != nil && oCurrentState != api.SceneOffStage {
	// 	Visit(n.currentScene, n.transStack, interpolation)
	// }

	// --------------------------------------------------------
	// Incoming or next scene
	// --------------------------------------------------------
	// if n.nextScene != nil {
	// 	iScene, _ := n.nextScene.(api.IScene)

	// 	iNextState, _ := iScene.State()

	// 	if iNextState != api.SceneOffStage {
	// 		Visit(n.nextScene, n.transStack, interpolation)
	// 	}
	// }

	// -------------------------------------------------------
	// Now that visible Scene(s) have been attached to the main Scene
	// node we can Visit the "Root" node.
	// -------------------------------------------------------
	Visit(n.root, n.transStack, interpolation)

	// When the current scene is the last scene to exit the stage
	// then the game is over.
	return n.currentScene != nil
}

// ShowState is for debugging purposes only
func ShowState(header string, no api.INode, footer string) {
	scene, _ := no.(api.IScene)

	curr, prev := scene.State()
	fmt.Print(no, " -- ", header, " {Curr: ")
	switch curr {
	case api.SceneOffStage:
		fmt.Print("SceneOffStage,")
	case api.SceneTransitionStartIn:
		fmt.Print("SceneTransitionStartIn,")
	case api.SceneTransitioningIn:
		fmt.Print("SceneTransitioningIn,")
	case api.SceneOnStage:
		fmt.Print("SceneOnStage,")
	case api.SceneTransitionStartOut:
		fmt.Print("SceneTransitionStartOut,")
	case api.SceneTransitioningOut:
		fmt.Print("SceneTransitioningOut,")
	case api.SceneExitedStage:
		fmt.Print("SceneExitedStage,")
	}

	fmt.Print(" Prev: ")
	switch prev {
	case api.SceneOffStage:
		fmt.Print("SceneOffStage")
	case api.SceneTransitionStartIn:
		fmt.Print("SceneTransitionStartIn")
	case api.SceneTransitioningIn:
		fmt.Print("SceneTransitioningIn")
	case api.SceneOnStage:
		fmt.Print("SceneOnStage")
	case api.SceneTransitionStartOut:
		fmt.Print("SceneTransitionStartOut")
	case api.SceneTransitioningOut:
		fmt.Print("SceneTransitioningOut")
	case api.SceneExitedStage:
		fmt.Print("SceneExitedStage")
	}

	fmt.Println("} ", footer)
}

func (n *nodeManager) setSceneState(node api.INode, state int) {
	scene, _ := node.(api.IScene)
	scene.Notify(state)
}

func (n *nodeManager) PopNode() api.INode {
	return n.stack.pop()
}

func (n *nodeManager) PushNode(node api.INode) {
	n.stack.nextNode = node
	n.stack.push(node)
}

func (n *nodeManager) ReplaceNode(node api.INode) {
	n.stack.replace(node)
}

// --------------------------------------------------------------------------
// Timing
// --------------------------------------------------------------------------

func (n *nodeManager) Update(msPerUpdate, secPerUpdate float64) {
	for _, target := range *n.timingTargets.Items() {
		if target != nil {
			target.Update(msPerUpdate, secPerUpdate)
		}
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

	for _, target := range *n.eventTargets.Items() {
		if target != nil {
			handled := target.Handle(event)

			if handled {
				break
			}
		}
	}
}

func (n *nodeManager) setNextNode() {
	if n.stack.hasRunningNode() {
		n.exitScene(n.stack.runningNode)
	}

	n.stack.runningNode = n.stack.nextNode
	n.stack.clearNextNode()

	// fmt.Println("NodeManager: new running node ", m.stack.runningNode)

	n.enterScene(n.stack.runningNode)
}

// End cleans up NodeManager by clearing the stack and calling all Exits
func (n *nodeManager) End() {
	// Dump the stack
	fmt.Println("End: Cleaning up scene stack.")
	if !n.stack.isEmpty() {
		pn := n.stack.top()

		for pn != nil {
			n.exitScene(pn)
			pn = n.stack.pop()
		}
	}

	n.eventTargets = nil
}

// -----------------------------------------------------
// Scene lifecycles
// -----------------------------------------------------

func (n *nodeManager) enterScene(node api.INode) {
	// fmt.Println("NodeManager: enterScene ", node)
	scene, _ := node.(api.IScene)
	scene.EnterScene(n)

	children := node.Children()
	for _, child := range children {
		n.enterNode(child)
	}
}

func (n *nodeManager) enterNode(node api.INode) {
	// fmt.Println("NodeManager: enterNode ", node)
	node.EnterNode(n)

	children := node.Children()
	for _, child := range children {
		n.enterNode(child)
	}
}

func (n *nodeManager) exitScene(node api.INode) bool {
	// fmt.Println("NodeManager: exitScene ", node)
	scene, _ := node.(api.IScene)
	pooled := scene.ExitScene(n)

	children := node.Children()
	for _, child := range children {
		n.exitNode(child)
	}

	return pooled
}

func (n *nodeManager) exitNode(node api.INode) {
	node.ExitNode(n)

	children := node.Children()
	for _, child := range children {
		n.exitNode(child)
	}
}

func (n *nodeManager) Debug() {
}

func (n nodeManager) String() string {
	return fmt.Sprintf("%s", n.stack)
}
