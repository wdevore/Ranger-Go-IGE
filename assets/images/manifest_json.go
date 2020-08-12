package images

// TextureTileJSON is a cell with a name
type TextureTileJSON struct {
	Name     string
	XYCoords []int
	STCoords []float32
}

// TextureManifestJSON main manifest header
type TextureManifestJSON struct {
	OutputPNG string
	Width     int
	Height    int
	Layer     float32
	Tiles     []*TextureTileJSON
}
