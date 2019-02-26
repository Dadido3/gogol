# gogol
A simple OpenCL implementation of `Game of Life` in go.

![](https://github.com/Dadido3/gogol/raw/master/screenshots/gol.png)

## How to build

#### Windows
1. Install `gcc` to make cgo work. Preferably use MinGW.
2. Install any OpenCL SDK. For NVIDIA graphic cards install the `CUDA Toolkit`.
3. Make sure cgo can find the OpenCL header files and libs. If you have installed the `CUDA Toolkit`, add the following environment variables (You may have to correct the version in the path):
	- `CGO_CFLAGS` `-IC:\PROGRA~1\NVIDIA~2\CUDA\v10.0\include\CL -IC:\PROGRA~1\NVIDIA~2\CUDA\v10.0\include`
	- `CGO_LDFLAGS` `-LC:\PROGRA~1\NVIDIA~2\CUDA\v10.0\lib\x64`
4. `go get` and `go build`

#### Linux
TODO: Put steps here
But you know what to do!

## How to use
Start executable by command line and use the `-device=` parameter to select a specific device. E.g. `-device=1` to select device 1. Avaliable devices with IDs will be printed when the executable is run.
