package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// StaticQuadNode is a generic texture ready node
type StaticQuadNode struct {
	nodes.Node

	color []float32

	shape api.IAtlasShape
}

// NewStaticQuadNode constructs a generic shape node
func NewStaticQuadNode(name string, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(StaticQuadNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	return o, nil
}

// Build configures the node
func (r *StaticQuadNode) Build(world api.IWorld) error {
	r.Node.Build(world)

	r.color = color.NewPaletteInt64(color.White).Array()

	r.shape = world.TextureAtlas().GenerateShape("Quad", gl.TRIANGLES)

	// Populate shape
	r.populate()

	return nil
}

func (r *StaticQuadNode) populate() {
	// These 2D vertices are interleaved with 2D texture coords
	vertices := []float32{
		// Pos    // Tex
		//x  y    z    w
		// 0.0, 0.0, 0.0, 0.0, // CCW
		// 1.0, 0.0, 1.0, 0.0,
		// 1.0, 1.0, 1.0, 1.0,
		// 0.0, 1.0, 0.0, 1.0,
		-0.5, -0.5, 0.0, 0.0, // CCW
		0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.0, 1.0,
	}

	r.shape.SetVertices(vertices)

	indices := []uint32{
		0, 1, 2,
		0, 2, 3,
		// 0, 2, 4,
		// 0, 4, 6,
	}

	r.shape.SetIndices(indices)

	r.shape.SetElementCount(len(indices))
}

// SetColor sets line color
func (r *StaticQuadNode) SetColor(color api.IPalette) {
	r.color = color.Array()
}

// SetAlpha sets the current color's alpha channel 0.0->1.0
func (r *StaticQuadNode) SetAlpha(alpha float32) {
	r.color[3] = alpha
}

// Draw renders shape
func (r *StaticQuadNode) Draw(model api.IMatrix4) {
	renG := r.World().UseRenderGraphic(api.TextureRenderGraphic)
	// renG.SetColor4(r.color)
	// Render texture on quad
	renG.Render(r.shape, model)
}
