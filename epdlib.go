package epaper

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

// Model contains the definitions of the device being used.
type Model struct {
	Width int
	Height int
	// TODO Color?
}

// EPaper represents the e-papaer device.
type EPaper struct {
	connection conn.Conn
	DataCommandSelection gpio.PinOut	// High: Data, Low: Command
	ChipSelection gpio.PinOut			// Low: active
	rst gpio.PinOut						// Low: active
	busy gpio.PinIO						// Low: active
	model Model
}

const (
	// Model7in5Width is the width size of the device display.
	Model7in5Width int = 640
	// Model7in5Height is the height of the device display.
	Model7in5Height int = 384

	// Model2in7Width is the width size of the device display.
	Model2in7Width int = 176
	// Model2in7Height is the height of the device display.
	Model2in7Height int = 264

	// ResetPin is the default pin where RST pin is connected.
	ResetPin string = "17"

	// DataCommandPin is the default pin where DC pin is connected.
	DataCommandPin string = "25"

	// ChipSelectionPin is the default pin where CS pin is connected.
	ChipSelectionPin string = "8"

	// BusyPin is the default pin where BUSY pin is connected.
	BusyPin string = "24"
)

var (
	// Model2in7bw represents the black-and-white EPD 2.7 inches display
	Model2in7bw = Model{Width: Model2in7Width, Height: Model2in7Height}

	// Model7in5 represents the EPD 7.5 inches display
	Model7in5 = Model{Width: Model7in5Width, Height: Model7in5Height}
)

const (
	// Source: https://www.waveshare.com/w/upload/2/2d/2.7inch-e-paper-Specification.pdf

	// CmdBoosterSoftStart is used during initialization.
	CmdBoosterSoftStart byte = 0x06

	// CmdDataStartTransimission1 TODO ?
	CmdDataStartTransimission1 byte = 0x10

	// CMdDeepSleep puts the screen on a low-power consumption. This should be done when the screen is not expected to be updated for a long time.
	CMdDeepSleep byte = 0x07

	// CmdDisplayRefresh TODO ?
	CmdDisplayRefresh byte = 0x12

	// CmdGetStatus = 0x71 // Used to check whne BUSY pin goes to low. Periph already takes care of this.

	// CmdLutForVcom sets the LUT for VCOM.
	CmdLutForVcom                   byte = 0x20

	// CmdLutBlue sets the LUT for White-to-White.
	CmdLutBlue                       byte = 0x21

	// CmdLutWhite sets the LUT for Black-to-White.
	CmdLutWhite                      byte = 0x22

	// CmdLutGray1 sets the LUT for White-to-Black.
	CmdLutGray1                     byte = 0x23

	// CmdLutGray2 sets the LUT for Black-to-Black.
	CmdLutGray2                     byte = 0x24

	// CmdLutRed0 ??
	CmdLutRed0                      byte = 0x25

	// CmdLutRed1 ??
	CmdLutRed1                      byte = 0x26

	// CmdLutRed2 ??
	CmdLutRed2                      byte = 0x27

	// CmdLutRed3 ??
	CmdLutRed3                      byte = 0x28

	//LUT_XON                        byte = 0x29

	// CmdPanelSetting is the code for PSR command.
	CmdPanelSetting byte = 0x00

	// CmdPartialDisplayRefresh is the code for PDRF command.
	CmdPartialDisplayRefresh byte = 0x16

	// CmdPllControl is the code for PLL command.
	CmdPllControl byte = 0x30

	// CmdPowerOff is the code for POF command.
	CmdPowerOff byte = 0x02

	// CmdPowerOn is the code for PON command.
	CmdPowerOn byte = 0x04

	// CmdPowerOptimization is the code for power optimization (not a documented command).
	CmdPowerOptimization = 0xf8

	// CmdPowerSetting is the code for PSR command.
	CmdPowerSetting byte = 0x01

	// CmdTconResolution is the code for TRES command.
	CmdTconResolution byte = 0x61

	// CmdTconSetting is the code for TCON command.
	CmdTconSetting byte = 0x60

	// CmdTemperatureCalibration is the code for TSE command.
	CmdTemperatureCalibration byte = 0x41

	// CmdVcmDcSetting is the code for VDCS command.
	CmdVcmDcSetting byte = 0x82

	// CmdVcomDataIntervalSet is the code for CDI command.
	CmdVcomDataIntervalSet byte = 0x50
)

