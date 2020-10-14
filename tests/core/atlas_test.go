package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/wdevore/Ranger-Go-IGE/engine"
	"github.com/wdevore/Ranger-Go-IGE/engine/rendering/atlas"
	"github.com/wdevore/Ranger-Go-IGE/extras/generators"
)

// go test -v -count=1 quadtree_test.go

func TestRunner(t *testing.T) {
	testMonoAtlasAddShape(t)
}

func testMonoAtlasAddShape(t *testing.T) {
	engine, err := engine.Construct("../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	mono := atlas.NewStaticMonoAtlas()
	err = mono.Configure("../../", world.Properties())
	if err != nil {
		panic(err)
	}

	vertices, indices, mode := generators.GenerateUnitHLineVectorShape()
	id := mono.AddShape("HLine", vertices, indices, mode)
	fmt.Println(id)

	vertices, indices, mode = generators.GenerateUnitVLineVectorShape()
	id = mono.AddShape("VLine", vertices, indices, mode)
	fmt.Println(id)

	vertices, indices, mode = generators.GenerateUnitPlusVectorShape()
	id = mono.AddShape("Plus", vertices, indices, mode)
	fmt.Println(id)

	mono.Shake()
}

func testMonoAtlasConfig(t *testing.T) {
	engine, err := engine.Construct("../..", "config.json")
	if err != nil {
		log.Fatal(err)
	}

	defer engine.End()

	world := engine.World()

	mono := atlas.NewStaticMonoAtlas()
	err = mono.Configure("../../", world.Properties())
	if err != nil {
		panic(err)
	}
}
