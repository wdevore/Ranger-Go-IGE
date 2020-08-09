package textures

// It requires a manifest doc that describes the sub-texture layouts
// For example:
//
// texture-atlas.png
// 256x256
// mine|0,64:64,64:64,128:0,128
// green ship|64,192:128,192:128,256:64,256
// orange ship|192,192:256,192:256,256:192,256
// ctype ship|0,192:64,192:64,256:0,256
// bomb|192,160:208,160:208,176:192,176
//

// TextureCoord is:
// uv coordinates start from the upper left corner (v-axis is facing down).
// st coordinates start from the lower left corner (t-axis is facing up).
// s = u;
// t = 1-v;
type textureCoord struct {
	s, t float32
}

func (tx *textureCoord) S() float32 {
	return tx.s
}

func (tx *textureCoord) T() float32 {
	return tx.t
}

func (tx *textureCoord) ST() (float32, float32) {
	return tx.s, tx.t
}

func (tx *textureCoord) Set(s, t float32) {
	tx.s = s
	tx.t = t
}
