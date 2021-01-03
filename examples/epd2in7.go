package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/otaviokr/go-epaper-lib"
)

func main() {
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
	time.Sleep(10000 * time.Millisecond)
	fmt.Println(" done.")

	fmt.Printf("Cleaning screen (flashes are expected)...")
	epd.ClearScreen()
	fmt.Println(" done.")

	// Print image.
	fmt.Printf("Printing screen 2 (Simple image) for 10s...")
	printImage(epd, "data/demo.png")
	time.Sleep(10000 * time.Millisecond)
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

	// Print image.
	fmt.Printf("Printing screen 4 (Simple image, rotated) for 10s...")
	printImageRotated(epd, "data/demo.png")
	time.Sleep(10000 * time.Millisecond)
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
	b := epd.Convert(m)
	epd.Display(b)
}

func printTextRotated(epd *epaper.EPaper, text string, fontSize float64, fontFile string) {
	m := epd.WriteRotate(text, fontSize, fontFile, true)
	r := epd.Rotate(m)
	b := epd.Convert(r)
	epd.Display(b)
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

	b := epd.Convert(m)
	epd.Display(b)
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

	x := epd.Rotate(m)

	b := epd.Convert(x)
	epd.Display(b)
}
