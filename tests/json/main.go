package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wdevore/Ranger-Go-IGE/engine/configuration"
)

func main() {
	var eConf configuration.Properties

	eConfFile, err := os.Open("../../engine/config.json")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer eConfFile.Close()
	bytes, err := ioutil.ReadAll(eConfFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	err = json.Unmarshal(bytes, &eConf)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	fmt.Println(eConf.Window)

	appConfFile, err2 := os.Open("config.json")
	if err2 != nil {
		log.Fatalln("ERROR:", err2)
	}

	defer appConfFile.Close()
	bytes2, err2 := ioutil.ReadAll(appConfFile)
	if err2 != nil {
		log.Fatalln("ERROR:", err2)
	}

	err2 = json.Unmarshal(bytes2, &eConf)
	if err2 != nil {
		log.Fatalln("ERROR:", err2)
	}

	fmt.Println(eConf.Window)
}
