package api

// ISingleTextureAtlasX is for single texture atlases.
type ISingleTextureAtlasX interface {
	SelectCoordsByIndex(index int)
	SpriteSheet() ISpriteSheet
}
