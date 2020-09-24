package quadtree

import (
	"fmt"
)

type quadTreeStack struct {
	nodes []*quadTreeNode
}

func newNodeStack() *quadTreeStack {
	o := new(quadTreeStack)
	return o
}

func (q *quadTreeStack) isEmpty() bool {
	return q.count() == 0
}

func (q *quadTreeStack) count() int {
	return len(q.nodes)
}

func (q *quadTreeStack) push(node *quadTreeNode) {
	// fmt.Println("NodeStack: pushing ", n.nextNode)
	q.nodes = append(q.nodes, node)
}

func (q *quadTreeStack) pop() *quadTreeNode {
	if !q.isEmpty() {
		topI := len(q.nodes) - 1 // Top element index
		popNode := q.nodes[topI] // Top element
		q.nodes = q.nodes[:topI] // Pop
		return popNode
	}

	// fmt.Println("NodeStack -- no nodes to pop")

	return nil
}

func (q *quadTreeStack) top() *quadTreeNode {
	topI := len(q.nodes) - 1
	return q.nodes[topI]
}

func (q quadTreeStack) String() string {
	s := "Stack:\n"
	for _, node := range q.nodes {
		s += fmt.Sprint(node)
	}

	return s
}
