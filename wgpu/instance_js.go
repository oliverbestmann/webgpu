//go:build js

package wgpu

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/oliverbestmann/webgpu/jsx"
)

func CreateInstance(descriptor *InstanceDescriptor) *Instance {
	gpu := js.Global().Get("navigator").Get("gpu")
	if !gpu.Truthy() {
		log.Println("WebGPU not supported")
		return nil
	}
	return &Instance{jsValue: gpu}
}

func (g *Instance) RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	adapter, ok := jsx.Await(g.jsValue.Call("requestAdapter", pointerToJS(options)))
	if !ok {
		return nil, fmt.Errorf("failed to request adapter: %s", rejectionReason(adapter))
	}
	if !adapter.Truthy() {
		return nil, fmt.Errorf("failed to request adapter: no WebGPU adapter available")
	}
	return &Adapter{jsValue: adapter}, nil
}

func (g *Instance) EnumerateAdapters(options *InstanceEnumerateAdapterOptons) []*Adapter {
	a, err := g.RequestAdapter(&RequestAdapterOptions{})
	if err != nil {
		log.Println(err)
		return nil
	}
	return []*Adapter{a}
}

// SurfaceDescriptor must contain a valid HTML canvas element on web.
type SurfaceDescriptor struct {
	// Canvas must be specified.
	Canvas js.Value

	Label string
}

func (g *Instance) CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	if descriptor.Canvas.IsUndefined() {
		panic("wgpu.Instance.CreateSurface: descriptor.Canvas must be specified")
	}
	jsContext := descriptor.Canvas.Call("getContext", "webgpu")
	return &Surface{jsContext}
}

func (g *Instance) GenerateReport() any { return nil } // no-op

// rejectionReason describes the value a promise was rejected with, usually a DOMException.
func rejectionReason(reason js.Value) string {
	if reason.IsUndefined() || reason.IsNull() {
		return "unknown error"
	}

	return reason.Call("toString").String()
}
