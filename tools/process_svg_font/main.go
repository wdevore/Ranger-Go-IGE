package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/geometry"
)

const (
	// Input
	svgFile = "../../assets/ZilSemiSlab.svg"
)

func main() {
	svgFile, err := os.Open(svgFile)
	if err != nil {
		panic(err)
	}

	defer svgFile.Close()

	bytes, err := ioutil.ReadAll(svgFile)
	if err != nil {
		panic(err)
	}

	svg, err := parseSvg(bytes)
	if err != nil {
		panic(err)
	}

	err = processElements(svg)
	if err != nil {
		panic(err)
	}
}

func processElements(svg *Svg) error {
	// First scan elements and find largest AABB across all characters

	// The SVG file must be structured as follows:
	// <?xml...>
	// ...
	// <g
	//    label="A"
	//    id="g65">
	//   <path
	//      id="path33"
	//      d="..."
	//      style="..."/>
	//   <path
	//      id="path34"
	//      d="..."
	//      style="..."/>
	// </g>
	// ...repeat group for each character
	maxAabb := geometry.NewRectangle()
	maxW := -math.MaxFloat32
	maxH := -math.MaxFloat32

	for _, grp := range svg.Groups {
		fmt.Println("--------------- grp: ", grp)

		// Transform vertices to global inkscape-space using group's tranform.
		translateTo(grp, grp.TX, grp.TY)

		// Get the AABB
		b := findAABB(grp)
		fmt.Println(b)
		fmt.Println(b.Left(), ", ", b.Top())

		// Now translate all coords to quadrant #3 with AABB's upper-left
		// to inkscape-space origin.
		translateToQ3(grp, float64(b.Left()), float64(b.Top()))

		// And now translate upwards to quadrant #2 ranger-space
		translateToQ2(grp, float64(b.Bottom()))

		// At this point the char should be in quadrant #2 with the
		// char positioned at the bottom-left.

		// Track max width and height for updating aabb
		maxW = math.Max(maxW, float64(b.Width()))
		maxH = math.Max(maxH, float64(b.Height()))

		// Calc the char's offset delta
		offsetX, offsetY := getOffset(grp, grp.TX, grp.TY)
		fmt.Println("offsets: ", offsetX, ", ", offsetY)

		offsetDx := offsetX - maxAabb.Left()
		offsetDy := offsetY - maxAabb.Top()
		fmt.Println("deltas: ", offsetDx, ", ", offsetDy)

	}

	// Update the aabb so we can normalize each character
	maxAabb.SetBySize(float32(maxW), float32(maxH))

	fmt.Println("L: ", maxAabb.Left(), ", T:", maxAabb.Top(), ", R:", maxAabb.Right(), ", B:", maxAabb.Bottom())

	// Now we need to normalize (aka unit-space) the coords based on the largest AABB.
	// Note: the char could be offset within the unit-space but the offset
	// is used only for positioning during rendering.

	// Write JSON .vfon file

	// --------------------------------------------------------
	// Misc
	// tx = 11.454058-5.5548377 = 5.90 = 5.8992203
	// ty = 18.368204-1.6841292 = = 16.68 = 16.6840748

	// x1 = 14.71715-5.5548377 = 9.1623123
	// y1 = 10.769834-1.6841292 = 9.0857048

	// x2 = 9.1623123+1.919514 = 11.0818263
	// y2 = 9.0857048

	// ix1 = 14.720847-5.5548377 = 9.1660093
	// iy1 = 15.512483-1.6841292 = 13.8283538

	fmt.Println("AABB: ", maxAabb)

	return nil
}

func getOffset(grp group, tx, ty float32) (offX, offY float32) {
	for _, pelem := range grp.Elements {
		switch e := pelem.(type) {
		case *text:
			return e.X + tx, e.Y + ty
		}
	}

	return 0, 0
}

// In Inkscape space
func findAABB(grp group) api.IRectangle {
	aabb := geometry.NewRectangle()

	// fmt.Println("\n---- Group element: ", grp)

	minX := math.MaxFloat32
	minY := math.MaxFloat32
	maxX := -math.MaxFloat32
	maxY := -math.MaxFloat32

	// Each group is a character.
	for _, pelem := range grp.Elements {
		switch e := pelem.(type) {
		// case *text:
		// 	fmt.Println("Text element: ", e)
		case *path:
			miX, miY := findTopLeft(e)
			minX = math.Min(minX, float64(miX))
			minY = math.Min(minY, float64(miY))

			maX, maY := findBottomRight(e)
			maxX = math.Max(maxX, maX)
			maxY = math.Max(maxY, maY)
		}
	}

	aabb.SetMinMax(float32(minX), float32(minY), float32(maxX), float32(maxY))

	return aabb
}

// inkscape space
func findTopLeft(pth *path) (minX, minY float64) {
	minX = math.MaxFloat32
	minY = math.MaxFloat32

	for _, v := range pth.Vertices {
		minX = math.Min(minX, float64(v.x))
		minY = math.Min(minY, float64(v.y))
	}

	return minX, minY
}

// inkscape space
func findBottomRight(pth *path) (maxX, maxY float64) {
	maxX = -math.MaxFloat32
	maxY = -math.MaxFloat32

	for _, v := range pth.Vertices {
		maxX = math.Max(maxX, float64(v.x))
		maxY = math.Max(maxY, float64(v.y))
	}

	return maxX, maxY
}

func translateTo(grp group, tx, ty float32) {
	for _, pelem := range grp.Elements {
		pth, ok := pelem.(*path)
		if ok {
			for _, v := range pth.Vertices {
				v.x += tx
				v.y += tx
			}
		}

		txt, ok := pelem.(*text)
		if ok {
			txt.X += tx
			txt.Y += ty
		}
	}
}

// In inkscape space
func translateToQ3(grp group, minX, minY float64) {
	for _, pelem := range grp.Elements {
		pth, ok := pelem.(*path)
		if ok {
			for _, v := range pth.Vertices {
				v.x -= float32(minX)
				v.y -= float32(minY)
			}
		}
	}
}

// in Ranger space
func translateToQ2(grp group, maxY float64) {
	for _, pelem := range grp.Elements {
		pth, ok := pelem.(*path)
		if ok {
			for _, v := range pth.Vertices {
				v.y -= float32(maxY)
			}
		}
	}
}

func processPath(pth *path, minX, minY float64, aabb api.IRectangle) {

	// bounds := geometry.NewRectangle()
	// Scan vertices of path and find top-left and then translate to quadrant #3.
	// Or find bottom-left and translate to origin which puts the vertices in
	// quadrant #2.

	// bounds.Expand(v.x, v.y)
	// aabb.Expand(v.x, v.y)
	// fmt.Println("Bounds: ", bounds)

}

func parseSvg(bytes []byte) (*Svg, error) {
	svg := new(Svg)

	err := xml.Unmarshal(bytes, svg)
	if err != nil {
		return nil, err
	}
	return svg, nil
}
