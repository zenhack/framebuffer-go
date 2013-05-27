// A simple utilty to dispaly an image from stdin (png or jpeg)
// on the linux framebuffer.
//
// use -w and -h to set the width and height of the display, by
// default 1440x900 is assumed.
package main

import (
	_ "image/jpeg"
	_ "image/png"

	"image"
	"image/draw"

	"os"

	"flag"
	"fmt"

	"github.com/zenhack/framebuffer-go"
)

var w = flag.Int("w", 1440, "Width of the display (pixels)")
var h = flag.Int("h", 900, "Height of the display (pixels)")

func main() {
	flag.Parse()

	pic, _, err := image.Decode(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error : Decoding image : ", err)
		os.Exit(1)
	}

	fb, err := framebuffer.Open("/dev/fb0", *w, *h)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error : Opening framebuffer : ", err)
		os.Exit(1)
	}

	rect := image.Rectangle{
		Min: fb.Bounds().Min,
		Max: image.Point{
			X: pic.Bounds().Max.X - pic.Bounds().Min.X + fb.Bounds().Min.X,
			Y: pic.Bounds().Max.Y - pic.Bounds().Min.Y + fb.Bounds().Min.Y,
		},
	}

	draw.Draw(fb, rect, pic, pic.Bounds().Min, draw.Over)
	fb.Flush()
	fb.Close()
}
