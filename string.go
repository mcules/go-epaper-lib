package epaper

import (
	"bytes"
	"image"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/hqbobo/text2pic"
)

// Write is used to prepare text to be printed on the display.
// It turns a string in an image, then calls the returned value as a parameter of Convert(), and later, call Display().
func (e *EPaper) Write(text string, fontSize float64, fontFile string) (image.Image) {
	return e.WriteRotate(text, fontSize, fontFile, false)
}

// WriteRotate is used to prepare text to be printed on the display.
// It turns a string in an image, then calls the returned value as a parameter of Convert(), and later, call Display().
// If rotate is TRUE, the text will be rotate 90 degree clockwise.
func (e *EPaper) WriteRotate(text string, fontSize float64, fontFile string, rotate bool) (image.Image) {
	// Read the font file.
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		panic(err)
	}

	// Parse the font read, so it can be used correctly.
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	// Convert the text to image.
	width := e.model.Width
	if rotate {
		width = e.model.Height
	}
	pic := text2pic.NewTextPicture(text2pic.Configure{Width: width, BgColor: text2pic.ColorWhite})
	pic.AddTextLine(text, fontSize, f, text2pic.ColorBlack, text2pic.Padding{Left: 0, Top: 0, Bottom: 0})

	var buffer bytes.Buffer
	err = pic.Draw(&buffer, text2pic.TypePng)
	if err != nil {
		// FIXME Better error handling.
		panic(err)
	}
	bufferReader := bytes.NewReader(buffer.Bytes())

	img, _, err := image.Decode(bufferReader)
	if err != nil {
		// FIXME Better error handling.
		panic(err)
	}

	return img
}
