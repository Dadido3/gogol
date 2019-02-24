// OpenCL game of life project main.go
package main

import (
	"flag"
	"fmt"
	_ "os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	width, height = 800, 800
)

func main() {

	deviceNumber := flag.Int("device", 0, "OpenCL device number")

	flag.Parse()

	err := initComputeDevice(*deviceNumber)
	if err != nil {
		panic(err)
	}

	world, err := NewWorld(width, height)
	if err != nil {
		panic(err)
	}

	/*f, err := os.Create("profile.out")
	if err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()*/

	//runTicker := time.NewTicker(time.Second / 60).C
	rateTicker := time.NewTicker(time.Second).C
	go func() {

		counter := 0
		profilingCounter := 0
		for {
			select {
			//case <-runTicker:
			default:
				/*if profilingCounter == 60 {
					pprof.StartCPUProfile(f)
				}
				if profilingCounter == 600 {
					pprof.StopCPUProfile()
				}*/
				counter += 10
				profilingCounter++
				err := world.Update(10)
				if err != nil {
					panic(err)
				}
			case <-rateTicker:
				fmt.Printf("Rate: %d\n", counter)
				counter = 0
			}

		}
	}()

	pixelgl.Run(func() {
		cfg := pixelgl.WindowConfig{
			Title:  "gogol",
			Bounds: pixel.R(0, 0, width+2, height+2),
			VSync:  true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		imd := imdraw.New(nil)
		for !win.Closed() {
			select {
			case img := <-world.imageUpdates:
				pic := pixel.PictureDataFromImage(img)
				sprite := pixel.NewSprite(pic, pic.Bounds())

				win.Clear(colornames.Black)
				imd.Clear()

				sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

				imd.Draw(win)
				win.Update()
			default: // Run with maximum frames, to handle input events
				win.Update()
			}

			if win.Pressed(pixelgl.MouseButtonLeft) {
				coords := Coord{int(win.MousePosition().X), height + 2 - int(win.MousePosition().Y)}
				world.newPixelQueue <- coords
			}

		}
	})
}
