package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/maths"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
	"github.com/wdevore/Ranger-Go-IGE/extras"
)

type triangleNode struct {
	nodes.Node

	color    []float32
	halfSide float32

	vertices   []float32
	oVertices  []float32
	oVertices2 []float32

	shape    api.IAtlasShape
	centered bool
	filled   bool

	drawBounds bool
	boundsNode api.INode
	afBounds   api.IAffineTransform
	model      api.IMatrix4
	mvp        api.IMatrix4
	mout       api.IMatrix4
}

// NewTriangleNode constructs a generic shape node
func NewTriangleNode(name string, centered, filled bool, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(triangleNode)
	o.Initialize(name)
	o.SetParent(parent)

	o.centered = centered
	o.filled = filled
	o.drawBounds = false
	o.afBounds = maths.NewTransform()
	o.model = maths.NewMatrix4()
	o.mvp = maths.NewMatrix4()
	o.mout = maths.NewMatrix4()

	parent.AddChild(o)

	if err := o.build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *triangleNode) build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	if r.filled {
		name := r.Name() + "::FilledTriangle"
		if r.centered {
			name = r.Name() + "::FilledCenteredTriangle"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.TRIANGLES)
	} else {
		name := r.Name() + "::OutlineTriangle"
		if r.centered {
			name = r.Name() + "::OutlineCenteredTriangle"
		}
		r.shape = world.Atlas().GenerateShape(name, gl.LINE_LOOP)
	}

	// The shape has been added to the atlas but is hasn't been
	// populated with this node's backing data.
	r.populate()

	var err error
	r.boundsNode, err = extras.NewStaticRectangleNode(
		r.vertices[0], r.vertices[1], r.vertices[3], r.vertices[7],
		"ORect", false, world, r.Parent(),
	)
	if err != nil {
		return err
	}
	// [-50 -50 0 50 -50 0 0 31.400002 0]
	// BL(-50.000,-50.000) - TR(50.000,31.400) - Width: 100.000 x Height: 81.400 | At: (0.000,-9.300)

	// r.boundsNode.SetScale(100.0)
	r.boundsNode.SetScaleComps(100.0, 81.4) // 1.228501229
	r.boundsNode.SetPosition(0.000, -9.300)
	r.boundsNode.SetVisible(true)
	rbn := r.boundsNode.(*extras.StaticRectangleNode)
	rbn.SetColor(color.NewPaletteInt64(color.Brick))

	return nil
}

func (r *triangleNode) populate() {
	const centerOffset = 0.0 //float32(math.Pi / 4 / 10)

	r.oVertices = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}
	r.oVertices2 = []float32{
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
		0.0, 0.0, 0.0,
	}

	if r.centered {
		r.halfSide = 0.5
		r.vertices = []float32{
			-r.halfSide, -r.halfSide + centerOffset, 0.0,
			r.halfSide, -r.halfSide + centerOffset, 0.0,
			0.0, 0.314 + centerOffset, 0.0,
		}
	} else {
		const top = float32(math.Pi / 10)
		r.halfSide = 1.0
		r.vertices = []float32{
			0.0, 0.0 + centerOffset, 0.0,
			r.halfSide, 0.0 + centerOffset, 0.0,
			0.0, top + centerOffset, 0.0,
		}
	}

	r.Bounds().SetBounds3D(r.vertices)

	r.shape.SetVertices(r.vertices)

	indices := []uint32{
		0, 1, 2,
	}

	r.shape.SetElementCount(len(indices))

	r.shape.SetIndices(indices)
}

// Vertices returns the shape's vertices
func (r *triangleNode) Vertices() []float32 {
	return r.vertices
}

// SideLength returns the scale length
func (r *triangleNode) SideLength() float32 {
	return r.halfSide * r.Scale() * 2
}

// HalfSide returns the scaled half side length.
func (r *triangleNode) HalfSide() float32 {
	return r.halfSide * r.Scale()
}

// SetColor sets line color
func (r *triangleNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *triangleNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *triangleNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.StaticRenderGraphic)
	renG.SetColor(r.color)
	renG.Render(r.shape, model)

	// This code renders the local BBox not the AABB.
	// nodeRender, isRenderType := r.boundsNode.(api.IRender)
	// if isRenderType {
	// 	nodeRender.Draw(model)
	// }

	r.model.Set(model)

	if r.drawBounds {
		// Transform vertices. We can't use just the Model matrix.
		// We must use the MVP matrix instead.
		vm := r.World().Viewspace()
		pm := r.World().Projection()

		// r.mvp.ToIdentity()
		r.mvp.Set(model)
		maths.Multiply4(pm, vm, r.mout)
		// maths.Multiply4(r.model, r.mout, r.mvp)

		// maths.Multiply4(pm, r.model, r.mout)
		// maths.Multiply4(r.mout, vm, r.mvp)
		// maths.Multiply4(r.mvp, r.model, r.mout)
		// r.mout.ScaleByComp(10.0, 10.0, 1.0)
		// r.model.TranslateBy2Comps(0.0, 1/1.2/10)

		// maths.Multiply4(vm, r.mvp, r.mout)
		// maths.Multiply4(r.mout, pm, r.mvp)
		// r.mvp.TranslateBy2Comps(0.0, 1/1.2/10)
		// r.mvp.ScaleByComp(1.0, 1.2, 1.0)

		// -------------------
		// Set bounds based on local coords not transformed coords
		// ----------------------

		// r.afBounds.SetFromMatrix(r.mvp)
		// r.afBounds.TransformVertices3D(r.vertices, r.oVertices)
		model.TransformVertices3D(r.vertices, r.oVertices)
		// vm.TransformVertices3D(r.oVertices, r.oVertices2)
		// pm.TransformVertices3D(r.oVertices2, r.oVertices)

		// r.mout.TransformVertices3D(r.vertices, r.oVertices)
		// r.model.TransformVertices3D(r.oVertices, r.oVertices2)
		fmt.Println("model ----------------")
		fmt.Println(model)
		fmt.Println("vm ----------------")
		fmt.Println(vm)
		fmt.Println("pm ----------------")
		fmt.Println(pm)
		fmt.Println("mvp ----------------")
		fmt.Println(r.mvp)
		fmt.Println(r.oVertices)

		bounds := r.Bounds()
		// Compute new AA bounds
		bounds.SetBounds3D(r.oVertices)
		fmt.Println(bounds)

		// Rescale and position rectangle node
		r.boundsNode.SetScaleComps(bounds.Width(), bounds.Height())
		r.boundsNode.SetPosition(bounds.Center().X(), bounds.Center().Y())
	}
}
