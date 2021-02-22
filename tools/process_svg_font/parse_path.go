package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// Some examples:
// https://github.com/fyne-io/fyne/blob/master/theme/svg.go

// UnmarshalXML Implements encoding.xml.Unmarshaler interface
func (p *path) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
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

func (p *path) parsePath(path string) error {
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

	p.Vertices = []*pathVertex{}

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

	p.Vertices = append(p.Vertices, &pathVertex{x: float32(x), y: float32(y)})

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
		case "c":
			return fmt.Errorf("unexpected command: %s for token: %s", command, token)
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

		p.Vertices = append(p.Vertices, &pathVertex{x: float32(x), y: float32(y)})
	}

	return nil
}

func (p *path) getCoord(command, coord string, cx, cy float64) (x, y float64) {

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
