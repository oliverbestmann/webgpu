#pragma once

#include <wgpu.h>

// Trampolines defined in wgpu_c_cb.go. Declared once here so the definitions are
// checked against these prototypes, and the prototypes against the header typedefs.
void gowebgpu_buffer_map_callback_c(WGPUMapAsyncStatus status, WGPUStringView message, void * userdata1, void * userdata2);
void gowebgpu_request_adapter_callback_c(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView message, void * userdata1, void * userdata2);
void gowebgpu_request_device_callback_c(WGPURequestDeviceStatus status, WGPUDevice device, WGPUStringView message, void * userdata1, void * userdata2);
void gowebgpu_device_lost_callback_c(WGPUDevice const * device, WGPUDeviceLostReason reason, WGPUStringView message, void * userdata1, void * userdata2);
void gowebgpu_queue_work_done_callback_c(WGPUQueueWorkDoneStatus status, WGPUStringView message, void * userdata1, void * userdata2);
