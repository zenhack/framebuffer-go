// An image/draw compatible interface to the linux framebuffer
//
// Use Open() to get a framebuffer object, draw on it using the
// facilities of image/draw, and call its Flush() method to sync changes
// to the display.
package framebuffer

// #include "fb.h"
// #include <stdlib.h> /* for C.free */
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"unsafe"
)

const (
	red   = 2
	green = 1
	blue  = 0
	x     = 3 // not sure what this does, but there's a slot for it.

	colorBytes = 4
)

// A framebuffer object. Obtain with Open() - the zero value is not useful.
// call Close() when finished to close the underlying file descriptor.
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

// Closes the framebuffer
func (fb *FrameBuffer) Close() error {
	return fb.file.Close()
}

// Opens/initializes the framebuffer with device node located at <filename>.
func Open(filename string) (*FrameBuffer, error) {
	var cFilename *C.char
	cFilename = C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	var info C.fb_info_t
	cErr := C.initfb(cFilename, &info)
	if cErr != 0 {
		return nil, fmt.Errorf("Error initializing framebuffer")
	}

	return &FrameBuffer{
		buf: make([]byte, info.fix_info.smem_len),
		file: os.NewFile(uintptr(info.fd), filename)}, nil
}
