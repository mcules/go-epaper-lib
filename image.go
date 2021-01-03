package epaper

import (
	"bytes"
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/transform"
)

// Convert the input image into a ready-to-deisplay byte buffer.
func (e *EPaper) Convert(img image.Image) []byte {
	var clearBackground byte = 0x00

	// Processing each line from the original image. If image is too large, we'll cap to the screen size.
	height := img.Bounds().Dy()
	if img.Bounds().Dy() > e.model.Height {
		height = e.model.Height
	}
	width := img.Bounds().Dx()
	if img.Bounds().Dx() > e.model.Width {
		width = e.model.Width
	}

	// Create the output array (each element represents 8 pixels, so we need a smaller array than the original matrix.)
	buffer := bytes.Repeat([]byte{0xFF}, e.lineWidth * e.model.Height)
	offset := 0
	var newValue byte = clearBackground
	for j := 0; j < height; j++ {
		for i:= 0; i < width; i++ {
			// Shift previous values before calculating the current one.
			newValue = newValue << 1

			// If color in pixel (x,y) is black, we mark it on the correct bit in the new element for the array.
			if color.Palette([]color.Color{color.Black, color.White}).Index(img.At(i, j)) == 1 {
				newValue |= 0x01
			}

			// A new byte is ready to be appended to buffer array.
			if i > 0 && (i + 1) % 8 == 0 {
				buffer[(((i +1) / 8) - 1) + (j * width) + offset] = newValue
				newValue = clearBackground
			}
		}
		offset += e.lineWidth - width
	}

	return buffer
}

// Rotate will rotate the image 90 degrees clockwise. Use it before calling convert, because convert will insert the image in the display representation matrix.
func (e *EPaper) Rotate(img image.Image) image.Image {
	return transform.Rotate(img, 90.0, &transform.RotationOptions{ResizeBounds: true, Pivot: &image.Point{0, 0}})
}
