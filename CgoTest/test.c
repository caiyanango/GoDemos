#include <stdint.h> // for uintptr_t

// A Go function
void MyGoPrint(uintptr_t handle);

// A C function
void myprint(uintptr_t handle) {
    MyGoPrint(handle);
}