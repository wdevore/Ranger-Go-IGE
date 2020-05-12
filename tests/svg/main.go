package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/rustyoz/svg"
)

// Some examples:
// https://github.com/fyne-io/fyne/blob/master/theme/svg.go

type Path struct {
	Id string `xml:"id,attr"`
	D  string `xml:"d, attr"`
}

type Rect struct {
	Id string `xml:"id,attr"`
}

type Group struct {
	Id          string
	Stroke      string
	StrokeWidth int32
	Fill        string
	FillRule    string
	Elements    []interface{}
}

type Svg struct {
	Title  string  `xml:"title"`
	Groups []Group `xml:"g"`
}

// Implements encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.Id = attr.Value
		case "stroke":
			g.Stroke = attr.Value
		case "stroke-width":
			if intValue, err := strconv.ParseInt(attr.Value, 10, 32); err != nil {
				return err
			} else {
				g.StrokeWidth = int32(intValue)
			}
		case "fill":
			g.Fill = attr.Value
		case "fill-rule":
			g.FillRule = attr.Value
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch tok := token.(type) {
		case xml.StartElement:
			var elementStruct interface{}

			switch tok.Name.Local {
			case "rect":
				elementStruct = &Rect{}
			case "path":
				elementStruct = &Path{}
			}

			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return err
			} else {
				g.Elements = append(g.Elements, elementStruct)
			}

			fmt.Println(tok.Name)

		case xml.EndElement:
			return nil
		}
	}
}

func ParseSvg(bytes []byte) *Svg {
	svg := new(Svg)

	err := xml.Unmarshal(bytes, svg)
	if err != nil {
		fmt.Printf("ParseSvg Error: %v\n", err)
		return nil
	}
	return svg
}

func main() {
	svgFile, err := os.Open("../../assets/BasicFont.svg")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer svgFile.Close()

	bytes, err := ioutil.ReadAll(svgFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	svg := ParseSvg(bytes)
	for _, elem := range svg.Groups[0].Elements {
		fmt.Println("Group element: ", elem)
	}
}

func test1() {
	// https://godoc.org/github.com/rustyoz/svg#ParseSvgFromReader
	svgFile, err := os.Open("../../assets/BasicFont.svg")
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	defer svgFile.Close()

	svgP, err := svg.ParseSvgFromReader(svgFile, "BasicFont", 1.0)

	if err != nil {
		log.Fatal(err)
	}

	isType, ok := svgP.Groups[0].Elements[0].(*svg.Path)

	if ok {
		fmt.Println(isType)
	}
	fmt.Println(svgP.Groups[0].Elements[0])
}
