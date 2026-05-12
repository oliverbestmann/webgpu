//go:build (linux || freebsd || netbsd || openbsd) && !android && ((!x11 && !wayland) || wayland)

package wgpuglfw

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/oliverbestmann/webgpu/wgpu"
)

func GetSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	switch glfw.GetPlatform() {

	case glfw.PlatformWayland:
		return &wgpu.SurfaceDescriptor{
			WaylandSurface: &wgpu.SurfaceSourceWaylandSurface{
				Display: unsafe.Pointer(glfw.GetWaylandDisplay()),
				Surface: unsafe.Pointer(w.GetWaylandWindow()),
			},
		}

	default:
		panic("Unsupported glfw platform. Build with no tags to support both x11 and wayland.")
	}
}
