package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/wdevore/Ranger-Go-IGE/assets/images"
)

func writeAsFlatFile(textureManifest, textureAtlasFile string, width, height int, blocks []*Block) {
	f, _ := os.Create(textureManifest)
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString(textureAtlasFile + "\n")
	w.WriteString(fmt.Sprintf("%dx%d\n", width, height))

	for _, block := range blocks {
		if block.fit != nil {
			fit := block.fit
			// Write the 4 corners of the block:
			// bottom-left, bottom-right, top-right, top-left

			//       texture-space (ST "not UV")
			//     D (0,1)        (1,1) C
			//  ^      *-----------*
			//  |      |           |
			//  |+Y    |     _     |
			//  |      |           |
			//  |      |           |
			//         *-----------*
			//     A (0,0)        (1,0) B

			botlefX := fit.x // A
			botlefY := fit.y
			botrigX := fit.x + block.w // B
			botrigY := fit.y
			toprigX := fit.x + block.w // C
			toprigY := fit.y + block.h
			toplefX := fit.x // D
			toplefY := fit.y + block.h

			w.WriteString(fmt.Sprintf("%s|%d,%d:%d,%d:%d,%d:%d,%d\n",
				block.name,
				botlefX, botlefY,
				botrigX, botrigY,
				toprigX, toprigY,
				toplefX, toplefY,
			))
		}
	}
	w.Flush()
}

func writeAsJSONFile(textureManifest, textureAtlasFile string, width, height int, blocks []*Block) {
	jOut := images.TextureManifestJSON{}
	jOut.Width = width
	jOut.Height = height
	jOut.OutputPNG = textureAtlasFile

	tiles := []*images.TextureTileJSON{}

	for _, block := range blocks {
		if block.fit != nil {
			fit := block.fit

			botlefX := fit.x // A
			botlefY := fit.y
			botrigX := fit.x + block.w // B
			botrigY := fit.y
			toprigX := fit.x + block.w // C
			toprigY := fit.y + block.h
			toplefX := fit.x // D
			toplefY := fit.y + block.h

			tile := &images.TextureTileJSON{
				Name: block.name,
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

			tiles = append(tiles, tile)
		}
	}
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
