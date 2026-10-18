//go:build js

package wgpu

import (
	"fmt"
	"syscall/js"

	"github.com/oliverbestmann/webgpu/jsx"
)

func (g *Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {
	var desc any = js.Undefined()
	if descriptor != nil {
		dm := descriptor.toJS().(map[string]any)
		// requestDevice rejects limits the browser does not know about,
		// so only pass those that are present in the adapter limits.
		if rl, ok := dm["requiredLimits"].(map[string]any); ok {
			al := g.jsValue.Get("limits")
			for k := range rl {
				if al.Get(k).IsUndefined() {
					delete(rl, k)
				}
			}
		}
		desc = dm
	}
	promise := g.jsValue.Call("requestDevice", desc)

	device, ok := jsx.Await(promise)
	if !ok {
		return nil, fmt.Errorf("failed to request device: %s", rejectionReason(device))
	}
	if !device.Truthy() {
		return nil, fmt.Errorf("failed to request device: no WebGPU device available")
	}
	return &Device{jsValue: device}, nil
}

func (g *Adapter) GetInfo() AdapterInfo {
	return AdapterInfo{} // TODO(kai): implement?
}

func (g *Adapter) GetLimits() Limits {
	return limitsFromJS(g.jsValue.Get("limits"))
}

func (g *Adapter) HasFeature(name FeatureName) bool {
	hasFeature := g.jsValue.Get("features").Call("has", name.String())
	return hasFeature.Bool()
}
