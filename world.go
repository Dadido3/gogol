// world
package main

import (
	"image"
	_ "image/color"
	"math/rand"

	"github.com/microo8/blackcl"
)

type Coord struct {
	X, Y int
}

type World struct {
	width, height int
	data1, data2  *blackcl.Bytes

	imageUpdates  chan image.Image
	newPixelQueue chan Coord
}

func NewWorld(width, height int) (world *World, err error) {
	data1, err := openGlDevice.NewBytes((width + 2) * (height + 2))
	if err != nil {
		return nil, err
	}

	data2, err := openGlDevice.NewBytes((width + 2) * (height + 2))
	if err != nil {
		return nil, err
	}

	world = &World{
		width:         width,
		height:        height,
		data1:         data1,
		data2:         data2,
		imageUpdates:  make(chan image.Image, 1),
		newPixelQueue: make(chan Coord, 10),
	}

	// Init data
	array := make([]byte, (width+2)*(height+2))
	for i, _ := range array {
		array[i] = byte(rand.Intn(2) * 255)
	}
	err = <-data1.Copy(array)
	if err != nil {
		return nil, err
	}
	err = <-data2.Copy(array)
	if err != nil {
		return nil, err
	}

	return
}

func (w *World) Update(iterations int) error {

loop:
	for {
		select {
		case newPixelPos := <-w.newPixelQueue:
			if newPixelPos.X >= 0 && newPixelPos.X < w.width+2 && newPixelPos.Y >= 0 && newPixelPos.Y < w.height+2 {
				array, err := w.data1.Data()
				if err != nil {
					return err
				}
				array[newPixelPos.X+newPixelPos.Y*(w.width+2)] = 255
				err = <-w.data1.Copy(array)
				if err != nil {
					return err
				}
			}
		default:
			break loop
		}
	}

	err := <-clKernel.Global(w.width, w.height).Local(20, 20).Run(w.data1, w.data2)
	if err != nil {
		return err
	}

	w.data1, w.data2 = w.data2, w.data1

	// If possible write the state to the graphics output channel
	if len(w.imageUpdates) == 0 {
		array, err := w.data1.Data()
		if err != nil {
			return err
		}

		img := image.NewGray(image.Rect(0, 0, w.width+2, w.height+2))
		img.Pix = array

		w.imageUpdates <- img
	}

	return nil
}