/**
PANEL_SETTING                  byte = 0x00
	POWER_SETTING                  byte = 0x01
	POWER_OFF                      byte = 0x02
	POWER_OFF_SEQUENCE_SETTING     byte = 0x03
	POWER_ON                       byte = 0x04
	POWER_ON_MEASURE               byte = 0x05
	BOOSTER_SOFT_START             byte = 0x06
	DEEP_SLEEP                     byte = 0x07
	DATA_START_TRANSMISSION_1      byte = 0x10
	DATA_STOP                      byte = 0x11
	DISPLAY_REFRESH                byte = 0x12
	IMAGE_PROCESS                  byte = 0x13
	LUT_FOR_VCOM                   byte = 0x20
	LUT_BLUE                       byte = 0x21
	LUT_WHITE                      byte = 0x22
	LUT_GRAY_1                     byte = 0x23
	LUT_GRAY_2                     byte = 0x24
	LUT_RED_0                      byte = 0x25
	LUT_RED_1                      byte = 0x26
	LUT_RED_2                      byte = 0x27
	LUT_RED_3                      byte = 0x28
	LUT_XON                        byte = 0x29
	PLL_CONTROL                    byte = 0x30
	TEMPERATURE_SENSOR_COMMAND     byte = 0x40
	TEMPERATURE_CALIBRATION        byte = 0x41
	TEMPERATURE_SENSOR_WRITE       byte = 0x42
	TEMPERATURE_SENSOR_READ        byte = 0x43
	VCOM_AND_DATA_INTERVAL_SETTING byte = 0x50
	LOW_POWER_DETECTION            byte = 0x51
	TCON_SETTING                   byte = 0x60
	TCON_RESOLUTION                byte = 0x61
	SPI_FLASH_CONTROL              byte = 0x65
	REVISION                       byte = 0x70
	GET_STATUS                     byte = 0x71
	AUTO_MEASUREMENT_VCOM          byte = 0x80
	READ_VCOM_VALUE                byte = 0x81
	VCM_DC_SETTING                 byte = 0x82
*/

// New creates a new instance of EPaper with default parameters.
func New(model Model) (*EPaper, error) {
	return NewCustom(DataCommandPin, ChipSelectionPin, ResetPin, BusyPin, model)
}

// NewCustom creates a new instance of EPaper with custom parameters.
func NewCustom(dcPin, csPin, rstPin, busyPin string, model Model) (*EPaper, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// DC Pin
	dc := gpioreg.ByName(dcPin)
	if dc == nil {
		return nil, errors.New("spi: failed to find DC pin")
	} else if dc == gpio.INVALID {
		return nil, errors.New("epaper: use nil for dc to use 3-wire mode, do not use gpio.INVALID")
	} else if err := dc.Out(gpio.Low); err != nil {
		return nil, err
	}

	// CS Pin
	cs := gpioreg.ByName(csPin)
	if cs == nil {
		return nil, errors.New("spi: failed to find CS pin")
	} else if err := cs.Out(gpio.Low); err != nil {
		return nil, err
	}

	// RST Pin
	rst := gpioreg.ByName(rstPin)
	if rst == nil {
		return nil, errors.New("spi: failed to find RST pin")
	} else if err := rst.Out(gpio.Low); err != nil {
		return nil, err
	}

	// BUSY Pin
	busy := gpioreg.ByName(busyPin)
	if busy == nil {
		return nil, errors.New("spi: failed to find BUSY pin")
	} else if err := busy.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		return nil, err
	}

	// SPI
	port, err := spireg.Open("")
	if err != nil {
		return nil, err
	}

	// TODO official python lib limits to 4 MHz
	connection, err := port.Connect(5 * physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		port.Close()
		return nil, err
	}

	// var width, height int

	// if model % 8 == 0 {
	// 	width = (EPaperWidth / 8)
	// } else {
	// 	width = (EPaperWidth / 8 + 1)
	// }

	// height = EPaperHeight

	e := &EPaper{
		connection: connection,
		DataCommandSelection: dc,
		ChipSelection: cs,
		rst: rst,
		busy: busy,
		// width: width,
		// height: height,
		model: model,
	}

	return e, nil
}

// Reset clear the display (it can also awaken the device).
func (e *EPaper) Reset() {
	level := gpio.High
	for i := 0; i < 3; i++ {
		e.rst.Out(level)
		time.Sleep(200 * time.Millisecond)
		level = !level
	}
}

func (e *EPaper) sendCommand(c byte) {
	e.DataCommandSelection.Out(gpio.Low)
	e.ChipSelection.Out(gpio.Low)
	e.connection.Tx([]byte{c}, nil)
	e.ChipSelection.Out(gpio.High)
}

func (e *EPaper) sendData(d byte) {
	e.DataCommandSelection.Out(gpio.High)
	e.ChipSelection.Out(gpio.Low)
	e.connection.Tx([]byte{d}, nil)
	e.ChipSelection.Out(gpio.High)
}

