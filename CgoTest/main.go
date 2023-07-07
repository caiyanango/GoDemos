package main

/*
#include <stdint.h>

void MyGoPrint(uintptr_t handle);
void myprint(uintptr_t handle);
*/
import "C"
import "runtime/cgo"

//export MyGoPrint
func MyGoPrint(handle C.uintptr_t) {
	h := cgo.Handle(handle)
	val := h.Value().(string)
	println(val)
	h.Delete()
}

func main() {
	val := "hello Go"
	C.myprint(C.uintptr_t(cgo.NewHandle(val)))
}
