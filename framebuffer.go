// An image/draw compatible interface to the linux framebuffer
//
// Use Open() to get a framebuffer object, draw on it using the
// facilities of image/draw, and call its Flush() method to sync changes
// to the display.
package framebuffer

import (
	"image"
	"image/color"

	"os"
)

const (
	red   = 2
	green = 1
	blue  = 0
	x     = 3 // not sure what this does, but there's a slot for it.

	colorBytes = 4
)

type FrameBuffer struct {
	buf  []byte
	h, w int
	file *os.File
}

func (fb *FrameBuffer) ColorModel() color.Model {
	return color.RGBAModel
}

func (fb *FrameBuffer) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: fb.w, Y: fb.h},
	}
}

func (fb *FrameBuffer) getPixelStart(x, y int) int {
	return (y*fb.w + x) * colorBytes
}

func (fb *FrameBuffer) At(x, y int) color.Color {
	pixelStart := fb.getPixelStart(x, y)
	return color.RGBA{
		R: fb.buf[pixelStart+red],
		G: fb.buf[pixelStart+green],
		B: fb.buf[pixelStart+blue],
		A: 0,
	}
}

func (fb *FrameBuffer) Set(x, y int, c color.Color) {
	pixelStart := fb.getPixelStart(x, y)
	r, g, b, _ := c.RGBA()
	fb.buf[pixelStart+red] = uint8(r)
	fb.buf[pixelStart+green] = uint8(g)
	fb.buf[pixelStart+blue] = uint8(b)
}

// Sync changes to video memory - nothing will actually appear on the
// screen until this is called.
func (fb *FrameBuffer) Flush() error {
	fb.file.Seek(0, 0)
	_, err := fb.file.Write(fb.buf)
	return err
}

// Opens/initializes the framebuffer with device node located at <filename>.
// width and height should be the width and height of the display, in pixels.
func Open(filename string, width, height int) (*FrameBuffer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &FrameBuffer{buf: make([]byte, height*width*colorBytes), w: width, h: height, file: file}, nil
}
