package wgpu

import (
	"errors"
	"runtime"
	"strings"
)

// #include "gen_wgpu_wrappers.h"
import "C"

import (
	"unsafe"
)

var allocErrorCallbackValue = makeTypedPool[errorCallback]()

type errorCallback struct {
	callers [1]uintptr
	errStr  *C.char
}

func acquireErrorCallback() *errorCallback {
	cb := allocErrorCallbackValue.Get()

	if cb.callers[0] != 0 {
		// someone else initialized this instance after calling Done()
		// This should never happen.
		panic("errorCallback missuse detected")
	}

	// get the location of the caller for error reporting
	runtime.Callers(2, cb.callers[:])

	if cb.errStr == nil {
		cb.errStr = calloc[C.char](int(C.error_buf_size()))

		// if this instance gets cleaned up, we need to release
		// the error memory too
		runtime.AddCleanup(cb, cfree, cb.errStr)
	}

	// clear string by setting the first byte to zero
	*cb.errStr = 0

	return cb
}

func cfree(ch *C.char) {
	free(ch)
}

func (v *errorCallback) Done() {
	if v.callers[0] == 0 {
		// someone else called done on this instance.
		// This should never happen.
		panic("errorCallback missuse detected")
	}

	// use the callers field to indicate that the handle is initialized
	v.callers[0] = 0

	// return the instance back to the pool
	allocErrorCallbackValue.Put(v)
}

func (v *errorCallback) ToPointer() *C.char {
	return v.errStr
}

func (v *errorCallback) ToError() error {
	if *v.errStr == 0 {
		return nil
	}

	frames := runtime.CallersFrames(v.callers[:])
	frame, _ := frames.Next()

	context := frame.Func.Name()

	// strip github.com/.../ from method name
	if idx := strings.LastIndexByte(context, '/'); idx >= 0 {
		context = context[idx+1:]
	}

	// convert the message to a string
	message := C.GoString(v.errStr)

	return Error{
		Context: context,
		Wrapped: errors.New(message),
	}
}

//export gowebgpu_error_callback_go
func gowebgpu_error_callback_go(_type C.WGPUErrorType, message C.WGPUStringView, userdata unsafe.Pointer) {
}
