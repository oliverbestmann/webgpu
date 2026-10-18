//go:build !js

package wgpu

import (
	"errors"

	_ "github.com/oliverbestmann/webgpu/libs-android"

	_ "github.com/oliverbestmann/webgpu/libs-darwin"

	_ "github.com/oliverbestmann/webgpu/libs-ios"

	_ "github.com/oliverbestmann/webgpu/libs-linux"

	"runtime"

	_ "github.com/oliverbestmann/webgpu/libs-windows"
)

/*

// Android
#cgo android,amd64 CFLAGS: -I${SRCDIR}/lib/android/amd64
#cgo android,386 CFLAGS: -I${SRCDIR}/lib/android/386
#cgo android,arm64 CFLAGS: -I${SRCDIR}/lib/android/arm64
#cgo android,arm CFLAGS: -I${SRCDIR}/lib/android/arm

// Linux
#cgo linux,!android,amd64 CFLAGS: -I${SRCDIR}/lib/linux/amd64
#cgo linux,!android,arm64 CFLAGS: -I${SRCDIR}/lib/linux/arm64

// iOS
#cgo ios,amd64 CFLAGS: -I${SRCDIR}/lib/ios/amd64
#cgo ios,arm64 CFLAGS: -I${SRCDIR}/lib/ios/arm64

// Darwin
#cgo darwin,!ios,amd64 CFLAGS: -I${SRCDIR}/lib/darwin/amd64
#cgo darwin,!ios,arm64 CFLAGS: -I${SRCDIR}/lib/darwin/arm64
#cgo darwin,!ios       LDFLAGS: -framework Foundation

// Windows
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/lib/windows/amd64
#cgo windows,arm64 CFLAGS: -I${SRCDIR}/lib/windows/arm64
#cgo windows,386 CFLAGS: -I${SRCDIR}/lib/windows/386

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <wgpu.h>

// The message is a WGPUStringView passed by value, and is not NUL-terminated
// unless length is WGPU_STRLEN.
static size_t logMessageLen(WGPUStringView msg) {
	if (msg.data == NULL) return 0;
	if (msg.length == WGPU_STRLEN) return strlen(msg.data);
	return msg.length;
}

#ifdef __ANDROID__
#include <android/log.h>
void logCallback_cgo(WGPULogLevel level, WGPUStringView msg, void *userdata) {
	(void)userdata;
	size_t len = logMessageLen(msg);
	char *text = malloc(len + 1);
	if (text == NULL) return;
	if (len > 0) memcpy(text, msg.data, len);
	text[len] = '\0';
	switch (level) {
	case WGPULogLevel_Error:
		__android_log_write(ANDROID_LOG_ERROR, "GoLogWGPU", text);
		break;
	case WGPULogLevel_Warn:
		__android_log_write(ANDROID_LOG_WARN, "GoLogWGPU", text);
		break;
	default:
		__android_log_write(ANDROID_LOG_INFO, "GoLogWGPU", text);
		break;
	}
	free(text);
}
#else
void logCallback_cgo(WGPULogLevel level, WGPUStringView msg, void *userdata) {
	(void)userdata;
	char const *level_str;
	switch (level) {
	case WGPULogLevel_Error:
		level_str = "Error";
		break;
	case WGPULogLevel_Warn:
		level_str = "Warn";
		break;
	case WGPULogLevel_Info:
		level_str = "Info";
		break;
	case WGPULogLevel_Debug:
		level_str = "Debug";
		break;
	case WGPULogLevel_Trace:
		level_str = "Trace";
		break;
	default:
		level_str = "Unknown Level";
	}
	fprintf(stderr, "[wgpu] [%s] %.*s\n", level_str,
		(int)logMessageLen(msg), msg.data == NULL ? "" : msg.data);
}
#endif

// Registered from C, not from Go: a Go-side cast to WGPULogCallback would hide a
// signature that no longer matches the header, and C type-checks it here.
static void registerLogCallback(void) {
	wgpuSetLogCallback(logCallback_cgo, NULL);
}

*/
import "C"

func init() {
	C.registerLogCallback()
}

func SetLogLevel(level LogLevel) {
	C.wgpuSetLogLevel(C.WGPULogLevel(level))
}

func GetVersion() Version {
	return Version(C.wgpuGetVersion())
}

func (g *Device) addRef() C.WGPUDevice {
	if g.ref == nil {
		panic(errors.New("device already released"))
	}

	C.wgpuDeviceAddRef(g.ref)

	return g.ref
}

type releaser interface{ release() }

func releaseOnGC[T releaser](value T) T {
	runtime.SetFinalizer(value, releaseNow)
	return value
}

func releaseNow(value releaser) {
	value.release()
}

// cBool converts the given Go bool to a C.WGPUBool.
func cBool(b bool) C.WGPUBool {
	if b {
		return 1
	}
	return 0
}

// goBool converts the given C.WGPUBool to a Go bool.
func goBool(b C.WGPUBool) bool {
	return b != 0
}
