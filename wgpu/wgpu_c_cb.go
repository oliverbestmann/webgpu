//go:build !js

package wgpu

/*
#cgo CFLAGS: -Werror=incompatible-pointer-types

#include <wgpu.h>

void gowebgpu_buffer_map_callback_c(WGPUMapAsyncStatus status, WGPUStringView message, void * userdata, void * userdata2) {
  extern void gowebgpu_buffer_map_callback_go(WGPUMapAsyncStatus status, void * userdata);
  gowebgpu_buffer_map_callback_go(status, userdata);
}

void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView message, void * userdata1, void * userdata2) {
  extern void gowebgpu_request_adapter_callback_go(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView message, void * userdata);
  gowebgpu_request_adapter_callback_go(status, adapter, message, userdata1);
}

void gowebgpu_request_device_callback_c(WGPURequestDeviceStatus status, WGPUDevice device, WGPUStringView message, void * userdata1, void * userdata2) {
  extern void gowebgpu_request_device_callback_go(WGPURequestDeviceStatus status, WGPUDevice device, WGPUStringView message, void * userdata);
  gowebgpu_request_device_callback_go(status, device, message, userdata1);
}

void gowebgpu_device_lost_callback_c(WGPUDevice const * device, WGPUDeviceLostReason reason, WGPUStringView message, void * userdata1, void * userdata2) {
  extern void gowebgpu_device_lost_callback_go(WGPUDeviceLostReason reason, WGPUStringView message, void * userdata);
  gowebgpu_device_lost_callback_go(reason, message, userdata1);
}

void gowebgpu_queue_work_done_callback_c(WGPUQueueWorkDoneStatus status, WGPUStringView message, void * userdata1, void * userdata2) {
  extern void gowebgpu_queue_work_done_callback_go(WGPUQueueWorkDoneStatus status, void * userdata);
  gowebgpu_queue_work_done_callback_go(status, userdata1);
}

// Each trampoline must match the callback typedef of the header, or the
// arguments arrive shifted at runtime.
static WGPUBufferMapCallback const gowebgpu_check_buffer_map = gowebgpu_buffer_map_callback_c;
static WGPURequestAdapterCallback const gowebgpu_check_request_adapter = gowebgpu_request_adapter_callback_c;
static WGPURequestDeviceCallback const gowebgpu_check_request_device = gowebgpu_request_device_callback_c;
static WGPUDeviceLostCallback const gowebgpu_check_device_lost = gowebgpu_device_lost_callback_c;
static WGPUQueueWorkDoneCallback const gowebgpu_check_queue_work_done = gowebgpu_queue_work_done_callback_c;

*/
import "C"
