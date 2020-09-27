package quadtree

import (
	"fmt"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
)

//
// .----------------------.----------------------.
// |                      |                      |
// |                      |                      |
// |                      |                      |
// |                      |                      |
// |       Quadrant       |       Quadrant       |
// |          1           |          2           |
// |                      |                      |
// |                      |                      |
// |                      |                      |
// .----------------------.----------------------.
// |                      |                      |
// |                      |                      |
// |                      |                      |
// |                      |                      |
// |       Quadrant       |       Quadrant       |
// |          4           |          3           |
// |                      |                      |
// |                      |                      |
// |                      |                      |
// .----------------------.----------------------.

const (
	quadNoFit = iota
	quadAdded

	indent = "   "
)

type quadTreeNode struct {
	quadrant1 *quadTreeNode
	quadrant2 *quadTreeNode
	quadrant3 *quadTreeNode
	quadrant4 *quadTreeNode

	// INodes held by this tree node
	nodes    []api.INode
	hasNodes bool

	divided bool

	// boundary is considered the parent of the quadrants.
	boundary api.IRectangle
	id       int // For debugging
}

func newQuadTreeNode() *quadTreeNode {
	o := new(quadTreeNode)
	o.divided = false
	o.boundary = geometry.NewRectangle()
	return o
}

// We decend until we either reach max-depth or we can't fit at the current
// depth.
// If we can't fit and we haven't reached max-depth then the node is placed
// in the parent of the current quadrant, which could even be the Root.
// Nodes that "sit" on boundaries will aways be placed in the parent.
func (q *quadTreeNode) add(node api.INode, depth int) int {
	// fmt.Println("quad: ", q.id, ", depth: ", depth, ", bounds: ", q.boundary)

	// Note: if depth == 0 Then the node must fit into one of the quadrants.
	// If it doesn't fit then it is added to the parent.

	nBounds := node.Bounds()
	// See if the INode's boundary fits this quadrant.
	fits := q.boundary.Contains(nBounds)
	state := quadNoFit

	if fits {
		if depth == 0 {
			// We have reached max depth and it fits this quadrant so just collect here.
			// fmt.Println("added ", node, " at Max, quad: ", q.id, ", depth: ", depth, ", bounds: ", q.boundary)
			q.nodes = append(q.nodes, node)
			q.hasNodes = true
			return quadAdded
		}

		// We haven't reached max depth so we need to sub divide prior to decending.
		if !q.divided {
			q.divide()
		}

		// --------------------------------------------------
		// Now attempt to place the node in one of the quadrants.
		// If it doesn't fit in any of the quadrants then place
		// it in the parent (aka this quadrant).
		// --------------------------------------------------

		state1 := q.quadrant1.add(node, depth-1)
		if state1 == quadAdded {
			return state1
		}

		state2 := q.quadrant2.add(node, depth-1)
		if state2 == quadAdded {
			return state2
		}

		state3 := q.quadrant3.add(node, depth-1)
		if state3 == quadAdded {
			return state3
		}

		state4 := q.quadrant4.add(node, depth-1)
		if state4 == quadAdded {
			return state4
		}

		if state1 == quadNoFit && state2 == quadNoFit && state3 == quadNoFit && state4 == quadNoFit {
			// The node didn't fit in any of the quadrants.
			// fmt.Println("added ", node, " quad: ", q.id, ", depth: ", depth, ", bounds: ", q.boundary)
			q.nodes = append(q.nodes, node)
			state = quadAdded
			q.hasNodes = true
		}
	}

	return state
}

func (q *quadTreeNode) remove(node api.INode) bool {
	removed := q.removeNode(node)
	if removed {
		return true
	}

	// Use node's AABB to determine which quadrant to check.
	if q.quadrant1.boundary.Contains(node.Bounds()) {
		removed = q.quadrant1.remove(node)
		if removed {
			return true
		}
	}

	if q.quadrant2.boundary.Contains(node.Bounds()) {
		removed = q.quadrant2.remove(node)
		if removed {
			return true
		}
	}

	if q.quadrant3.boundary.Contains(node.Bounds()) {
		removed = q.quadrant3.remove(node)
		if removed {
			return true
		}
	}

	if q.quadrant4.boundary.Contains(node.Bounds()) {
		removed = q.quadrant4.remove(node)
		if removed {
			return true
		}
	}

	return false
}

func (q *quadTreeNode) query(boundary api.IRectangle, nodes *[]api.INode) {
	// Collect INodes in this quadrant.
	if len(q.nodes) > 0 {
		for _, n := range q.nodes {
			if boundary.Intersects(n.Bounds()) {
				*nodes = append(*nodes, n)
			}
		}
	}

	if !q.divided {
		return
	}

	// As long as the boundary intersects the quadrant then there is a
	// potential for INodes to intersect 'boundary'.
	if q.quadrant1.boundary.Intersects(boundary) {
		q.quadrant1.query(boundary, nodes)
	}

	if q.quadrant2.boundary.Intersects(boundary) {
		q.quadrant2.query(boundary, nodes)
	}

	if q.quadrant3.boundary.Intersects(boundary) {
		q.quadrant3.query(boundary, nodes)
	}

	if q.quadrant4.boundary.Intersects(boundary) {
		q.quadrant4.query(boundary, nodes)
	}
}

func (q *quadTreeNode) traverse(quadrantCB api.QuadrantBoundsFunc) {
	quadrantCB(q.boundary)

	if !q.divided {
		return
	}

	q.quadrant1.traverse(quadrantCB)
	q.quadrant2.traverse(quadrantCB)
	q.quadrant3.traverse(quadrantCB)
	q.quadrant4.traverse(quadrantCB)
}

