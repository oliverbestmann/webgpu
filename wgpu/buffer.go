//go:build !js

package wgpu

/*
#include "gen_wgpu_wrappers.h"
extern void gowebgpu_buffer_map_callback_c(WGPUMapAsyncStatus status, WGPUStringView message, void *userdata, void *userdata2);
*/
import "C"
import (
	"unsafe"
)

func (p *Buffer) Destroy() {
	C.wgpuBufferDestroy(p.ref)
}

func (p *Buffer) GetMappedRange(offset, size uint) []byte {
	buf := C.wgpuBufferGetMappedRange(p.ref, C.size_t(offset), C.size_t(size))
	return unsafe.Slice((*byte)(buf), size)
}

func (p *Buffer) GetSize() uint64 {
	return uint64(C.wgpuBufferGetSize(p.ref))
}

func (p *Buffer) GetUsage() BufferUsage {
	return BufferUsage(C.wgpuBufferGetUsage(p.ref))
}

//export gowebgpu_buffer_map_callback_go
func gowebgpu_buffer_map_callback_go(status C.WGPUMapAsyncStatus, userdata unsafe.Pointer) {
	handle := lookupHandle(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(BufferMapCallback)
	if ok {
		cb(MapAsyncStatus(status))
	}
}

func (p *Buffer) TryMapAsync(mode MapMode, offset uint64, size uint64, callback BufferMapCallback) error {
	callbackHandle := newHandle(callback)

	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuBufferMapAsync(
		p.device,
		errh.ToPointer(),
		p.ref,
		C.WGPUMapMode(mode),
		C.size_t(offset),
		C.size_t(size),
		C.WGPUBufferMapCallbackInfo{
			mode:      C.WGPUCallbackMode_AllowSpontaneous,
			callback:  C.WGPUBufferMapCallback(C.gowebgpu_buffer_map_callback_c),
			userdata1: callbackHandle.ToPointer(),
		},
	)

	return errh.ToError()
}

func (p *Buffer) TryUnmap() error {
	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuBufferUnmap(
		p.device,
		errh.ToPointer(),
		p.ref,
	)

	return errh.ToError()
}
