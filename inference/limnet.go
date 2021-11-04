package main

// #cgo CPPFLAGS: -I/usr/local/include/torch/csrc/api/include
// #cgo CXXFLAGS: -std=c++20
// #cgo LDFLAGS: -L../build -L/libtorch/lib -llimnet -lstdc++ -lc10 -ltorch_cpu
// #include "limnet.h"
import "C"

func main() {
	print(C.initialize(false))
}