func (q *quadTreeNode) removeNode(node api.INode) bool {
	// Does the node exist?
	remIn := -1
	for i, n := range q.nodes {
		if n == node {
			remIn = i
			break
		}
	}

	if remIn >= 0 {
		// Remove it and adjust array
		q.nodes[remIn] = q.nodes[len(q.nodes)-1] // Copy last element to index i.
		q.nodes[len(q.nodes)-1] = nil            // Erase last element (write zero value).
		q.nodes = q.nodes[:len(q.nodes)-1]       // Truncate slice.

		q.hasNodes = len(q.nodes) > 0

		return true
	}

	return false
}

// clearQuadrant removes both empty quadrants and items
func (q *quadTreeNode) clearQuadrant(quad *quadTreeNode, lvl int) {
	// fmt.Println("Enter Quad: ", quad.id, ", Lvl: ", lvl)
	if quad.divided {
		q.clearQuadrant(quad.quadrant1, lvl+1)
		q.clearQuadrant(quad.quadrant2, lvl+1)
		q.clearQuadrant(quad.quadrant3, lvl+1)
		q.clearQuadrant(quad.quadrant4, lvl+1)

		nodeStack.push(quad.quadrant1)
		quad.quadrant1 = nil
		nodeStack.push(quad.quadrant2)
		quad.quadrant2 = nil
		nodeStack.push(quad.quadrant3)
		quad.quadrant3 = nil
		nodeStack.push(quad.quadrant4)
		quad.quadrant4 = nil

		quad.divided = false
		quad.hasNodes = false
		quad.nodes = nil
	}
}

// cleanQuadrant Only removes empty quadrants
func (q *quadTreeNode) cleanQuadrant(quad *quadTreeNode, lvl int) bool {
	// fmt.Println("Enter Quad: ", quad.id, ", Lvl: ", lvl)
	// Decend down to the leaf. If all quadrants in the decent path
	// had zero items then we can remove all quadrants along that
	// path.
	hasItems := quad.hasNodes

	// if hasItems {
	// 	fmt.Println("Has Item: Quad: ", quad.id, ", Lvl: ", lvl)
	// }

	if quad.divided {
		hi1 := q.cleanQuadrant(quad.quadrant1, lvl+1)
		hasItems = hasItems || hi1
		hi2 := q.cleanQuadrant(quad.quadrant2, lvl+1)
		hasItems = hasItems || hi2
		hi3 := q.cleanQuadrant(quad.quadrant3, lvl+1)
		hasItems = hasItems || hi3
		hi4 := q.cleanQuadrant(quad.quadrant4, lvl+1)
		hasItems = hasItems || hi4

		hi := hi1 || hi2 || hi3 || hi4

		if !hi {
			// fmt.Println("Removing quads: ", quad.id, ", Lvl: ", lvl)
			nodeStack.push(quad.quadrant1)
			quad.quadrant1 = nil
			nodeStack.push(quad.quadrant2)
			quad.quadrant2 = nil
			nodeStack.push(quad.quadrant3)
			quad.quadrant3 = nil
			nodeStack.push(quad.quadrant4)
			quad.quadrant4 = nil

			quad.divided = false
		}
	}

	return hasItems
}

func (q *quadTreeNode) divide() {
	// q.boundary is the encompassing bounds
	hw := q.boundary.Width() / 2.0
	hh := q.boundary.Height() / 2.0
	cx := q.boundary.Left() + hw
	cy := q.boundary.Bottom() + hh

	if !nodeStack.isEmpty() {
		q.quadrant1 = nodeStack.pop()
	} else {
		q.quadrant1 = newQuadTreeNode()
	}
	q.quadrant1.id = 1
	q.quadrant1.boundary.SetMinMax(
		q.boundary.Left(), cy,
		cx, q.boundary.Top(),
	)

	if !nodeStack.isEmpty() {
		q.quadrant2 = nodeStack.pop()
	} else {
		q.quadrant2 = newQuadTreeNode()
	}
	q.quadrant2.id = 2
	q.quadrant2.boundary.SetMinMax(
		cx, cy,
		q.boundary.Right(), q.boundary.Top(),
	)

	if !nodeStack.isEmpty() {
		q.quadrant3 = nodeStack.pop()
	} else {
		q.quadrant3 = newQuadTreeNode()
	}
	q.quadrant3.id = 3
	q.quadrant3.boundary.SetMinMax(
		cx, q.boundary.Bottom(),
		q.boundary.Right(), cy,
	)

	if !nodeStack.isEmpty() {
		q.quadrant4 = nodeStack.pop()
	} else {
		q.quadrant4 = newQuadTreeNode()
	}
	q.quadrant4.id = 4
	q.quadrant4.boundary.SetMinMax(
		q.boundary.Left(), q.boundary.Bottom(),
		cx, cy,
	)

	q.divided = true
}

func (q *quadTreeNode) toString(header string, level int) string {
	// Print any nodes held by this quadrant indented
	idn := ""
	for i := level * 3; i > 0; i-- {
		idn += " "
	}

	var s = ""
	if len(q.nodes) > 0 {
		s = fmt.Sprintf("%s%s {%d} [%0.3f, %0.3f] Div: %t\n", idn, header, level, q.boundary.Width(), q.boundary.Height(), q.divided)
	} else {
		s = fmt.Sprintf("%s%s {%d}\n", idn, header, level)
	}

	for _, n := range q.nodes {
		s += idn + indent + fmt.Sprint(n) + "\n"
	}

	if q.divided {
		s += q.quadrant1.toString("Quad1", level+1)
		s += q.quadrant2.toString("Quad2", level+1)
		s += q.quadrant3.toString("Quad3", level+1)
		s += q.quadrant4.toString("Quad4", level+1)
	}

	return s
}
