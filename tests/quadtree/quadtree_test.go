package main

import (
	"fmt"
	"testing"

	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
	"github.com/wdevore/Ranger-Go-IGE/extras"
	"github.com/wdevore/Ranger-Go-IGE/extras/quadtree"
)

// go test -v -count=1 quadtree_test.go

func TestRunner(t *testing.T) {
	testQuadtree(t)
}

func testQuadtree(t *testing.T) {
	tree := quadtree.NewQuadTree()
	tree.SetMaxDepth(5)
	tree.SetBoundary(0.0, 0.0, 500.0, 500.0)

	node0, _ := extras.NewStaticNilNode("Rect")
	node0.SetPosition(10.0+250.0, 10.0)
	node0.SetBoundBySize(10, 10)
	tree.Add(node0)

	node1, _ := extras.NewStaticNilNode("Rect")
	node1.SetPosition(50.0+250.0, 50.0+250.0)
	node1.SetBoundBySize(20, 20)
	tree.Add(node1)

	// Add a node directly in the center insures that the node
	// is placed at the root as it can't be placed in any of
	// the quadrants.
	node2, _ := extras.NewStaticNilNode("Rect")
	node2.SetPosition(250.0, 250.0)
	node2.SetBoundBySize(20, 20)
	tree.Add(node2)

	// fmt.Println("Clean 1 ------------------------------")
	// tree.Clear()
	// fmt.Println(tree)

	tree.Remove(node0)
	fmt.Println(tree)

	fmt.Println("Clean 2 ------------------------------")
	tree.Clean()
	fmt.Println(tree)

	// tree.Remove(node1)
	// fmt.Println(tree)

	// tree.Remove(node2)
	// fmt.Println(tree)
}

func testRectangleIntersects(t *testing.T) {
	r := geometry.NewRectangle()
	r.Set(100.0, 100.0, 100.0, 100.0)

	o := geometry.NewRectangle()
	o.Set(20.0, 20.0, 50.0, 50.0)

	intersects := r.Intersects(o)

	if intersects {
		t.Fatal("Wasn't expecting an intersection")
	}

	o.Set(120.0, 120.0, 50.0, 50.0)

	intersects = r.Intersects(o)

	if !intersects {
		t.Fatal("Was expecting an intersection")
	}

	o.Set(50.0, 50.0, 50.0, 50.0)

	intersects = r.Intersects(o)

	if !intersects {
		t.Fatal("Was expecting an intersection")
	}

	o.Set(75.0, 150.0, 50.0, 50.0)

	intersects = r.Intersects(o)

	if !intersects {
		t.Fatal("Was expecting an intersection")
	}

	fmt.Println("------------------------------ Contains ")

	o.Set(110.0, 110.0, 5.0, 5.0)

	contains := r.Contains(o)

	if !contains {
		t.Fatal("Was expecting Contains")
	}

	o.Set(290.0, 290.0, 175.0, 175.0)
	contains = r.Contains(o)

	if contains {
		t.Fatal("Wasn't expecting Contains")
	}

	o.Set(110.0, 110.0, 50.0, 50.0)

	contains = r.Contains(o)

	if !contains {
		t.Fatal("Was expecting Contains")
	}

	o.Set(110.0, 110.0, 150.0, 150.0)

	contains = r.Contains(o)

	if contains {
		t.Fatal("Wasn't expecting Contains")
	}
}

func testPointsInRectangle(t *testing.T) {

	r := geometry.NewRectangle()
	r.Set(10.0, 10.0, 100.0, 100.0)

	p := geometry.NewPoint()
	p.SetByComp(5.0, 5.0)

	inside := r.PointContained(p)

	if inside {
		t.Fatal("Wasn't expecting point inside")
	}

	p.SetByComp(15.0, 15.0)

	inside = r.PointContained(p)

	if !inside {
		t.Fatal("Was expecting point inside")
	}

	p.SetByComp(150.0, 15.0)

	inside = r.PointContained(p)

	if inside {
		t.Fatal("Wasn't expecting point inside")
	}

	p.SetByComp(10.0, 10.0)

	inside = r.PointContained(p)

	if inside {
		t.Fatal("Wasn't expecting point inside")
	}

	p.SetByComp(110.0, 10.0)

	inside = r.PointContained(p)

	if inside {
		t.Fatal("Wasn't expecting point inside")
	}

	p.SetByComp(50.0, 50.0)

	inside = r.PointContained(p)

	if !inside {
		t.Fatal("Was expecting point inside")
	}

	// ----------------------- PointInside()
	p.SetByComp(50.0, 50.0)

	inside = r.PointInside(p)

	if !inside {
		t.Fatal("Was expecting point inside")
	}

	p.SetByComp(10.0, 50.0)

	inside = r.PointInside(p)

	if !inside {
		t.Fatal("Was expecting point inside")
	}

	p.SetByComp(110.0, 50.0)

	inside = r.PointInside(p)

	if inside {
		t.Fatal("Wasn't expecting point inside")
	}
}
