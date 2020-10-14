package generators

import (
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"
)

// These shape generators define vertices and indices that represent
// unit-sized vector shapes. These shapes are then loaded into an Atlas.

// GenerateUnitHLineVectorShape builds a horizontal unit length line.
func GenerateUnitHLineVectorShape() (vertices []float32, indices []uint32, mode int) {
	vertices = []float32{
		-0.5, 0.0, 0.0,
		0.5, 0.0, 0.0,
	}

	indices = []uint32{
		0, 1,
	}

	return vertices, indices, gl.LINES
}

// GenerateUnitVLineVectorShape builds a vertical unit length line.
func GenerateUnitVLineVectorShape() (vertices []float32, indices []uint32, mode int) {
	vertices = []float32{
		0.0, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	indices = []uint32{
		0, 1,
	}

	return vertices, indices, gl.LINES
}

// GenerateUnitRectangleVectorShape builds a rectangle unit length line.
func GenerateUnitRectangleVectorShape(centered bool, forFilling bool) (vertices []float32, indices []uint32, mode int) {
	if centered {
		vertices = []float32{
			-0.5, -0.5, 0.0,
			0.5, -0.5, 0.0,
			0.5, 0.5, 0.0,
			-0.5, 0.5, 0.0,
		}
	} else {
		vertices = []float32{
			0.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			1.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
		}
	}

	// These indices can be used with the same vertices. You don't need
	// separate vertex arrays.
	if forFilling {
		indices = []uint32{
			0, 1, 2,
			0, 2, 3,
		}
		mode = gl.TRIANGLES
	} else {
		indices = []uint32{
			0, 1, 2, 3, // CCW
		}
		mode = gl.LINE_LOOP
	}

	return vertices, indices, mode
}

// GenerateUnitTriangleVectorShape builds a triangle with equal length sides.
func GenerateUnitTriangleVectorShape(forFilling bool) (vertices []float32, indices []uint32, mode int) {
	vertices = []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, math.Pi / 10.0, 0.0,
	}

	indices = []uint32{
		0, 1, 2,
	}

	if forFilling {
		mode = gl.TRIANGLES
	} else {
		mode = gl.LINE_LOOP
	}

	return vertices, indices, mode
}

// GenerateUnitPlusVectorShape builds a plus-sign of unit length
func GenerateUnitPlusVectorShape() (vertices []float32, indices []uint32, mode int) {
	vertices = []float32{
		-0.5, 0.0, 0.0,
		0.5, 0.0, 0.0,
		0.0, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	indices = []uint32{
		0, 1, 2, 3,
	}

	return vertices, indices, gl.LINES
}

// GenerateUnitZBarVectorShape builds a "Z" shape for testing purposes.
func GenerateUnitZBarVectorShape(forFilling bool) (vertices []float32, indices []uint32, mode int) {
	vertices = []float32{
		-0.1, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.5, -0.4, 0.0,
		0.1, -0.4, 0.0,
		0.1, 0.5, 0.0,
		-0.5, 0.5, 0.0,
		-0.5, 0.4, 0.0,
		-0.1, 0.4, 0.0,
	}

	// These indices can be used with the same vertices. You don't need
	// separate vertex arrays.
	if forFilling {
		indices = []uint32{
			0, 1, 2,
			2, 3, 0,
			0, 3, 4,
			0, 4, 7,
			6, 7, 4,
			6, 4, 5,
		}
		mode = gl.TRIANGLES
	} else {
		indices = []uint32{
			0, 1, 2, 3, 4, 5, 6, 7, // CCW
		}
		mode = gl.LINES
	}

	return vertices, indices, mode
}

// GenerateUnitCircleVectorShape builds a circle with radius 0.5.
func GenerateUnitCircleVectorShape(segments int, forFilling bool) (vertices []float32, indices []uint32, mode int) {
	radius := 0.5 // diameter of 1.0
	step := math.Pi / float64(segments)

	index := uint32(0)

	if forFilling {
		// Filled circles have a center point for the Fan fill algorithm
		vertices = append(vertices, 0.0, 0.0, 0.0)

		// Reference the center point
		indices = append(indices, 0)

		index++
		mode = gl.TRIANGLES
	} else {
		mode = gl.LINES
	}

	for i := 0.0; i < 2.0*math.Pi; i += step {
		x := math.Cos(i) * radius
		y := math.Sin(i) * radius
		vertices = append(vertices, float32(x), float32(y), 0.0)
		indices = append(indices, index)
		index++
	}

	return vertices, indices, mode
}

// GenerateUnitArcVectorShape builds a arc/pie with radius 0.5.
func GenerateUnitArcVectorShape(startAngle, endAngle float64, segments int, forFilling bool) (vertices []float32, indices []uint32, mode int) {
	radius := 0.5 // diameter of 1.0
	step := (endAngle - startAngle) / float64(segments)

	index := uint32(0)

	vertices = append(vertices, 0.0, 0.0, 0.0)

	// Reference the center point
	indices = append(indices, 0)

	if forFilling {
		mode = gl.TRIANGLES
	} else {
		mode = gl.LINE_LOOP
	}

	index++

	for i := startAngle; i <= startAngle+endAngle; i += step {
		x := math.Cos(i) * radius
		y := math.Sin(i) * radius
		vertices = append(vertices, float32(x), float32(y), 0.0)
		indices = append(indices, index)
		index++
	}

	return vertices, indices, mode
}
