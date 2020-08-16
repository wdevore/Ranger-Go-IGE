package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/wdevore/Ranger-Go-IGE/assets/images"
)

func writeAsJSONFile(textureManifest, textureAtlasFile string, width, height int) {
	jOut := images.TextureManifestJSON{}
	jOut.Width = width
	jOut.Height = height
	jOut.Layer = 0.0
	jOut.OutputPNG = textureAtlasFile

	tiles := []*images.TextureTileJSON{}

	cnt := 0
	cw := 9
	ch := 9

	x := 0
	y := 0
	c := 97
	for rc := 0; rc < 14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	c = 65
	for rc := 0; rc < 14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	for rc := 0; rc < 26-14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	c = 33
	for rc := 0; rc < 14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	for rc := 0; rc < 14; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	x = 0
	y += 9
	for rc := 0; rc < 4; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	c = 91
	for rc := 0; rc < 6; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	c = 123
	for rc := 0; rc < 4; rc++ {
		genTiles(x, y, c, cw, ch, float32(width), float32(height), &tiles)
		x += 9
		c++
		cnt++
	}

	fmt.Println("Rune Count: ", cnt)

	jOut.Count = cnt
	jOut.Tiles = tiles

	indentedJSON, _ := json.MarshalIndent(jOut, "", "  ")

	dataPath, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON written to: ", dataPath+textureManifest)
	err = ioutil.WriteFile(dataPath+"/"+textureManifest, indentedJSON, 0644)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}

func genTiles(x, y, c, cw, ch int, width, height float32, tiles *[]*images.TextureTileJSON) {
	botlefX := x // A
	botlefY := y
	botrigX := x + cw // B
	botrigY := y
	toprigX := x + cw // C
	toprigY := y + ch
	toplefX := x // D
	toplefY := y + ch

	tile := &images.TextureTileJSON{
		Name: string(rune(c)),
		XYCoords: []int{
			botlefX,
			botlefY,
			botrigX,
			botrigY,
			toprigX,
			toprigY,
			toplefX,
			toplefY,
		},
		STCoords: []float32{
			float32(botlefX) / float32(width),
			1.0 - (float32(botlefY) / float32(height)),
			float32(botrigX) / float32(width),
			1.0 - (float32(botrigY) / float32(height)),
			float32(toprigX) / float32(width),
			1.0 - (float32(toprigY) / float32(height)),
			float32(toplefX) / float32(width),
			1.0 - (float32(toplefY) / float32(height)),
		},
	}

	*tiles = append(*tiles, tile)
}
