//go:build cgo

package wgpu

// #include <stdlib.h>
import "C"
import "unsafe"

func sizeOf[T any]() C.size_t {
	var tZero T
	size := unsafe.Sizeof(tZero)
	return C.size_t(size)
}

func calloc[T any](n int) *T {
	return (*T)(C.calloc(C.size_t(n), sizeOf[T]()))
}

func free[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}
