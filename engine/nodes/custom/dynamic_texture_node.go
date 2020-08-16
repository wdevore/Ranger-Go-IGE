package custom

import (
	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/wdevore/Ranger-Go-IGE/api"
	"github.com/wdevore/Ranger-Go-IGE/engine/nodes"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/color"
)

// DynamicTextureNode is a dynamic texture ready node.
// This node contains an index to a TextureAtlas.
type DynamicTextureNode struct {
	nodes.Node

	shape api.IAtlasShape

	textureMan         api.ITextureManager
	verticesAndTexture []float32
	textureIndexes     []int
	index              int

	color []float32
}

// NewDynamicTextureNode constructs a generic shape node
func NewDynamicTextureNode(name string, startIndex int, textureMan api.ITextureManager, world api.IWorld, parent api.INode) (api.INode, error) {
	o := new(DynamicTextureNode)
	o.Initialize(name)
	o.SetParent(parent)
	parent.AddChild(o)

	if err := o.Build(world); err != nil {
		return nil, err
	}

	o.index = startIndex
	o.textureMan = textureMan

	o.color = color.NewPaletteInt64(color.Transparent).Array()

	return o, nil
}

// Build configures the node
func (d *DynamicTextureNode) Build(world api.IWorld) error {
	d.Node.Build(world)

	d.shape = world.ShapeAtlas().GenerateShape("DynamicTextureQuad", gl.TRIANGLES)

	return nil
}

// Populate ...
func (d *DynamicTextureNode) Populate() {
	// These 2D vertices are interleaved with 2D texture coords
	// The s,t coords are sourced by the manifest based on index
	idx := d.textureIndexes[d.index]
	coords := d.textureMan.GetSTCoords(0, idx)
	c := *coords
	d.verticesAndTexture = []float32{
		// Pos           Tex
		//x  y       z/s    w/t
		-0.5, -0.5, c[0], c[1], // CCW
		0.5, -0.5, c[2], c[3],
		0.5, 0.5, c[4], c[5],
		-0.5, 0.5, c[6], c[7],
	}

	// d.verticesAndTexture = []float32{
	// 	// Pos           Tex
	// 	//x  y       z/s    w/t
	// 	// -0.5, -0.5, 0.0, 0.0, // CCW
	// 	// 0.5, -0.5, 1.0, 0.0,
	// 	// 0.5, 0.5, 1.0, 1.0,
	// 	// -0.5, 0.5, 0.0, 1.0,
	// }

	// fmt.Println(d.verticesAndTexture)
	d.shape.SetVertices(d.verticesAndTexture)

	indices := []uint32{
		0, 1, 2,
		0, 2, 3,
	}

	d.shape.SetIndices(indices)

	d.shape.SetElementCount(len(indices))
}

// Draw renders shape
func (d *DynamicTextureNode) Draw(model api.IMatrix4) {
	renG := d.World().UseRenderGraphic(api.TextureRenderGraphic)

	renG.SetColor4(d.color)

	idx := d.textureIndexes[d.index]
	coords := d.textureMan.GetSTCoords(0, idx)

	renG.UpdateTexture(coords)

	// Render texture on quad
	renG.Render(d.shape, model)
}

// SetIndexes defines the indices into the texture atlas
func (d *DynamicTextureNode) SetIndexes(indexes []int) {
	d.textureIndexes = indexes
}

// SetColor ...
func (d *DynamicTextureNode) SetColor(colr []float32) {
	d.color = colr
}

// SelectCoordsByIndex is called after a render graphic has been configured
func (d *DynamicTextureNode) SelectCoordsByIndex(index int) {
	d.index = index
	// Fetch s,t texture coords from texture atlas
	idx := d.textureIndexes[index]
	coords := d.textureMan.GetSTCoords(0, idx)

	// Call VBO's update.
	renG := d.World().GetRenderGraphic(api.TextureRenderGraphic)

	renG.UpdateTexture(coords)
}
