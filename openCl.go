// openCl
package main

import (
	"fmt"

	"github.com/microo8/blackcl"
)

var openGlDevice *blackcl.Device
var clKernel *blackcl.Kernel

const golKernel = `
__kernel void gol(__global uchar* src, __global uchar* dest) {
	const int2 pos = {get_global_id(0)+1, get_global_id(1)+1};
	const int width = get_global_size(0)+2;
	uchar counter = 0;

	for (int y1 = -1; y1 <= 1; y1++) {
		for (int x1 = -1; x1 <= 1; x1++) {
			if ((x1 != 0 || y1 != 0) && src[pos.x+x1+width*(pos.y+y1)]) {
				counter ++;
			}
		}
	}
	
	uchar new = (counter == 3 || (counter == 2 && src[pos.x+width*pos.y]));
	
	dest[pos.x+width*pos.y] = new * 255;
	//dest[pos.x+width*pos.y] = src[pos.x+width*pos.y];
}`

func initComputeDevice() error {
	devices, err := blackcl.GetDevices(blackcl.DeviceTypeAll)
	if err != nil {
		return err
	}
	openGlDevice = devices[2]
	fmt.Println(openGlDevice.Name())

	openGlDevice.AddProgram(golKernel)
	clKernel = openGlDevice.Kernel("gol")

	return nil
}
