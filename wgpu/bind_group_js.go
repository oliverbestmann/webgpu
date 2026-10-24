//go:build js

package wgpu

import (
	"syscall/js"
)

// BufferBindingLayout as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpubufferbindinglayout
type BufferBindingLayout struct {
	Type             BufferBindingType
	HasDynamicOffset bool
	MinBindingSize   uint64
}

func (g BufferBindingLayout) toJS() any {
	result := make(map[string]any)
	result["type"] = enumToJS(g.Type)
	result["hasDynamicOffset"] = g.HasDynamicOffset
	result["minBindingSize"] = g.MinBindingSize
	return result
}

// SamplerBindingLayout as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpusamplerbindinglayout
type SamplerBindingLayout struct {
	Type SamplerBindingType
}

func (g SamplerBindingLayout) toJS() any {
	result := make(map[string]any)
	result["type"] = enumToJS(g.Type)
	return result
}

// TextureBindingLayout as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gputexturebindinglayout
type TextureBindingLayout struct {
	SampleType    TextureSampleType
	ViewDimension TextureViewDimension
	Multisampled  bool
}

func (g TextureBindingLayout) toJS() any {
	result := make(map[string]any)
	result["sampleType"] = enumToJS(g.SampleType)
	result["viewDimension"] = enumToJS(g.ViewDimension)
	result["multisampled"] = g.Multisampled
	return result
}

// StorageTextureBindingLayout as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpustoragetexturebindinglayout
type StorageTextureBindingLayout struct {
	Access        StorageTextureAccess
	Format        TextureFormat
	ViewDimension TextureViewDimension
}

func (g StorageTextureBindingLayout) toJS() any {
	result := make(map[string]any)
	result["access"] = enumToJS(g.Access)
	result["format"] = enumToJS(g.Format)
	result["viewDimension"] = enumToJS(g.ViewDimension)
	return result
}

// ExternalTextureBindingLayout as described:
type ExternalTextureBindingLayout struct {
	jsValue any
}

func (g ExternalTextureBindingLayout) toJS() any {
	return g.jsValue
}

// BindGroupLayoutEntry as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpubindgrouplayoutentry
type BindGroupLayoutEntry struct {
	Binding         uint32
	Visibility      ShaderStage
	Buffer          BufferBindingLayout
	Sampler         SamplerBindingLayout
	Texture         TextureBindingLayout
	StorageTexture  StorageTextureBindingLayout
	ExternalTexture ExternalTextureBindingLayout
}

func (g BindGroupLayoutEntry) toJS() any {
	result := make(map[string]any)
	result["binding"] = g.Binding
	result["visibility"] = uint32(g.Visibility)
	switch {
	case g.Buffer != BufferBindingLayout{}:
		result["buffer"] = g.Buffer.toJS()
	case g.Sampler != SamplerBindingLayout{}:
		result["sampler"] = g.Sampler.toJS()
	case g.Texture != TextureBindingLayout{}:
		result["texture"] = g.Texture.toJS()
	case g.StorageTexture != StorageTextureBindingLayout{}:
		result["storageTexture"] = g.StorageTexture.toJS()
	case !g.ExternalTexture.jsValue.(js.Value).IsUndefined():
		result["externalTexture"] = g.ExternalTexture.toJS()
	}
	return result
}

func (g BindGroupLayoutDescriptor) toJS() any {
	return map[string]any{
		"entries": mapSlice(g.Entries, func(entry BindGroupLayoutEntry) any {
			return entry.toJS()
		}),
	}
}

var (
	bindGroupDescriptor     = jsObject.New()
	bindGroupEntryArray     = jsArray.New()
	bindGroupEntries        = reusedObjects{ctor: jsObject}
	bindGroupBufferBindings = reusedObjects{ctor: jsObject}
)

// toJS fills the reused entry object at idx.
func (g *BindGroupEntry) toJS(idx int) js.Value {
	var resource js.Value
	switch {
	case g.Sampler != nil:
		resource = g.Sampler.jsValue
	case g.TextureView != nil:
		resource = g.TextureView.jsValue
	default:
		resource = bindGroupBufferBindings.get(idx)
		resource.Set("buffer", pointerToJS(g.Buffer))
		resource.Set("offset", g.Offset)
		resource.Set("size", uint64ToJS(g.Size))
	}

	result := bindGroupEntries.get(idx)
	result.Set("binding", g.Binding)
	result.Set("resource", resource)
	return result
}

func (g BindGroupDescriptor) toJS() js.Value {
	bindGroupEntryArray.Set("length", len(g.Entries))
	for idx := range g.Entries {
		bindGroupEntryArray.SetIndex(idx, g.Entries[idx].toJS(idx))
	}

	bindGroupDescriptor.Set("layout", pointerToJS(g.Layout))
	bindGroupDescriptor.Set("entries", bindGroupEntryArray)
	return bindGroupDescriptor
}
