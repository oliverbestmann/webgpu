//go:build cgo

package wgpu

// #include <stdlib.h>
import "C"
import "unsafe"

func sizeOf[T any]() uintptr {
	var tZero T
	return unsafe.Sizeof(tZero)
}

func calloc[T any](n int) *T {
	return (*T)(C.calloc(C.size_t(n), C.size_t(sizeOf[T]())))
}

func callocOne[T any]() *T {
	return calloc[T](1)
}

func callocSlice[T any](n int) (*T, []T) {
	ptr := calloc[T](n)
	return ptr, unsafe.Slice(ptr, n)
}

func free[T any](ptr *T) {
	C.free(unsafe.Pointer(ptr))
}
