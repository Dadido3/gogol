// OpenCL test project main.go
package main

import (
	"fmt"
	"image"
	_ "image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/microo8/blackcl"
)

func readImage(d *blackcl.Device, path string) (*blackcl.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	i, err := d.NewImageFromImage(img)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func writeImage(img *blackcl.Image, path string) error {
	receivedImg, err := img.Data()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, receivedImg)
}

const golKernel = `
__constant sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_REPEAT | CLK_FILTER_NEAREST;
__kernel void gol(__read_only image2d_t src, __write_only image2d_t dest) {
	const int2 pos = {get_global_id(0), get_global_id(1)};
	int counter = 0;
	int self = 0;

	if (read_imagei(src, sampler, pos).x != 0)
		self = 1;

	for (int y1 = -1; y1 <= 1; y1++) {
		for (int x1 = -1; x1 <=1; x1++) {
			if (x1 != 0 || y1 != 0) {
				int2 tempPos = {pos.x+x1, pos.y+y1};
				if (read_imagei(src, sampler, tempPos).x != 0)
					counter ++;
			}
		}
	}
	
	int new = (counter == 3 || (counter == 2 && self));
	
	float4 pixel = {new, new, new, 1};
	write_imagef(dest, pos, pixel);
}`

func main() {
	devices, err := blackcl.GetDevices(blackcl.DeviceTypeDefault)
	if err != nil {
		log.Fatal(err)
	}
	d := devices[0]
	defer d.Release()
	imgA, err := readImage(d, "data/init.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imgA.Release()
	d.AddProgram(golKernel)
	k := d.Kernel("gol")
	imgB, err := d.NewImage(blackcl.ImageTypeRGBA, imgA.Bounds())
	if err != nil {
		log.Fatal(err)
	}
	defer imgB.Release()

	for i := 0; i < 100; i++ {
		err = <-k.Global(imgA.Bounds().Dx(), imgA.Bounds().Dy()).Local(1, 1).Run(imgA, imgB)
		if err != nil {
			log.Fatal(err)
		}

		imgA, imgB = imgB, imgA
		fmt.Print("#")

	}

	err = writeImage(imgB, "data/output.png")
	if err != nil {
		log.Fatal(err)
	}
}
