package cgo

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func ToCString(s string) (cStr *C.char, free func()) {

	if s == "" {
		return nil, func() {}
	}

	cStr = C.CString(s)
	free = func() { C.free(unsafe.Pointer(cStr)) }
	return
}
