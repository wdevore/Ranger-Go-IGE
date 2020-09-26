package quadtree

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
)

var (
	nodeStack = newNodeStack()
)

type quadTree struct {
	// How many items can be held
	capacity int // Currently not used
	maxDepth int

	root *quadTreeNode
}

// NewQuadTree creates a new QuadTree object.
func NewQuadTree() api.IQuadTree {
	o := new(quadTree)
	o.initialize()
	return o
}

func (q *quadTree) initialize() {
	q.root = newQuadTreeNode()

	// Pre populate the stack pool
	for i := 0; i < 100; i++ {
		n := newQuadTreeNode()
		nodeStack.push(n)
	}
}

// Add returns true if the 'node' was added based on the boundary of the tree.
func (q *quadTree) Add(node api.INode) bool {
	return q.root.add(node, q.maxDepth) == quadAdded
}

// Remove node from tree
func (q *quadTree) Remove(node api.INode) bool {
	return q.root.remove(node)
}

// Query returns a collection of INodes based on a AABB rectangle.
func (q *quadTree) Query(boundary api.IRectangle, nodes *[]api.INode) {
	q.root.query(boundary, nodes)
}

// Clear entire tree by removing both empty quadrants and items
func (q *quadTree) Clear() {
	q.root.clearQuadrant(q.root, 0)
}

// Clean entire tree by removing both empty quadrants
func (q *quadTree) Clean() {
	q.root.cleanQuadrant(q.root, 0)
}

func (q *quadTree) SetBoundary(x, y, w, h float32) {
	q.root.boundary.Set(x, y, w, h)
}

func (q *quadTree) Boundary() api.IRectangle {
	return q.root.boundary
}

// SetCapacity sets the limit on how many INodes can be held
// by any one quadrant. Once a quadrant is 'filled' the quadrant
// must be divided with all the INodes being redistributed among
// the new divided quadrants.
// TODO -- not implemented yet.
func (q *quadTree) SetCapacity(capacity int) {
	q.capacity = capacity
}

func (q *quadTree) Capacity() int {
	return q.capacity
}

func (q *quadTree) SetMaxDepth(depth int) {
	q.maxDepth = depth
}

func (q *quadTree) MaxDepth() int {
	return q.maxDepth
}

// String returns a formatted string of the tree.
// ---------- QuadTree ---------------
// Root {0} [500.000, 500.000] Div: true
//    |'Rect' (2)|
//    Quad1 {1}
//    Quad2 {1}
//       Quad1 {2}
//       Quad2 {2}
//       Quad3 {2}
//       Quad4 {2}
//          Quad1 {3}
//          Quad2 {3}
//          Quad3 {3}
//          Quad4 {3}
//             Quad1 {4}
//             Quad2 {4} [31.250, 31.250] Div: true
//                |'Rect' (1)|
//                Quad1 {5}
//                Quad2 {5}
//                Quad3 {5}
//                Quad4 {5}
//             Quad3 {4}
//             Quad4 {4}
//    Quad3 {1}
//       Quad1 {2}
//       Quad2 {2}
//       Quad3 {2}
//       Quad4 {2}
//          Quad1 {3}
//          Quad2 {3}
//          Quad3 {3}
//          Quad4 {3}
//             Quad1 {4}
//             Quad2 {4}
//             Quad3 {4}
//             Quad4 {4}
//                Quad1 {5}
//                Quad2 {5}
//                Quad3 {5}
//                Quad4 {5} [15.625, 15.625] Div: false
//                   |'Rect' (0)|   <-- removed
//    Quad4 {1}
// ---------- QuadTree ---------------
func (q quadTree) String() string {
	s := fmt.Sprintf("---------- QuadTree ---------------\n")
	s += fmt.Sprintf("%s\n", q.root.toString("Root", 0))
	s += fmt.Sprintf("---------- QuadTree ---------------\n")
	return s
}
