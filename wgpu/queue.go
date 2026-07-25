//go:build !js

package wgpu

/*
#include "gen_wgpu_wrappers.h"
extern void gowebgpu_queue_work_done_callback_c(WGPUQueueWorkDoneStatus status, void * userdata);
*/
import "C"
import (
	"unsafe"
)

//export gowebgpu_queue_work_done_callback_go
func gowebgpu_queue_work_done_callback_go(status C.WGPUQueueWorkDoneStatus, message C.WGPUStringView, userdata unsafe.Pointer) {
	handle := lookupHandle(userdata)
	defer handle.Delete()

	cb, ok := handle.Value().(QueueWorkDoneCallback)
	if ok {
		cb(QueueWorkDoneStatus(status))
	}
}

func (p *Queue) OnSubmittedWorkDone(callback QueueWorkDoneCallback) {
	handle := newHandle(callback)

	C.wgpuQueueOnSubmittedWorkDone(p.ref, C.WGPUQueueWorkDoneCallbackInfo{
		mode:      C.WGPUCallbackMode_AllowProcessEvents,
		callback:  C.WGPUQueueWorkDoneCallback(C.gowebgpu_queue_work_done_callback_c),
		userdata1: handle.ToPointer(),
	})
}

func (p *Queue) Submit(commands ...*CommandBuffer) (submissionIndex SubmissionIndex) {
	commandCount := len(commands)
	if commandCount == 0 {
		r := C.wgpuQueueSubmitForIndex(p.ref, 0, nil)
		return SubmissionIndex(r)
	}

	commandRefs, commandRefsSlice := callocSlice[C.WGPUCommandBuffer](commandCount)
	defer free(commandRefs)

	for i, v := range commands {
		commandRefsSlice[i] = v.ref
	}

	r := C.wgpuQueueSubmitForIndex(
		p.ref,
		C.size_t(commandCount),
		commandRefs,
	)
	return SubmissionIndex(r)
}

func (p *Queue) TryWriteBuffer(buffer *Buffer, bufferOffset uint64, data []byte) error {
	errh := acquireErrorCallback()
	defer errh.Done()

	size := len(data)
	if size == 0 {
		C.go_wgpuQueueWriteBuffer(
			p.device,
			errh.ToPointer(),
			p.ref,
			buffer.ref,
			C.uint64_t(bufferOffset),
			nil,
			0,
		)

		return errh.ToError()
	}

	C.go_wgpuQueueWriteBuffer(
		p.device,
		errh.ToPointer(),
		p.ref,
		buffer.ref,
		C.uint64_t(bufferOffset),
		unsafe.Pointer(&data[0]),
		C.size_t(size),
	)

	return errh.ToError()
}

func (p *Queue) TryWriteTexture(destination *TexelCopyTextureInfo, data []byte, dataLayout *TexelCopyBufferLayout, writeSize *Extent3D) error {
	var dst C.WGPUTexelCopyTextureInfo
	if destination != nil {
		dst = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var layout C.WGPUTexelCopyBufferLayout
	if dataLayout != nil {
		layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(dataLayout.Offset),
			bytesPerRow:  C.uint32_t(dataLayout.BytesPerRow),
			rowsPerImage: C.uint32_t(dataLayout.RowsPerImage),
		}
	}

	var writeExtent C.WGPUExtent3D
	if writeSize != nil {
		writeExtent = C.WGPUExtent3D{
			width:              C.uint32_t(writeSize.Width),
			height:             C.uint32_t(writeSize.Height),
			depthOrArrayLayers: C.uint32_t(writeSize.DepthOrArrayLayers),
		}
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	size := len(data)
	if size == 0 {
		C.go_wgpuQueueWriteTexture(
			p.device,
			errh.ToPointer(),
			p.ref,
			&dst,
			nil,
			0,
			&layout,
			&writeExtent,
		)

		return errh.ToError()
	}

	C.go_wgpuQueueWriteTexture(
		p.device,
		errh.ToPointer(),
		p.ref,
		&dst,
		unsafe.Pointer(&data[0]),
		C.size_t(size),
		&layout,
		&writeExtent,
	)

	return errh.ToError()
}
