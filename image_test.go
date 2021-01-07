package epaper_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	"github.com/anthonynsimon/bild/paint"
	"github.com/otaviokr/go-epaper-lib"
	"periph.io/x/periph/conn/gpio"
)

func TestAddLayerNoTansparency(t *testing.T) {
	expectedDisplayResult := []byte{
		0x13, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x1e, 0xff, 0x1e, 0xff, 0x1e, 0xff, 0x1e, 0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x12,
	}

	// Create a dummy "epaper"
	// (to create a real one, use the example source code, this won't work!)
	debug := new(bytes.Buffer)
	e, err := epaper.NewCustom("", "", "", "", ModelSim, true, debug)
	if err != nil {
		t.Fatal(err)
	}

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Start up the epaper display and monitor the output..
	e.Init()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	// Validating data sent to dummy device.
	errorMsg := validateByteSlice(debug.Bytes(), ExpectedInitResult, "Input function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Clearing the screen - this is default process.
	e.ClearScreen()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), ExpectedClearScreenResult, "ClearScreen function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Create a black square.
	dest := paint.FloodFill(image.Rect(0, 0, 10, 10), image.Point{X: 0, Y: 0}, color.RGBA{A: 255}, 0)
	e.AddLayer(dest, 0, 0, false)

	// Put a smaller white square on the center.
	mask := paint.FloodFill(image.Rect(0, 0, 10, 10), image.Point{X: 0, Y: 0}, color.RGBA{A: 255}, 0)
	for i := 3; i < 7; i++ {
		for j := 3; j < 7; j++ {
			mask.SetRGBA(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	e.AddLayer(mask, 0, 0, false)

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Print the image on "screen".
	e.PrintDisplay()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), expectedDisplayResult, "PrintDisplay function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()
}

func TestAddLayerWithTansparency(t *testing.T) {
	expectedDisplayResult := []byte{
		0x13, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x12,
	}

	// Create a dummy "epaper"
	// (to create a real one, use the example source code, this won't work!)
	debug := new(bytes.Buffer)
	e, err := epaper.NewCustom("", "", "", "", ModelSim, true, debug)
	if err != nil {
		t.Fatal(err)
	}

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Start up the epaper display and monitor the output..
	e.Init()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	// Validating data sent to dummy device.
	errorMsg := validateByteSlice(debug.Bytes(), ExpectedInitResult, "Input function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Clearing the screen - this is default process.
	e.ClearScreen()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), ExpectedClearScreenResult, "ClearScreen function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Create a black square.
	dest := paint.FloodFill(image.Rect(0, 0, 10, 10), image.Point{X: 0, Y: 0}, color.RGBA{A: 255}, 0)
	e.AddLayer(dest, 0, 0, false)

	// Put a white square on top. Since it is transparent (white = transparent), the square should remain black.
	mask := paint.FloodFill(image.Rect(0, 0, 10, 10), image.Point{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255}, 0)
	// for i := 3; i < 7; i++ {
	// 	for j := 3; j < 7; j++ {
	// 		mask.SetRGBA(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	// 	}
	// }
	e.AddLayer(mask, 0, 0, true)

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Print the image on "screen".
	e.PrintDisplay()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), expectedDisplayResult, "PrintDisplay function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()
}

func TestRotate(t *testing.T) {
	expectedDisplayResult := []byte{
		0x13, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f,
		0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x12,
	}

	// Create a dummy "epaper"
	// (to create a real one, use the example source code, this won't work!)
	debug := new(bytes.Buffer)
	e, err := epaper.NewCustom("", "", "", "", ModelSim, true, debug)
	if err != nil {
		t.Fatal(err)
	}

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Start up the epaper display and monitor the output..
	e.Init()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	// Validating data sent to dummy device.
	errorMsg := validateByteSlice(debug.Bytes(), ExpectedInitResult, "Input function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Clearing the screen - this is default process.
	e.ClearScreen()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), ExpectedClearScreenResult, "ClearScreen function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()

	// Create two black horizontal strips.
	dest := paint.FloodFill(image.Rect(0, 0, 16, 16), image.Point{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255}, 0)
	for i := 0; i < 15; i++ {
		for j := 0; j < 6; j++ {
			dest.SetRGBA(i, j, color.RGBA{A: 255})
		}
		for k := 12; k < 16; k++ {
			dest.SetRGBA(i, k, color.RGBA{A: 255})
		}
	}

	// Rotate the image
	final := e.Rotate(dest)
	e.AddLayer(final, 0, 0, false)

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Print the image on "screen".
	e.PrintDisplay()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	errorMsg = validateByteSlice(debug.Bytes(), expectedDisplayResult, "PrintDisplay function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()
}
