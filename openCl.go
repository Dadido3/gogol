// openCl
package main

import (
	"fmt"
	_ "io/ioutil"

	"github.com/Dadido3/blackcl"
)

var openGlDevice *blackcl.Device
var clKernel *blackcl.Kernel

const golKernel = `
__kernel void gol(__global char* src, __global char* dest) {
	const int2 pos = {get_global_id(0), get_global_id(1)};
	const int width = get_global_size(0)+2;
	int counter = 0;
	int address = pos.x+width*pos.y;

	for (int y1 = -1; y1 <= 1; y1++) {
		for (int x1 = -1; x1 <= 1; x1++) {
			if ((x1 != 0 || y1 != 0) && src[address + x1 + y1*width]) {
				counter ++;
			}
		}
	}
	
	int new = (counter == 3 || (counter == 2 && src[address]));
	
	dest[address] = new * 255;
	//dest[address] = src[address];
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

	/*p := openGlDevice.AddProgram(golKernel)
	bin, err := p.GetBinaries()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("golAssembly.ptx", bin[0], 0644)
	if err != nil {
		return err
	}*/

	return nil
}
