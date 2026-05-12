//go:build (linux || freebsd || netbsd || openbsd) && !android

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

	case glfw.PlatformWayland:
		return &wgpu.SurfaceDescriptor{
			WaylandSurface: &wgpu.SurfaceSourceWaylandSurface{
				Display: unsafe.Pointer(glfw.GetWaylandDisplay()),
				Surface: unsafe.Pointer(w.GetWaylandWindow()),
			},
		}

	default:
		panic("unsupported glfw platform")
	}
}
