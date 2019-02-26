// openCl
package main

import (
	"fmt"

	"github.com/Dadido3/blackcl"
)

var openGlDevice *blackcl.Device
var clKernel *blackcl.Kernel

const golKernel = `
__kernel void gol(__global char* src, __global char* dest) {
	const int2 pos = {get_global_id(0), get_global_id(1)};
	const int width = get_global_size(0)+2;
	char counter = 0;

	for (int y1 = -1; y1 <= 1; y1++) {
		for (int x1 = -1; x1 <= 1; x1++) {
			if ((x1 != 0 || y1 != 0) && src[pos.x+x1+width*(pos.y+y1)]) {
				counter ++;
			}
		}
	}
	
	char new = (counter == 3 || (counter == 2 && src[pos.x+width*pos.y]));
	
	dest[pos.x+width*pos.y] = new * 255;
	//dest[pos.x+width*pos.y] = src[pos.x+width*pos.y];
}`

func initComputeDevice(deviceNumber int) error {
	devices, err := blackcl.GetDevices(blackcl.DeviceTypeAll)
	if err != nil {
		return err
	}

	for i, device := range devices {
		if i == deviceNumber {
			fmt.Printf("-> ID %d: %s (%s, %s, %s)\n", i, device.Name(), device.DriverVersion(), device.OpenCLCVersion(), device.Version())
		} else {
			fmt.Printf("   ID %d: %s (%s, %s, %s)\n", i, device.Name(), device.DriverVersion(), device.OpenCLCVersion(), device.Version())
		}
	}
	openGlDevice = devices[deviceNumber]

	openGlDevice.AddProgram(golKernel)
	clKernel = openGlDevice.Kernel("gol")

	return nil
}
