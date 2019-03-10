// world
package main

import (
	"image"
	"math/rand"

	"github.com/Dadido3/blackcl"
)

type coord struct {
	X, Y int
}

type world struct {
	width, height int
	data1, data2  *blackcl.Bytes

	imageUpdates  chan image.Image
	newPixelQueue chan coord
}

func newWorld(width, height int) (w *world, err error) {
	data1, err := openGlDevice.NewBytes((width + 2) * (height + 2))
	if err != nil {
		return nil, err
	}

	data2, err := openGlDevice.NewBytes((width + 2) * (height + 2))
	if err != nil {
		return nil, err
	}

	w = &world{
		width:         width,
		height:        height,
		data1:         data1,
		data2:         data2,
		imageUpdates:  make(chan image.Image, 1),
		newPixelQueue: make(chan coord, 10),
	}

	// Init data
	array := make([]byte, (width+2)*(height+2))
	for i := range array {
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

func (w *world) update(iterations int) error {

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

	for i := 0; i < iterations-1; i++ {
		_, err := clKernel.GlobalOffset(1, 1).Global(w.width, w.height).Run(false, nil, w.data1, w.data2)
		if err != nil {
			return err
		}
		w.data1, w.data2 = w.data2, w.data1
	}
	event, err := clKernel.GlobalOffset(1, 1).Global(w.width, w.height).Run(true, nil, w.data1, w.data2)
	if err != nil {
		return err
	}
	defer event.Release()
	event.Wait()

	w.data1, w.data2 = w.data2, w.data1

	// If possible write the state to the graphics output channel
	if len(w.imageUpdates) == 0 {
		go func() {
			array, _ := w.data1.Data()

			img := image.NewGray(image.Rect(0, 0, w.width+2, w.height+2))
			img.Pix = array

			w.imageUpdates <- img
		}()
	}

	return nil
}
