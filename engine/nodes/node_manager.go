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

	// Stack of nodes
	stack *nodeStack

	transStack api.ITransformStack

	timingTargets api.INodeList
	eventTargets  api.INodeList

	nextScene    api.INode
	currentScene api.INode

	// Used during PreVisit
	preNode api.INode

	// Used during PostVisit
	postNode api.INode

	projection *display.Projection
	viewport   *display.Viewport

	projLoc int32
	viewLoc int32

	preM4  api.IMatrix4
	postM4 api.IMatrix4
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

	o.stack = newNodeStack()
	o.transStack = newTransformStack()

	o.timingTargets = NewNodeList()
	o.eventTargets = NewNodeList()

	o.preM4 = maths.NewMatrix4()
	o.postM4 = maths.NewMatrix4()
	return o
}

func (n *nodeManager) Configure() error {
	// Setup view/projection matrix composition

	n.configureProjections(n.world)

	// -------------------------------------------------------
	// Default Shader
	programID := n.world.Shader().Program()
	n.world.Shader().Use()

	n.projLoc = gl.GetUniformLocation(programID, gl.Str("projection\x00"))
	if n.projLoc < 0 {
		return errors.New("NodeManager: couldn't find 'projection' uniform variable")
	}

	n.viewLoc = gl.GetUniformLocation(programID, gl.Str("view\x00"))
	if n.viewLoc < 0 {
		return errors.New("NodeManager: couldn't find 'view' uniform variable")
	}

	pm := n.projection.Matrix().Matrix()
	gl.UniformMatrix4fv(n.projLoc, 1, false, &pm[0])

	vm := n.world.Viewspace().Matrix()
	gl.UniformMatrix4fv(n.viewLoc, 1, false, &vm[0])
	// -------------------------------------------------------

	// -------------------------------------------------------
	// Texture Shader
	programID = n.world.TextureShader().Program()
	n.world.TextureShader().Use()

	projLoc := gl.GetUniformLocation(programID, gl.Str("projection\x00"))
	if projLoc < 0 {
		return errors.New("NodeManager: couldn't find 'projection' uniform variable")
	}

	viewLoc := gl.GetUniformLocation(programID, gl.Str("view\x00"))
	if viewLoc < 0 {
		return errors.New("NodeManager: couldn't find 'view' uniform variable")
	}

	gl.UniformMatrix4fv(projLoc, 1, false, &pm[0])
	gl.UniformMatrix4fv(viewLoc, 1, false, &vm[0])
	// -------------------------------------------------------

	identity := maths.NewMatrix4()

	// Initialize will make an Identity matrix the current matrix ready to be
	// placed on the stack on the first save() call.
	n.transStack.Initialize(identity)

	dvr := n.world.Properties().Window.DeviceRes
	n.postM4.SetTranslate3Comp(float32(-dvr.Width/2+10.0), float32(-dvr.Height/2)+10.0, 0.0)

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

// Typically called by your "main" code
func (n *nodeManager) SetPreNode(node api.INode) {
	n.preNode = node
}

func (n *nodeManager) SetPostNode(node api.INode) {
	n.postNode = node
}

func (n *nodeManager) PreVisit() {
	// Custom node activities, for example, clear background
	// using custom nodes.
	if n.preNode != nil {
		nodeRender, isRenderType := n.preNode.(api.IRender)
		if isRenderType {
			n.preNode.Update(0.0, 0.0)
			nodeRender.Draw(n.preM4)
		} else {
			log.Fatalf("NodeManager.PreVisit: oops, PreNode '%s' doesn't implement IRender.Draw method", n.preNode)
		}
	}
}

func (n *nodeManager) Begin() bool {
	if n.stack.isEmpty() || n.stack.count() < 2 {
		return false
	}

	n.currentScene = n.stack.pop()
	n.enterNodes(n.currentScene)

	return true
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
	oScene, ok := n.currentScene.(api.IScene)
	if !ok {
		panic("Current Scene '" + n.currentScene.Name() + "' doesn't implement IScene")
	}

	oCurrentState, _ := oScene.State()

	switch oCurrentState {
	case api.SceneOffStage:
		// The current scene is off stage which means we need to tell it
		// to begin transitioning onto the stage.
		ShowState("NM C: ", n.currentScene, " Notify SceneTransitionStartIn")
		n.setSceneState(n.currentScene, api.SceneTransitionStartIn)
	case api.SceneTransitionStartOut:
		// The current scene wants to transition off the stage.
		// Notify it that it can do so.
		ShowState("NM C: ", n.currentScene, " Notify SceneTransitioningOut")
		n.setSceneState(n.currentScene, api.SceneTransitioningOut)

		// At the same time we need to tell the next scene (if there is one) that it can
		// start transitioning onto the stage.
		if n.stack.isEmpty() {
			fmt.Println("---- Stack empty ------")
			n.nextScene = nil
		} else {
			n.nextScene = n.stack.pop()
			fmt.Println("NM C: Pop next scene: ", n.nextScene)
			n.enterNodes(n.nextScene)
			ShowState("NM NS: ", n.nextScene, " Notify SceneTransitionStartIn")
			n.setSceneState(n.nextScene, api.SceneTransitionStartIn)
		}
	case api.SceneExitedStage:
		// The current scene has finished leaving the stage.
		ShowState("NM NS: ", n.currentScene, "")
		n.exitNodes(n.currentScene) // Let it cleanup and exit.

		// Promote next-scene to current-scene
		n.currentScene = n.nextScene

		// Attempt to bring another scene into play.
		if !n.stack.isEmpty() {
			n.nextScene = n.stack.pop()
			fmt.Println("NM NS: Popped next scene: ", n.nextScene)
			n.enterNodes(n.nextScene)
		}
	}

	if n.currentScene != nil && oCurrentState != api.SceneOffStage {
		Visit(n.currentScene, n.transStack, interpolation)
	}

	// --------------------------------------------------------
	// Incoming or next scene
	// --------------------------------------------------------
	if n.nextScene != nil {
		iScene, ok := n.nextScene.(api.IScene)
		if !ok {
			panic("Incoming Scene '" + n.nextScene.Name() + "' doesn't implement IScene")
		}

		iNextState, _ := iScene.State()

		if iNextState != api.SceneOffStage {
			Visit(n.nextScene, n.transStack, interpolation)
		}
	}

	// When the current scene is nil then there are no more scenes
	// to visit, it's time to end the game.
	return n.currentScene != nil
}

// ShowState ---
func ShowState(header string, no api.INode, footer string) {
	scene, _ := no.(api.IScene)

	curr, prev := scene.State()
	fmt.Print(no, " -- ", header)
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

	fmt.Println(footer)
}

func (n *nodeManager) PostVisit() {
	// Custom node activities, for example, FPS visual
	if n.postNode != nil {
		nodeRender, isRenderType := n.postNode.(api.IRender)
		if isRenderType {
			n.postNode.Update(0.0, 0.0)
			nodeRender.Draw(n.postM4)
		} else {
			log.Fatalf("NodeManager.PostVisit: oops, PostNode '%s' doesn't implement IRender.Draw method", n.postNode)
		}
	}
}

// func (n *nodeManager) oldVisit(interpolation float64) bool {
// 	if n.stack.isEmpty() {
// 		return false
// 	}

// 	// fmt.Println("NodeManager: visiting ", m.stack.runningNode)

// 	if n.stack.hasNextNode() {
// 		n.setNextNode()
// 	}

// 	n.transStack.Save()

// 	runningScene := n.stack.runningNode.(api.IScene)

// 	currentState, _ := runningScene.State()

// 	if currentState == api.SceneReplaceTake {
// 		repl := runningScene.GetReplacement()
// 		// fmt.Println("NodeManager: SceneReplaceTake with ", repl)
// 		if repl != nil {
// 			n.stack.replace(repl)
// 			// Immediately switch to the new replacement node
// 			if n.stack.hasNextNode() {
// 				n.setNextNode()
// 			}
// 		} else {
// 			n.exitNodes(n.stack.runningNode)
// 			n.stack.pop()
// 		}
// 	}

// 	// Visit the running node
// 	Visit(n.stack.runningNode, n.transStack, interpolation)

// 	n.transStack.Restore()

// 	return true // continue to draw.
// }

func (n *nodeManager) setSceneState(node api.INode, state int) {
	scene, ok := node.(api.IScene)
	if !ok {
		panic("Scene '" + node.Name() + "' doesn't implement IScene")
	}
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
		n.exitNodes(n.stack.runningNode)
	}

	n.stack.runningNode = n.stack.nextNode
	n.stack.clearNextNode()

	// fmt.Println("NodeManager: new running node ", m.stack.runningNode)

	n.enterNodes(n.stack.runningNode)
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
	return fmt.Sprintf("%s", n.stack)
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
