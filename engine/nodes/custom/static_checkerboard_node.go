package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// TODO replace with a interleved vertex/color atlas
// StaticCheckerboardNode is a checkerboard
type StaticCheckerboardNode struct {
	nodes.Node

	evenColor []float32
	oddColor  []float32
	size      float32

	shape api.IAtlasShape

	m4 api.IMatrix4
}

// NewStaticCheckerboardNode constructs a generic shape node
func NewStaticCheckerboardNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticCheckerboardNode)
	o.Initialize(name)
	o.SetParent(parent)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticCheckerboardNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.evenColor = color.NewPaletteInt64(color.LightGray).Array()  // Lighter
	r.oddColor = color.NewPaletteInt64(color.LighterGray).Array() // Darker

	r.m4 = maths.NewMatrix4()
	r.size = 250

	name := "Checkerboard"
	r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	return nil
}

func (r *StaticCheckerboardNode) populate() {
	var vertices []float32

	vertices = []float32{
		0.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
	}

	r.shape.SetVertices(vertices)

	var indices []uint32

	indices = []uint32{
		0, 1, 2,
		0, 2, 3,
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// SetEvenColor sets even color
func (r *StaticCheckerboardNode) SetEvenColor(color api.IPalette) {
	r.evenColor = color.Array()
}

// SetEvenAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticCheckerboardNode) SetEvenAlpha(alpha float32) {
	r.evenColor[3] = alpha
}

// SetOddColor sets odd color
func (r *StaticCheckerboardNode) SetOddColor(color api.IPalette) {
	r.oddColor = color.Array()
}

// SetOddAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticCheckerboardNode) SetOddAlpha(alpha float32) {
	r.oddColor[3] = alpha
}

// PointInside checks if point inside shape's polygon
func (r *StaticCheckerboardNode) PointInside(p api.IPoint) bool {
	return false
}

// SetSize sets the size of a check box
func (r *StaticCheckerboardNode) SetSize(size float32) {
	r.size = size
}

// Draw renders shape
func (r *StaticCheckerboardNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)

	xflip := false
	yFlip := true

	// Re-render the same square at different positions.
	dvr := r.World().Properties().Window.DeviceRes
	hw := float32(dvr.Width) / 2
	hh := float32(dvr.Height) / 2

	// Top-Right quadrant +,+
	for y := float32(0); y < hh; y += r.size {
		for x := float32(0); x < hw; x += r.size {
			r.m4.SetTranslate3Comp(x, y, 0.0)
			r.m4.ScaleByComp(r.size, r.size, 1.0)

			if xflip {
				renG.SetColor4(r.evenColor) // lighter
			} else {
				renG.SetColor4(r.oddColor) // darker
			}

			renG.Render(r.shape, r.m4)
			xflip = !xflip
		}
		xflip = yFlip
		yFlip = !yFlip
	}

	xflip = true
	yFlip = false
	// Bottom-Right quadrant +,-
	for y := float32(-r.size); y > -hh-r.size; y -= r.size {
		for x := float32(0); x < hw; x += r.size {
			r.m4.SetTranslate3Comp(x, y, 0.0)
			r.m4.ScaleByComp(r.size, r.size, 1.0)

			if xflip {
				renG.SetColor4(r.evenColor)
			} else {
				renG.SetColor4(r.oddColor)
			}

			renG.Render(r.shape, r.m4)
			xflip = !xflip
		}
		xflip = yFlip
		yFlip = !yFlip
	}

	xflip = true
	yFlip = false
	// Top-Left quadrant -,+
	for y := float32(0); y < hh; y += r.size {
		for x := float32(-r.size); x > -hw-r.size; x -= r.size {
			r.m4.SetTranslate3Comp(x, y, 0.0)
			r.m4.ScaleByComp(r.size, r.size, 1.0)

			if xflip {
				renG.SetColor4(r.evenColor) // lightest
			} else {
				renG.SetColor4(r.oddColor) // lighter (darker)
			}

			renG.Render(r.shape, r.m4)
			xflip = !xflip
		}
		xflip = yFlip
		yFlip = !yFlip
	}

	xflip = false
	yFlip = true
	// Bottom-Left quadrant -,-
	for y := float32(-r.size); y > -hh-r.size; y -= r.size {
		for x := float32(-r.size); x > -hw-r.size; x -= r.size {
			r.m4.SetTranslate3Comp(x, y, 0.0)
			r.m4.ScaleByComp(r.size, r.size, 1.0)

			if xflip {
				renG.SetColor4(r.evenColor)
			} else {
				renG.SetColor4(r.oddColor)
			}

			renG.Render(r.shape, r.m4)
			xflip = !xflip
		}
		xflip = yFlip
		yFlip = !yFlip
	}
}
