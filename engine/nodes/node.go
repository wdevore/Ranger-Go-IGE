package nodes

import (
	"fmt"
	"log"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
)

var ids = 0

// Node is an embedded type used by all nodes.
type Node struct {
	id      int
	name    string
	visible bool
	dirty   bool

	parent api.INode
	world  api.IWorld

	Transform
	Group
}

// NewNode constructs a raw base node. Only the Engine should
// construct this base node.
func NewNode() api.INode {
	o := new(Node)
	return o
}

// ID returns the internally generated Id.
func (n *Node) ID() int {
	return n.id
}

// SetID overrides the autogenerated Id
func (n *Node) SetID(id int) {
	n.id = id
}

// Initialize called by base objects from their Initialize
func (n *Node) Initialize(name string) {
	n.id = ids
	ids++
	n.name = name
	n.visible = true
	n.dirty = true

	n.initializeTransform()
	n.initializeGroup()
}

// Build builds this nodes internal geometry
func (n *Node) Build(world api.IWorld) error {
	n.world = world
	return nil
}

// World returns cached world object
func (n *Node) World() api.IWorld {
	return n.world
}

// InitializeWithID called by base objects from their Initialize
func (n *Node) InitializeWithID(id int, name string) {
	n.id = id
	n.name = name
	n.visible = true
	n.dirty = true
}

// Visit traverses "down" the heirarchy while space-mappings traverses upward.
func Visit(node api.INode, transStack api.ITransformStack, interpolation float64) {
	// fmt.Println("Node: visiting ", node)
	if !node.IsVisible() {
		return
	}

	transStack.Save()

	// Because position and angles are dependent
	// on lerping we perform interpolation first.
	node.Interpolate(interpolation)

	aft := node.CalcTransform()

	current := transStack.ApplyAffine(aft)

	nodeRender, isRenderType := node.(api.IRender)
	if isRenderType {
		// Send a copy so the node can optionally modify it.
		nodeRender.Draw(current)
	} else {
		log.Fatalf("Node: oops, %s doesn't implement IRender.Draw method", node)
	}

	children := node.Children()

	if len(children) > 0 {
		for _, child := range children {
			filter, isFilterType := child.(api.IFilter)
			if isFilterType {
				filter.Visit(transStack, interpolation)
			} else {
				Visit(child, transStack, interpolation)
			}
		}
	}

	transStack.Restore()
}

// SetParent binds upward parent.
func (n *Node) SetParent(parent api.INode) {
	n.parent = parent
}

// Parent returns any defined parent
func (n *Node) Parent() api.INode {
	return n.parent
}

// HasParent indicates if this node has a parent, which most do except the root.
func (n *Node) HasParent() bool {
	return n.parent != nil
}

// IsVisible indicates visibility, default is "true"
func (n *Node) IsVisible() bool {
	return n.visible
}

// SetVisible changes the visibility of the node
func (n *Node) SetVisible(visible bool) {
	n.visible = visible
}

// Interpolate is used for blending time based properties.
func (n *Node) Interpolate(interpolation float64) {
	// fmt.Println("Node Interpolate on: ", n)
}

// IsDirty indicates if the node has been modified.
func (n *Node) IsDirty() bool {
	return n.dirty
}

// SetDirty marks a node dirty state.
func (n *Node) SetDirty(dirty bool) {
	n.dirty = dirty
}

// RippleDirty propagates a dirty state to children.
func (n *Node) RippleDirty(dirty bool) {
	for _, child := range n.Children() {
		child.RippleDirty(dirty)
	}

	n.SetDirty(dirty)
}

// Update updates the time properties of a node.
func (n *Node) Update(msPerUpdate, secPerUpdate float64) {
}

// Draw provides a default render--which is to draw nothing.
// You should override this in your custom node if your node
// needs to perform custom rendering.
func (n *Node) Draw(m4 api.IMatrix4) {
	// fmt.Println("Node: Draw ", n)
}

// GetBucket returns a buffer for capturing transformed vertices
func (n *Node) GetBucket() []api.IPoint {
	return nil
}

// Handle may handle an IO event
func (n *Node) Handle(event api.IEvent) bool {
	return false
}

// -----------------------------------------------------
// Scene lifecycles
// -----------------------------------------------------

// Transition specifies what action should happen when transitioning.
// The default is no action.
func (n *Node) Transition() int {
	return api.SceneNoAction
}

// EnterNode called when a node is entering the stage
func (n *Node) EnterNode(man api.INodeManager) {
	// fmt.Println("Node: node enter")
}

// ExitNode called when a node is exiting stage
func (n *Node) ExitNode(man api.INodeManager) {
}

// -----------------------------------------------------
// ITransform defaults
// -----------------------------------------------------

// CalcTransform calculates a matrix based on the
// current transform properties
func (n *Node) CalcTransform() api.IAffineTransform {
	aft := n.aft

	if n.IsDirty() {
		pos := n.position
		aft.MakeTranslate(pos.X(), pos.Y())

		rot := n.rotation
		if rot != 0.0 {
			aft.Rotate(rot)
		}

		sx, sy := n.ScaleComps()
		if sx != 1.0 || sy != 1.0 {
			aft.Scale(sx, sy)
		}

		// Invert...
		aft.InvertTo(n.inverse)
	}

	return aft
}

var p = geometry.NewPoint()

// SetPosition overrides transform's method
func (n *Node) SetPosition(x, y float32) {
	n.Transform.SetPosition(x, y)
	n.RippleDirty(true)
}

// SetRotation overrides transform's method
func (n *Node) SetRotation(radians float64) {
	n.Transform.SetRotation(radians)
	n.RippleDirty(true)
}

// SetScale overrides transform's method
func (n *Node) SetScale(scale float32) {
	n.Transform.SetScale(scale)
	n.RippleDirty(true)
}

// Name returns the node's string name
func (n *Node) Name() string {
	return n.name
}

// -------------------------------------------------------------------
// INodeGroup implementations
// -------------------------------------------------------------------

// -------------------------------------------------------------------
// Misc
// -------------------------------------------------------------------

// PrintTree prints the tree relative to this node.
func PrintTree(node api.INode) {
	fmt.Println("---------- Tree ---------------")
	printBranch(0, node)

	children := node.Children()
	if children != nil {
		printSubTree(children, 1)
	}

	fmt.Println("-------------------------------")
}

func printSubTree(children []api.INode, level int) {
	for _, child := range children {
		subChildren := child.Children()
		printBranch(level, child)
		if subChildren != nil {
			printSubTree(subChildren, level+1)
		}
	}
}

const indent = "   "

func printBranch(level int, node api.INode) {
	// If a node's name begins with "::" then don't print it.
	// This is handy for particle systems or parent nodes with
	// lots of cloned children.
	if node.Name()[0:2] == "::" {
		return
	}

	for i := 0; i < level; i++ {
		fmt.Print(indent)
	}

	fmt.Println(node)
}

func (n Node) String() string {
	return fmt.Sprintf("|'%s' (%d)|", n.name, n.id)
}
