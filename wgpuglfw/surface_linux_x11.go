//go:build (linux || freebsd || netbsd || openbsd) && !android && ((!x11 && !wayland) || x11)

package wgpuglfw

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/oliverbestmann/webgpu/wgpu"
)

func GetSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	switch glfw.GetPlatform() {

	case glfw.PlatformX11:
		return &wgpu.SurfaceDescriptor{
			XlibWindow: &wgpu.SurfaceSourceXlibWindow{
				Display: unsafe.Pointer(glfw.GetX11Display()),
				Window:  uint32(w.GetX11Window()),
			},
		}

	default:
		panic("Unsupported glfw platform. Build with no tags to support both x11 and wayland.")
	}
}
