# go-epaper-lib
A golang lib to use with Waveshare(tm) e-Paper HAT for Raspberry Pi.

```
The current version is developed/tested only for the 2.7 inches, black-and-white HAT.
Hopefully, new models will be available in the future, but, at the moment,
keep in mind that all info in this repo is related exclusively to the 2.7" model.
```

# Intro for the uninitiated

This is a library to help you use Waveshare(tm)'s [e-paper](https://en.wikipedia.org/wiki/Electronic_paper) display in [Go](https://golang.org/). If you are looking for libs in other languages (like Python or C, check the official Waveshare website!).

This lib is still on its initial pre-release version. That means, it is still under development, no guarantees of retrocompatibility between commits. Still, pushed commits should be functional.

# IMPORTANT NOTICE FOR RASPBERRY PI USERS

If you want to use the e-paper with Raspberry Pi, you need to enable the SPI kernel modules.

- Run the command `sudo raspi-config` (notice that you need root access).
- Access menu `3 Interface Options`, then `P4 SPI` and finally, select `Yes`.
- Reboot the Raspberry PI

```
If you don't follow these commands, you may see an error when running the lib similar to:

spireg: no port found. Did you run Init()?
```

# The basics

This is the minimum code to have something printed on the display. For more details, check the documentation and/or the demo codes at `examples/`.

```go
package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/otaviokr/go-epaper-lib"
)

func main() {
	// New EPaperDisplay handler for model 2.7 inches, black-and-white..
	epd, err := epaper.New(epaper.Model2in7bw)
	if err != nil {
		panic(err)
	}

	// Run mandatory initialization.
	epd.Init()

	// Clear screen to avoid undesired meshes.
	epd.ClearScreen()

	// Parse the string and process it with TTF font.
  m := epd.Write(text, fontSize, fontFile)

  // Add the resulting image into the buffer.
  epd.AddLayer(m, 0, 0, false)

  // Print the buffer onto the display..
  epd.Display()

  // Just a timeout before clearing the screen.
	time.Sleep(10000 * time.Millisecond)
	epd.ClearScreen()

  // Always put the screen to sleep when it will not be refreshed shortly.
	epd.Sleep()
}
```

# Functionalities

All of these functionalities are demonstrated in the example programs at `examples/`.

- **Write text**: Write a string in the display. Tabs and newline chars are not recognized, nor bold, italic, underline etc.
- **Write text, rotated**: the text is written rotated 90 degrees clockwise.
- **Display image**: the image is cropped if it is larger than the display size. Only black-and-white PNG files allowed.
- **Display image, rotated**: the image is rotated 90 degrees clockwise and it will be cropped if larger than display.

# Coordinate system

```
  +-----------------------------------------------+
  | O ::::::::::::::::::::::::: O                 |
  |   +-----------------------------------+  K4[] |
  | 2 | (176,0)                 (176,264) |       |
  | 7 |                                   |  K3[] |
  | H |                                   |       |
  | A |                                   |  K2[] |
  | T | (0,0)                    (0, 264) |       |
  |   +-----------------------------------+  K1[] |
  +-----------------------------------------------+

  The schematics of the display above show the coordinates on the screen and orientation.
```

Each byte sent fills 8 pixels on line (X), so always group the pixels in sets of 8 before sending them to be printed. Bit value 0 means "black", while 1 means "white".

So, we need to convert an image represented as a matrix where each element is a pixel, to an array where each element represented 8 pixels on X and the lines are concatenated.

```
Image matrix:
| 0 1 2 3 4 a b c d e f g h i j 9 |
| A B C D E k l m n o p q r s t F |

becomes :
[01234abc defghij9 ABCDEklm nopqrstF]
```
If the image is larger than the display, the image is cropped.

IMPORTANT! The image is cropped during rotate AND translation, so the order of the commands are relevant!

# How the output is composed

## Image

The PNG image is read into a `image.Image` object. This object can be manipulated (e.g., rotated) and then turned into a matrix where each element represents a pixel.

The matrix is then converted into an byte array where each element represents 8 pixels (i.e., each bit of the element is a pixel, where 0 is black and 1 is white). The lines of the matrix are concatenated after each other in the array.

So, for example, considering x(1,2) as the element on the image matrix X, at line 1 and column 2, the resulting array would be:

```
[ x(0,0), x(1,0), x(2,0), ..., x(174,0), x(175,0), x(0,1), x(1,1), x(2,1), ..., x(175,1), x(0,2), ..., x(174, 263), x(175, 263) ]
```

## Text

The string is rendered with the TTF file defined and turned into an image. From this point on, text is handled the same way as an ordinary image.

There are a few details on how functions `WriteRotate()` and `Rotate()` work on texts:

- `WriteRotate()` will consider the line length as the **height** of the display, instead of the **width**, making the line longer;
- You can use `Rotate()` on a string after it has been converted into an image, but keep in mind that the length of the text was already determined as the **width** of the display, so most likely the text will not take the entire display;

# Next features / Fixes

These are some idea I'd like to implement:

- [x] Prints image (cropping if necessary)
- [x] Prints text (custom font, custom font size)
- [x] Rotate image / text
- [x] Position text and image on display

- [ ] Write text with attributes (bold, italic)
- [ ] Compose screen (overlays)
- [ ] Print negative (if black, print as white and vice-versa)
- [ ] Program functionalities for buttons (key1 to key4 on e-Paper HAT)
- [ ] Partial refresh (i.e., update just a region of the display, instead of the whole display)
- [ ] Improve clearscreen time
- [ ] Text seems to be fading the closer it gets to the end of the "line"...

# Other Notes

The images and fonts used on the examples are for demonstration purposes only, and unfortunately, I don't have the details about them anymore. To the best of my knowledge these are all free-to-use resources available on the internet - if that is not true, please let me know and I'll remove them.

Also, even if these are indeed free-to-use but you are the author, drop a note and I'll be glad to add the credits here.

- Font used in the demo code is called "Pastel Colors". Source and author unknown...
- Image used in the demo code had no title or author defined...
