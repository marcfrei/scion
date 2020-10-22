package main

// See https://github.com/enobufs/go-calls-c-pointer

// #include <stdlib.h>
// typedef void (*data_callback)(char *data);
// void makeCallback(char *data, data_callback cb);
import "C"

import (
	"unsafe"

	"github.com/scionproto/scion/go/lib/snet"
)


//export RunServer
func RunServer(sciondAddr, local *C.char, cb C.data_callback) int {
	localAddr, err := snet.ParseUDPAddr(C.GoString(local))
	if err != nil {
		return 1
	}

	runServer(C.GoString(sciondAddr), *localAddr, func (data string) {
		cdata := C.CString(data)
		C.makeCallback(cdata, cb)
		C.free(unsafe.Pointer(cdata))
	})
	return 0
}

