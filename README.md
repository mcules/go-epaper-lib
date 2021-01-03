# go-epaper-lib
A golang lib to use with Waveshare(tm) e-Paper HAT for Raspberry Pi.

The current version is developed/tested only for the 2.7 inches, black-and-white HAT. Hopefully, new models will be available in the future.

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

# Next features

These are some idea I'd like to implement:

- [x] Prints image (cropping if necessary)
- [x] Prints text (custom font, custom font size)

- [ ] Write text with custom font and attributes (bold, italic)
- [ ] Position text and image on display
- [ ] Rotate image / text
- [ ] Compose screen (overlays)
- [ ] Print negative (if black, print as white and vice-versa)
- [ ] Program functionalities for buttons (key1 to key4 on e-Paper HAT)