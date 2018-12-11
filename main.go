package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Text    string   `xml:",chardata"`
	Xmlns   string   `xml:"xmlns,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	ViewBox string   `xml:"viewBox,attr"`
	G       struct {
		Text    string `xml:",chardata"`
		ID      string `xml:"id,attr"`
		Style   string `xml:"style,attr"`
		Polygon []struct {
			Text          string `xml:",chardata"`
			ID            string `xml:"id,attr"`
			Points        string `xml:"points,attr"`
			Fill          string `xml:"fill,attr"`
			FillOpacity   string `xml:"fill-opacity,attr"`
			StrokeOpacity string `xml:"stroke-opacity,attr"`
		} `xml:"polygon"`
	} `xml:"g"`
}

type svgProperties struct {
	fillSegments string
	color        string
	name         string
}

var (
	red   = "#FF0000"
	green = "#00FF00"
	blue  = "#0000FF"
	cyan  = "#00FFFF"
	black = "#000000"
	white = "#FFFFFF"
	gray  = "#999999"
)


// Check segmentLetters.png to get a better grasp of the mapping
var letter2Segments = map[string]string{
	"0": "abcdef", // O uses a,b,c,d,e and f segments of 7 segment display
	"1": "ef",
	"2": "abedg",
	"3": "abcdg",
	"4": "bcfg",
	"5": "acdfg",
	"6": "acdefg",
	"7": "abc",
	"8": "abcdefg",
	"9": "abcdfg",
	"d": "bcdeg",
	"h": "cefg",
	"y": "bcdfg",
	"H": "bcefg",
	"I": "ef",
	"L": "def",
	"P": "abefg",
	"U": "bcdef",
}

// Corresponding SVG files will be created for each array element
var svgOut = []svgProperties{
	{letter2Segments["8"], black, "black8.svg"},
	{letter2Segments["6"], red, "red6.svg"},
	{letter2Segments["U"], cyan, "cyanU.svg"},
	{letter2Segments["L"], gray, "grayL.svg"},
}

func main() {
	// Do not delete templates/7seg8_v1.svg!
	// This SVG file is taken from Wikimedia.org
	dat, _ := ioutil.ReadFile("templates/7seg8_v1.svg")
	var data Svg
	xml.Unmarshal(dat, &data)

	for _, v := range svgOut {
		fillStr := v.fillSegments
		polygons := data.G.Polygon
		for i, _ := range polygons {
			polygon := &polygons[i]
			polygonID := polygon.ID
			if strings.Contains(fillStr, polygonID) {
				polygon.Fill = v.color
				polygon.FillOpacity = "1"
				polygon.StrokeOpacity = "1"
			} else {
				polygon.FillOpacity = "0"
				polygon.StrokeOpacity = "0"
			}
		}
		result, err := xml.Marshal(data)
		if err != nil {
			fmt.Errorf("svg data in struct could not be converted into string!")
		}

		file, err := os.Create(v.name)
		defer file.Close()
		if err != nil {
			fmt.Errorf("%v file cannot be created!", v.name)
		}
		file.WriteString(string(result))
	}
}
