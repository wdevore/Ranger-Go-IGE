package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rustyoz/svg"
)

// ####################################################
// This is an old prelimary test. See the tools folder
// for a more feature complete code base.
// ####################################################

// Some examples:
// https://github.com/fyne-io/fyne/blob/master/theme/svg.go

type PathVertex struct {
	x, y float32
}

type Text struct {
	ID    string  `xml:"id,attr"`
	X     float32 `xml:"x,attr"`
	Y     float32 `xml:"y,attr"`
	Style string  `xml:"style,attr"`
}

type Path struct {
	ID       string `xml:"id,attr"`
	D        string `xml:"d,attr"`
	Style    string `xml:"style,attr"`
	Vertices []PathVertex
	Closed   bool // True = LINE_LOOP, False = LINE_STRIP
}

type Rectangle struct {
	ID     string  `xml:"id,attr"`
	X      float32 `xml:"x,attr"`
	Y      float32 `xml:"y,attr"`
	Height float32 `xml:"height,attr"`
	Width  float32 `xml:"width,attr"`
}

type Group struct {
	ID       string
	Label    string
	Elements []interface{} // <path>, id ...
}

type Svg struct {
	Title  string  `xml:"title"`
	Groups []Group `xml:"g"`
}

// UnmarshalXML Implements encoding.xml.Unmarshaler interface
func (p *Path) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			p.ID = attr.Value
		case "d":
			p.parsePath(attr.Value)
			p.D = attr.Value
		case "style":
			p.Style = attr.Value
		default:
			// fmt.Println("Path UnmarshalXML: ", attr.Value)
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.StartElement:
		case xml.EndElement:
			return nil
		}
	}
}

func (p *Path) parsePath(path string) error {
	// Example:
	// d="M 0.65009511,8.0172702 3.4857168,0.29375086 6.6257727,8.0719752"
	// OR relative
	// d="m 23.643773,17.674347 -0.0076,-7.117348 h -0.561806 l -3.052303,4.661477 -0.03616,0.463303 4.970314,0.01889"
	// OR closed
	// d="m 188.0518,2.5272938 1.39562,0.4586884 0.71797,1.3376465 ... z"

	// If "m" then the trailing vertices are relative.
	// If "M" then all vertices are absolute.
	// If M...z then it is a closed path
	// Each time we hit a capital command the current-point is reset to
	// the coord, otherwise C-P changes by a delta.

	p.Vertices = []PathVertex{}

	// if strings.Contains(path, "h ") {
	// 	fmt.Println("stop")
	// }
	tokens := strings.Split(path, " ")
	command := tokens[0]
	relative := command == "m"

	// The first vertex is always absolute
	coord := strings.Split(tokens[1], ",")

	// x,y are the current point
	x, err := strconv.ParseFloat(coord[0], 32)
	if err != nil {
		fmt.Println(x)
		return err
	}
	y, err := strconv.ParseFloat(coord[1], 32)
	if err != nil {
		return err
	}

	p.Vertices = append(p.Vertices, PathVertex{x: float32(x), y: float32(y)})

	// Begin scanning the rest of the path.
	for i := 2; i < len(tokens); i++ {
		// s[i] is either a coordinate or command
		token := tokens[i]

		// Check if command.
		switch token {
		case "z", "Z":
			// It is closed a closed path, we mark it and finish.
			p.Closed = true
			return nil
		case "l", "v", "h":
			// The next coord is relative to current-point.
			relative = true
			command = token
			continue // Move to coord.
		case "L", "V", "H":
			// The next coord is absolute.
			relative = false
			command = token
			continue // Move to absolute coord.
			// default:
			// 	return errors.New("Unknown command in path: " + token)
		}

		coordX, coordY := p.getCoord(command, token, x, y)

		command = ""

		if relative {
			x += coordX
			y += coordY
		} else {
			x = coordX
			y = coordY
		}

		p.Vertices = append(p.Vertices, PathVertex{x: float32(x), y: float32(y)})
	}

	return nil
}

func (p *Path) getCoord(command, coord string, cx, cy float64) (x, y float64) {

	switch command {
	case "v":
		x = 0.0
		y, _ = strconv.ParseFloat(coord, 64)
	case "V":
		x = cx
		y, _ = strconv.ParseFloat(coord, 64)
	case "h":
		x, _ = strconv.ParseFloat(coord, 64)
		y = 0.0
	case "H":
		x, _ = strconv.ParseFloat(coord, 64)
		y = cy
	default:
		s := strings.Split(coord, ",")
		x, _ = strconv.ParseFloat(s[0], 64)
		y, _ = strconv.ParseFloat(s[1], 64)
	}

	return x, y
}

// UnmarshalXML Implements encoding.xml.Unmarshaler interface
func (g *Group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.ID = attr.Value
		case "label":
			g.Label = attr.Value
		default:
			// if attr.Value == "A" {
			// 	fmt.Println("Group UnmarshalXML: ", attr.Value)
			// }
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
				elementStruct = &Rectangle{}
			case "text":
				elementStruct = &Text{}
			case "path":
				elementStruct = &Path{}
			case "g":
				elementStruct = &Group{}
			default:
				fmt.Println("Default tok.Name: ", tok.Name)
				continue
			}

			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return err
			} else {
				g.Elements = append(g.Elements, elementStruct)
			}

		case xml.EndElement:
			return nil
		}
	}
}

func ParseSvg(bytes []byte) (*Svg, error) {
	svg := new(Svg)

	err := xml.Unmarshal(bytes, svg)
	if err != nil {
		return nil, err
	}
	return svg, nil
}

func main() {
	test2()
}

func test2() {
	// BasicFont.svg
	svgFile, err := os.Open("../../assets/CascadiaCode.svg")
	if err != nil {
		panic(err)
	}

	defer svgFile.Close()

	bytes, err := ioutil.ReadAll(svgFile)
	if err != nil {
		panic(err)
	}

	svg, err := ParseSvg(bytes)
	if err != nil {
		panic(err)
	}

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