func (e *EPaper) waitUntilIdle() {
	for e.busy.Read() == gpio.Low {
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *EPaper) turnOnDisplay() {
	e.sendCommand(CmdDisplayRefresh)
	time.Sleep(100 * time.Millisecond)
	e.waitUntilIdle()
}

// Init initializes the display config.
// It should be only used when you put the device to sleep and need to re-init the device.
func (e *EPaper) Init() {
	e.Reset()

	// ORIGINAL: e.send(CmdPowerSetting, []byte{0x37, 0x00})
	e.send(CmdPowerSetting, []byte{0x03, 0x00, 0x2b, 0x2b, 0x09})
	// e.send(CmdPanelSetting, []byte{0xcf, 0x08})
	// ORIGINAL: e.send(CmdBoosterSoftStart, []byte{0xc7, 0xcc, 0x28})
	e.send(CmdBoosterSoftStart, []byte{0x07, 0x07, 0x17})

	// Power optimizations (new)
	e.send(CmdPowerOptimization, []byte{0x60, 0xa5})
	e.send(CmdPowerOptimization, []byte{0x89, 0xa5})
	e.send(CmdPowerOptimization, []byte{0x90, 0x00})
	e.send(CmdPowerOptimization, []byte{0x93, 0x2a})
	e.send(CmdPowerOptimization, []byte{0xa0, 0xa5})
	e.send(CmdPowerOptimization, []byte{0xa1, 0x00})
	e.send(CmdPowerOptimization, []byte{0x73, 0x41})

	e.send(CmdPartialDisplayRefresh, []byte{0x00})

	e.send(CmdPowerOn, nil)

	e.waitUntilIdle()

	e.send(CmdPanelSetting, []byte{0xaf})

	// ORIGINAL: e.send(CmdPllControl, []byte{0x3c})
	e.send(CmdPllControl, []byte{0x3a})	// 3A 100Hz, 29 150Hz, 39 200Hz, 31 171Hz

	// e.send(CmdTemperatureCalibration, []byte{0x00})
	// e.send(CmdVcomDataIntervalSet, []byte{0x77})
	// e.send(CmdTconSetting, []byte{0x22})
	// e.send(CmdTconResolution, []byte{
	// 	byte(e.model.Width >> 8), byte(e.model.Width & 0xff),
	// 	byte(e.model.Height >> 8), byte(e.model.Height & 0xff)})
	// ORIGINAL: e.send(CmdVcmDcSetting, []byte{0x1e})
	e.send(CmdVcmDcSetting, []byte{0x12})

	// e.send(0xe5, []byte{0x03})

	e.send(CmdLutForVcom, Model2in7LutVcomDc)
	e.send(CmdLutBlue, Model2in7LutWw)
	e.send(CmdLutWhite, Model2in7LutBw)
	e.send(CmdLutGray1, Model2in7LutWb)
	e.send(CmdLutGray2, Model2in7LutBb)

}

func (e *EPaper) send(cmd byte, data []byte) {
	e.sendCommand(cmd)
	if data != nil {
		for _, d := range data {
			e.sendData(d)
		}
	}
}

// ClearScreen erases anything that is on screen.
func (e *EPaper) ClearScreen() {
	data := make([]byte, e.model.Height * e.model.Width * 4)
	for i := range data {
		data[i] = 0x33
	}

	e.send(CmdDataStartTransimission1, data)
	e.turnOnDisplay()
}

// Display takes a byte buffer and updates the screen.
func (e *EPaper) Display(img []byte) {
	// TODO You must understand!
	e.sendCommand(CmdDataStartTransimission1)

	for j := 0; j < e.model.Height; j++ {
		for i := 0; i < e.model.Width; i++ {
			dataBlack := ^img[i+j*e.model.Height]

			for k := 0; k < 8; k++ {
				var data byte

				if dataBlack&0x80 > 0 {
					data = 0x00
				} else {
					data = 0x03
				}

				data <<= 4
				dataBlack <<= 1
				k++

				if dataBlack&0x80 > 0 {
					data |= 0x00
				} else {
					data |= 0x03
				}

				dataBlack <<= 1

				e.sendData(data)
			}
		}
	}

	e.turnOnDisplay()
}

// Sleep put the display in power-saving mode.
// You can use Reset() to awaken and Init() to re-initialize the display.
func (e *EPaper) Sleep() {
	// TODO official python lib does it differently
	// self.SPI.close()
	// self.GPIO.outpu(RST_PIN, 0)
	// self.GPIO.output(DC_PIN, 0)
	// self.GPIO.cleanup()

	e.sendCommand(CmdPowerOff)
	e.waitUntilIdle()
	e.send(CMdDeepSleep, []byte{0xA5})
}

// Convert the input image into a ready-to-deisplay byte buffer.
func (e *EPaper) Convert(img image.Image) []byte {
	var byteToSend byte = 0x00
	var bgColor = 1

	buffer := bytes.Repeat([]byte{0x00}, e.model.Width * e.model.Height)

	for i := 0; i < e.model.Width; i++ {
		for j := 0; j < e.model.Height; j++ {
			bit := bgColor

			if i < img.Bounds().Dx() && j < img.Bounds().Dy() {
				bit = color.Palette([]color.Color{color.Black, color.White}).Index(img.At(i, j))
			}

			if bit == 1 {
				byteToSend |= 0x80 >> (uint32(i) % 8)
			}

			if i % 8 == 7 {
				buffer[(i / 8) + (j * e.model.Width)] = byteToSend
				byteToSend = 0x00
			}
		}
	}

	return buffer
}