package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/mcules/go-epaper-lib"
)

func main() {
	var waitTime time.Duration = 5000
	fmt.Println("Starting demo for 2.7\" Display...")

	// Create new EPaperDisplay handler.
	epd, err := epaper.New(epaper.Model2in7bw)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	// Run mandatory initialization.
	fmt.Printf("Initializing e-paper...")
	epd.Init()
	fmt.Println(" done.")

	// Clear screen to avoid undesired meshes.
	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print text.
	fmt.Printf("Printing screen 1 (Simple text) for 10s...")
	printText(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf")
	time.Sleep(waitTime * time.Millisecond)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print image.
	fmt.Printf("Printing screen 2 (Simple image) for 10s...")
	printImage(epd, "data/demo.png")
	time.Sleep(waitTime * time.Millisecond)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print text, rotate 90 degree clockwise.
	fmt.Printf("Printing screen 3 (Simple text, rotated)...")
	printTextRotated(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf")
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print image, rotated.
	fmt.Printf("Printing screen 4 (Simple image, rotated) for 10s...")
	printImageRotated(epd, "data/demo.png")
	time.Sleep(waitTime * time.Millisecond)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print text, changing starting point.
	fmt.Printf("Printing screen 5 (Simple text on arbitrary point on screen)...")
	printTextPosition(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf", 30, -32)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print image, rotated, on arbitrary position..
	fmt.Printf("Printing screen 6 (Simple image, rotated, on arbitrary position) for 10s...")
	printImageRotatedPosition(epd, "data/demo.png", 30, -32)
	time.Sleep(waitTime * time.Millisecond)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print text, rotated, changing starting point.
	fmt.Printf("Printing screen 7 (Simple text, rotated, on arbitrary point on screen)...")
	printTextRotatedPosition(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf", 30, -32)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Composite images, no transparency..
	fmt.Printf("Printing screen 8 (Two layers, no transparency)...")
	printTwoLayerWithTranparency(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf", 30, 30, "data/demo.png", false)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Composite images, with transparency..
	fmt.Printf("Printing screen 9 (Two layers, with transparency)...")
	printTwoLayerWithTranparency(epd, "Hi, I'm an e-paper display! This is a demo.", 8, "data/font.ttf", 30, 30, "data/demo.png", true)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	fmt.Printf("Putting display into low-power sleep...")
	epd.Sleep()
	fmt.Println(" done.")

	fmt.Println("... finished.")
}

func printText(epd *epaper.EPaper, text string, fontSize float64, fontFile string) {
	m := epd.Write(text, fontSize, fontFile)
	epd.AddLayer(m, 0, 0, false)
	epd.PrintDisplay()
}

func printTextRotated(epd *epaper.EPaper, text string, fontSize float64, fontFile string) {
	m := epd.WriteRotate(text, fontSize, fontFile, true)
	r := epd.Rotate(m)
	epd.AddLayer(r, 0, 0, false)
	epd.PrintDisplay()
}

func printTextPosition(epd *epaper.EPaper, text string, fontSize float64, fontFile string, x, y int) {
	m := epd.Write(text, fontSize, fontFile)
	epd.AddLayer(m, 30, 30, true)
	epd.PrintDisplay()
}

func printTextRotatedPosition(epd *epaper.EPaper, text string, fontSize float64, fontFile string, x, y int) {
	m := epd.WriteRotate(text, fontSize, fontFile, true)
	r := epd.Rotate(m)
	epd.AddLayer(r, x, y, false)
	epd.PrintDisplay()
}

func printImage(epd *epaper.EPaper, imageFile string) {
	reader, err := os.Open(imageFile)
	if err != nil {
		fmt.Printf("ERROR while loading image: %+v\n", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("ERROR while decoding image: %+v\n", err)
	}

	epd.AddLayer(m, 0, 0, false)
	epd.PrintDisplay()
}

func printImageRotated(epd *epaper.EPaper, imageFile string) {
	reader, err := os.Open(imageFile)
	if err != nil {
		fmt.Printf("ERROR while loading image: %+v\n", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("ERROR while decoding image: %+v\n", err)
	}

	r := epd.Rotate(m)

	epd.AddLayer(r, 0, 0, false)
	epd.PrintDisplay()
}

func printImageRotatedPosition(epd *epaper.EPaper, imageFile string, x, y int) {
	reader, err := os.Open(imageFile)
	if err != nil {
		fmt.Printf("ERROR while loading image: %+v\n", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("ERROR while decoding image: %+v\n", err)
	}

	r := epd.Rotate(m)

	epd.AddLayer(r, x, y, false)
	epd.PrintDisplay()
}

func printTwoLayerWithTranparency(epd *epaper.EPaper, text string, fontSize float64, fontFile string, x, y int, imageFile string, transparent bool) {
	reader, err := os.Open(imageFile)
	if err != nil {
		fmt.Printf("ERROR while loading image: %+v\n", err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		fmt.Printf("ERROR while decoding image: %+v\n", err)
	}

	epd.AddLayer(m, 0, 0, false)

	t := epd.Write(text, fontSize, fontFile)
	epd.AddLayer(t, x, y, transparent)
	epd.PrintDisplay()
}