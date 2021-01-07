package epaper_test

import (
	"bytes"
	"testing"

	"github.com/otaviokr/go-epaper-lib"
	"periph.io/x/periph/conn/gpio"
)

func TestWrite(t *testing.T) {
	expectedDisplayResult := []byte{
		0x13, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xfc, 0xff, 0xf0,
		0xff, 0xe1, 0xff, 0xe7, 0xff, 0xc7, 0xff, 0x8f, 0xff, 0x9f, 0xff, 0x1f, 0xff, 0x1f, 0xff, 0xbf, 0xff, 0xfe,
		0xff, 0xfe, 0xff, 0xfe, 0xff, 0x12,
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

	// Write a sentence.
	display := e.Write("Hi", 8, "examples/data/font.ttf")
	e.AddLayer(display, 0, 0, false)

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Print the image on "screen".
	e.PrintDisplay()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	t.Log(debug.Bytes())

	errorMsg = validateByteSlice(debug.Bytes(), expectedDisplayResult, "PrintDisplay function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()
}

func TestWriteRotate(t *testing.T) {
	expectedDisplayResult := []byte{
		0x13, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe3, 0xff, 0xc3, 0xff, 0x93, 0xff, 0xb2, 0xff, 0x72, 0xff, 0xf2,
		0xff, 0xe2, 0xff, 0xe2, 0xff, 0x80, 0xff, 0xc6, 0xff, 0xe6, 0xff, 0xe7, 0xff, 0xe7, 0xff, 0xef, 0xff, 0xff,
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

	// Write a sentence.
	display := e.WriteRotate("Hi", 4, "examples/data/font.ttf", true)
	e.AddLayer(display, 0, 0, false)

	// Forcing the BUSY to High to avoid being blocked because of WaitUntilIdle().
	// Do not do this on real cases!
	e.Busy.Out(gpio.High)

	// Print the image on "screen".
	e.PrintDisplay()

	// Resetting BUSY...
	e.Busy.Out(gpio.Low)

	t.Log(debug.Bytes())

	errorMsg = validateByteSlice(debug.Bytes(), expectedDisplayResult, "PrintDisplay function")
	if len(errorMsg) > 0 {
		t.Fatal(errorMsg)
	}
	debug.Reset()
}