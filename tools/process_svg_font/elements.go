package main

import "fmt"

// Some examples:
// https://github.com/fyne-io/fyne/blob/master/theme/svg.go

type pathVertex struct {
	x, y float32
}

type text struct {
	ID    string  `xml:"id,attr"`
	X     float32 `xml:"x,attr"`
	Y     float32 `xml:"y,attr"`
	Style string  `xml:"style,attr"`
}

func (t text) String() string {
	return fmt.Sprintf("{ID:%s, XY:%0.3f,%0.3f}", t.ID, t.X, t.Y)
}

// path is manually parsed using UnmarshalXML, thus no `xml:...` tags
type path struct {
	ID       string
	D        string
	Style    string
	Vertices []*pathVertex
	Closed   bool // True = LINE_LOOP, False = LINE_STRIP
}

func (p path) String() string {
	// return fmt.Sprintf("{ID:%s, Closed:%v, Verts:%v}", p.ID, p.Closed, p.Vertices)
	return fmt.Sprintf("{ID:%s, Closed:%v}", p.ID, p.Closed)
}

type rectangle struct {
	ID     string  `xml:"id,attr"`
	X      float32 `xml:"x,attr"`
	Y      float32 `xml:"y,attr"`
	Height float32 `xml:"height,attr"`
	Width  float32 `xml:"width,attr"`
}

func (r rectangle) String() string {
	return fmt.Sprintf("{ID:%s, XY:%0.3f,%0.3f}, Size:%0.3fx%0.3f", r.ID, r.X, r.Y, r.Width, r.Height)
}

// group is manually parsed using UnmarshalXML, thus no `xml:...` tags
type group struct {
	ID       string
	Label    string
	TX, TY   float32
	Elements []interface{} // <path>, id ...
}

func (g group) String() string {
	return fmt.Sprintf("{ID:%s, Label:%s}", g.ID, g.Label)
}

// Svg --
type Svg struct {
	Title  string  `xml:"title"`
	Groups []group `xml:"g"`
}
