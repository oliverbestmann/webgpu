//go:build js

package wgpu

import (
	"fmt"
	"syscall/js"

	"github.com/oliverbestmann/webgpu/jsx"
)

func (g *Adapter) RequestDevice(descriptor *DeviceDescriptor) (*Device, error) {
	desc := js.Undefined()
	if descriptor != nil {
		desc = descriptor.toJS()
		if descriptor.RequiredLimits != nil {
			// requestDevice rejects limits the browser does not know about,
			// so only pass those that are present in the adapter limits.
			limits := descriptor.RequiredLimits.toJS()
			al := g.jsValue.Get("limits")
			for k := range limits {
				if al.Get(k).IsUndefined() {
					delete(limits, k)
				}
			}
			desc.Set("requiredLimits", limits)
		}
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

// GetFeatures as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuadapter-features
func (g *Adapter) GetFeatures() []FeatureName {
	return featuresFromJS(g.jsValue.Get("features"))
}

// featuresFromJS converts a GPUSupportedFeatures set, skipping features this
// package does not know.
func featuresFromJS(set js.Value) []FeatureName {
	names := js.Global().Get("Array").Call("from", set)

	features := make([]FeatureName, 0, names.Length())
	for i := range names.Length() {
		if feature, ok := enumFromJS[FeatureName](names.Index(i).String()); ok {
			features = append(features, feature)
		}
	}

	return features
}

func (g *Adapter) HasFeature(name FeatureName) bool {
	hasFeature := g.jsValue.Get("features").Call("has", name.String())
	return hasFeature.Bool()
}
