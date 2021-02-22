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
func (g *group) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			g.ID = attr.Value
		case "label":
			g.Label = attr.Value
		case "transform":
			// Ex translate(-5.5548377,-1.6841292)
			attr.Value = strings.TrimRight(attr.Value, ")")
			attr.Value = strings.TrimLeft(attr.Value, "translate(")
			f := strings.Split(attr.Value, ",")
			v, err := strconv.ParseFloat(f[0], 32)
			if err != nil {
				fmt.Println(f[0])
				return err
			}
			g.TX = float32(v)

			v, err = strconv.ParseFloat(f[1], 32)
			if err != nil {
				fmt.Println(f[0])
				return err
			}
			g.TY = float32(v)
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
				elementStruct = &rectangle{}
			case "text":
				elementStruct = &text{}
			case "path":
				elementStruct = &path{}
			case "g":
				elementStruct = &group{}
			default:
				fmt.Println("Default tok.Name: ", tok.Name)
				continue
			}

			if err = decoder.DecodeElement(elementStruct, &tok); err != nil {
				return err
			}

			g.Elements = append(g.Elements, elementStruct)

		case xml.EndElement:
			return nil
		}
	}
}
