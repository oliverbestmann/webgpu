package wgpu

import "unsafe"

func sizeOf[T any]() uintptr {
	var tZero T
	return unsafe.Sizeof(tZero)
}
