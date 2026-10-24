//go:build js

package wgpu

import (
	"fmt"
	"syscall/js"
)

var (
	jsObject = js.Global().Get("Object")
	jsArray  = js.Global().Get("Array")
)

// reusedObjects hands out JS objects that are filled again on every call instead of
// creating new ones, which avoids the object creation and the finalizer syscall/js
// registers for each new object. This is safe because WebGPU copies a descriptor when
// the method is called, and Go on wasm runs on a single thread.
// Every fill must set all fields, using undefined for absent ones.
type reusedObjects struct {
	ctor    js.Value
	objects []js.Value
}

func (c *reusedObjects) get(idx int) js.Value {
	for len(c.objects) <= idx {
		c.objects = append(c.objects, c.ctor.New())
	}
	return c.objects[idx]
}

// enumToJS converts the given non-bit-flag enum value to a type that
// can be passed as an argument to JavaScript. Bit flag enums should be
// passed as a uint.
//
// The JS strings are cached: passing a Go string to JS creates a new JS string,
// a heap allocation and a finalizer every time.
func enumToJS(s fmt.Stringer) js.Value {
	ss := s.String()
	if ss == "undefined" {
		return js.Undefined()
	}

	value, ok := enumStrings[ss]
	if !ok {
		value = js.ValueOf(ss)
		enumStrings[ss] = value
	}

	return value
}

// enumStrings needs no lock, Go on wasm runs on a single thread.
var enumStrings = map[string]js.Value{}

// pointerToJS converts the given pointer value to a type that can be
// passed as an argument to JavaScript. It must implement a toJS method.
// Returning js.Value instead of any avoids a heap allocation per call for handles.
func pointerToJS[T, R any, P interface {
	*T
	toJS() R
}](v P) js.Value {
	if v == nil {
		return js.Undefined()
	}
	return js.ValueOf(v.toJS())
}

// uint32ToJS converts the given uint32 value to a type that can be
// passed as an argument to JavaScript.
func uint32ToJS(v uint32) any {
	if v == LimitU32Undefined {
		return js.Undefined()
	}
	return v
}

// uint64ToJS converts the given uint64 value to a type that can be
// passed as an argument to JavaScript.
func uint64ToJS(v uint64) any {
	if v == LimitU64Undefined {
		return js.Undefined()
	}
	return v
}

func (g Color) toJS() any {
	return []any{g.R, g.G, g.B, g.A}
}

func (g Extent3D) toJS() any {
	return []any{g.Width, g.Height, g.DepthOrArrayLayers}
}

func (g Origin3D) toJS() any {
	return []any{g.X, g.Y, g.Z}
}

func (g *RequestAdapterOptions) toJS() any {
	result := make(map[string]any)
	result["powerPreference"] = enumToJS(g.PowerPreference)
	result["forceFallbackAdapter"] = g.ForceFallbackAdapter
	return result
}

func (g *DeviceDescriptor) toJS() any {
	result := make(map[string]any)
	result["label"] = g.Label
	result["requiredFeatures"] = mapSlice(g.RequiredFeatures, func(f FeatureName) any { return f.String() })
	if g.RequiredLimits != nil {
		result["requiredLimits"] = g.RequiredLimits.toJS()
	}
	return result
}

func (g *SurfaceConfiguration) toJS() any {
	result := make(map[string]any)
	result["usage"] = uint32(g.Usage)
	result["format"] = enumToJS(g.Format)
	result["alphaMode"] = enumToJS(g.AlphaMode)
	result["viewFormats"] = mapSlice(g.ViewFormats, func(f TextureFormat) any {
		return enumToJS(f)
	})
	return result
}

func (g *TextureDescriptor) toJS() any {
	return map[string]any{
		"label":         g.Label,
		"usage":         uint32(g.Usage),
		"dimension":     enumToJS(g.Dimension),
		"size":          g.Size.toJS(),
		"format":        enumToJS(g.Format),
		"mipLevelCount": g.MipLevelCount,
		"sampleCount":   g.SampleCount,
		"viewFormats": mapSlice(g.ViewFormats, func(f TextureFormat) any {
			return enumToJS(f)
		}),
	}
}

func (g *TextureViewDescriptor) toJS() any {
	return map[string]any{
		"label":           g.Label,
		"format":          enumToJS(g.Format),
		"dimension":       enumToJS(g.Dimension),
		"baseMipLevel":    g.BaseMipLevel,
		"mipLevelCount":   g.MipLevelCount,
		"baseArrayLayer":  g.BaseArrayLayer,
		"arrayLayerCount": g.ArrayLayerCount,
		"aspect":          enumToJS(g.Aspect),
	}
}

func (g *CommandEncoderDescriptor) toJS() any {
	return map[string]any{"label": g.Label}
}

func (g *RenderBundleEncoderDescriptor) toJS() any {
	return map[string]any{
		"label":              g.Label,
		"colorFormats":       mapSlice(g.ColorFormats, func(s TextureFormat) any { return enumToJS(s) }),
		"depthStencilFormat": enumToJS(g.DepthStencilFormat),
		"sampleCount":        g.SampleCount,
		"depthReadOnly":      g.DepthReadOnly,
		"stencilReadOnly":    g.StencilReadOnly,
	}
}

func (g *CommandBufferDescriptor) toJS() any {
	return map[string]any{"label": g.Label}
}

func (g *BufferDescriptor) toJS() any {
	return map[string]any{
		"label":            g.Label,
		"size":             g.Size,
		"usage":            uint32(g.Usage),
		"mappedAtCreation": g.MappedAtCreation,
	}
}

func (g *TexelCopyBufferInfo) toJS() any {
	return map[string]any{
		"buffer":       pointerToJS(g.Buffer),
		"offset":       g.Layout.Offset,
		"bytesPerRow":  g.Layout.BytesPerRow,
		"rowsPerImage": uint32ToJS(g.Layout.RowsPerImage),
	}
}

func (g *TexelCopyTextureInfo) toJS() any {
	return map[string]any{
		"texture":  pointerToJS(g.Texture),
		"mipLevel": g.MipLevel,
		"origin":   g.Origin.toJS(),
		"aspect":   enumToJS(g.Aspect),
	}
}

func (g *TexelCopyBufferLayout) toJS() any {
	return map[string]any{
		"offset":       g.Offset,
		"bytesPerRow":  g.BytesPerRow,
		"rowsPerImage": g.RowsPerImage,
	}
}

func (g *QuerySetDescriptor) toJS() any {
	return map[string]any{
		"label": g.Label,
		"type":  enumToJS(g.Type),
		"count": g.Count,
	}
}

func (g *PassTimestampWrites) toJS() any {
	if g == nil || g.QuerySet == nil {
		return nil
	}
	result := map[string]any{"querySet": g.QuerySet.toJS()}
	if g.BeginningOfPassWriteIndex != QuerySetIndexUndefined {
		result["beginningOfPassWriteIndex"] = g.BeginningOfPassWriteIndex
	}
	if g.EndOfPassWriteIndex != QuerySetIndexUndefined {
		result["endOfPassWriteIndex"] = g.EndOfPassWriteIndex
	}
	return result
}

var (
	renderPassDescriptor       = jsObject.New()
	renderPassColorArray       = jsArray.New()
	renderPassColorAttachments = reusedObjects{ctor: jsObject}
	renderPassClearValues      = reusedObjects{ctor: jsArray}
	renderPassDepthStencil     = jsObject.New()
)

func (g *RenderPassDescriptor) toJS() js.Value {
	renderPassColorArray.Set("length", len(g.ColorAttachments))
	for idx := range g.ColorAttachments {
		renderPassColorArray.SetIndex(idx, g.ColorAttachments[idx].toJS(idx))
	}

	var writes any = js.Undefined()
	if w := g.TimestampWrites.toJS(); w != nil {
		writes = w
	}

	result := renderPassDescriptor
	result.Set("colorAttachments", renderPassColorArray)
	result.Set("depthStencilAttachment", pointerToJS(g.DepthStencilAttachment))
	result.Set("timestampWrites", writes)
	return result
}

// toJS fills the reused attachment object at idx.
func (g *RenderPassColorAttachment) toJS(idx int) js.Value {
	clearValue := renderPassClearValues.get(idx)
	clearValue.SetIndex(0, g.ClearValue.R)
	clearValue.SetIndex(1, g.ClearValue.G)
	clearValue.SetIndex(2, g.ClearValue.B)
	clearValue.SetIndex(3, g.ClearValue.A)

	result := renderPassColorAttachments.get(idx)
	result.Set("view", g.View.jsValue)
	result.Set("loadOp", enumToJS(g.LoadOp))
	result.Set("storeOp", enumToJS(g.StoreOp))
	result.Set("clearValue", clearValue)
	result.Set("resolveTarget", pointerToJS(g.ResolveTarget))
	return result
}

func (g *RenderPassDepthStencilAttachment) toJS() js.Value {
	result := renderPassDepthStencil
	result.Set("view", pointerToJS(g.View))
	result.Set("depthClearValue", g.DepthClearValue)
	result.Set("stencilClearValue", g.StencilClearValue)
	result.Set("depthReadOnly", g.DepthReadOnly)
	result.Set("stencilReadOnly", g.StencilReadOnly)

	// https://www.w3.org/TR/webgpu/#abstract-opdef-gpurenderpassdepthstencilattachment-gpurenderpassdepthstencilattachment-valid-usage
	var depthLoadOp, depthStoreOp, stencilLoadOp, stencilStoreOp any = js.Undefined(), js.Undefined(), js.Undefined(), js.Undefined()
	if !g.DepthReadOnly {
		depthLoadOp, depthStoreOp = enumToJS(g.DepthLoadOp), enumToJS(g.DepthStoreOp)
	}
	if !g.StencilReadOnly {
		stencilLoadOp, stencilStoreOp = enumToJS(g.StencilLoadOp), enumToJS(g.StencilStoreOp)
	}

	result.Set("depthLoadOp", depthLoadOp)
	result.Set("depthStoreOp", depthStoreOp)
	result.Set("stencilLoadOp", stencilLoadOp)
	result.Set("stencilStoreOp", stencilStoreOp)
	return result
}

func (g *RenderPipelineDescriptor) toJS() any {
	result := make(map[string]any)
	if g.Layout != nil {
		result["layout"] = pointerToJS(g.Layout)
	} else {
		result["layout"] = "auto"
	}
	result["vertex"] = g.Vertex.toJS()
	result["primitive"] = g.Primitive.toJS()
	result["depthStencil"] = pointerToJS(g.DepthStencil)
	result["multisample"] = g.Multisample.toJS()
	result["fragment"] = pointerToJS(g.Fragment)
	return result
}

func (g *SamplerDescriptor) toJS() any {
	result := make(map[string]any)
	result["addressModeU"] = enumToJS(g.AddressModeU)
	result["addressModeV"] = enumToJS(g.AddressModeV)
	result["addressModeW"] = enumToJS(g.AddressModeW)
	result["magFilter"] = enumToJS(g.MagFilter)
	result["minFilter"] = enumToJS(g.MinFilter)
	result["mipmapFilter"] = enumToJS(g.MipmapFilter)
	result["lodMinClamp"] = g.LodMinClamp
	result["lodMaxClamp"] = g.LodMaxClamp
	result["compare"] = enumToJS(g.Compare)
	result["maxAnisotropy"] = g.MaxAnisotropy
	return result
}

func (g *ProgrammableStageDescriptor) toJS() any {
	return map[string]any{
		"module":     pointerToJS(g.Module),
		"entryPoint": g.EntryPoint,
	}
}

// toJS converts the limits to a JS GPUSupportedLimits-style object,
// only including the limits that are not undefined.
func (g *Limits) toJS() map[string]any {
	result := make(map[string]any)
	u32 := func(name string, v uint32) {
		if v != LimitU32Undefined {
			result[name] = v
		}
	}
	u64 := func(name string, v uint64) {
		if v != LimitU64Undefined {
			result[name] = v
		}
	}
	u32("maxTextureDimension1D", g.MaxTextureDimension1D)
	u32("maxTextureDimension2D", g.MaxTextureDimension2D)
	u32("maxTextureDimension3D", g.MaxTextureDimension3D)
	u32("maxTextureArrayLayers", g.MaxTextureArrayLayers)
	u32("maxBindGroups", g.MaxBindGroups)
	u32("maxBindingsPerBindGroup", g.MaxBindingsPerBindGroup)
	u32("maxDynamicUniformBuffersPerPipelineLayout", g.MaxDynamicUniformBuffersPerPipelineLayout)
	u32("maxDynamicStorageBuffersPerPipelineLayout", g.MaxDynamicStorageBuffersPerPipelineLayout)
	u32("maxSampledTexturesPerShaderStage", g.MaxSampledTexturesPerShaderStage)
	u32("maxSamplersPerShaderStage", g.MaxSamplersPerShaderStage)
	u32("maxStorageBuffersPerShaderStage", g.MaxStorageBuffersPerShaderStage)
	u32("maxStorageTexturesPerShaderStage", g.MaxStorageTexturesPerShaderStage)
	u32("maxUniformBuffersPerShaderStage", g.MaxUniformBuffersPerShaderStage)
	u64("maxUniformBufferBindingSize", g.MaxUniformBufferBindingSize)
	u64("maxStorageBufferBindingSize", g.MaxStorageBufferBindingSize)
	u32("minUniformBufferOffsetAlignment", g.MinUniformBufferOffsetAlignment)
	u32("minStorageBufferOffsetAlignment", g.MinStorageBufferOffsetAlignment)
	u32("maxVertexBuffers", g.MaxVertexBuffers)
	u64("maxBufferSize", g.MaxBufferSize)
	u32("maxVertexAttributes", g.MaxVertexAttributes)
	u32("maxVertexBufferArrayStride", g.MaxVertexBufferArrayStride)
	u32("maxInterStageShaderVariables", g.MaxInterStageShaderVariables)
	u32("maxColorAttachments", g.MaxColorAttachments)
	u32("maxColorAttachmentBytesPerSample", g.MaxColorAttachmentBytesPerSample)
	u32("maxComputeWorkgroupStorageSize", g.MaxComputeWorkgroupStorageSize)
	u32("maxComputeInvocationsPerWorkgroup", g.MaxComputeInvocationsPerWorkgroup)
	u32("maxComputeWorkgroupSizeX", g.MaxComputeWorkgroupSizeX)
	u32("maxComputeWorkgroupSizeY", g.MaxComputeWorkgroupSizeY)
	u32("maxComputeWorkgroupSizeZ", g.MaxComputeWorkgroupSizeZ)
	u32("maxComputeWorkgroupsPerDimension", g.MaxComputeWorkgroupsPerDimension)
	u32("maxImmediateSize", g.MaxImmediateSize)
	return result
}

func limitsFromJS(j js.Value) Limits {
	return Limits{
		MaxTextureDimension1D:                     uint32(j.Get("maxTextureDimension1D").Int()),
		MaxTextureDimension2D:                     uint32(j.Get("maxTextureDimension2D").Int()),
		MaxTextureDimension3D:                     uint32(j.Get("maxTextureDimension3D").Int()),
		MaxTextureArrayLayers:                     uint32(j.Get("maxTextureArrayLayers").Int()),
		MaxBindGroups:                             uint32(j.Get("maxBindGroups").Int()),
		MaxDynamicUniformBuffersPerPipelineLayout: uint32(j.Get("maxDynamicUniformBuffersPerPipelineLayout").Int()),
		MaxDynamicStorageBuffersPerPipelineLayout: uint32(j.Get("maxDynamicStorageBuffersPerPipelineLayout").Int()),
		MaxSampledTexturesPerShaderStage:          uint32(j.Get("maxSampledTexturesPerShaderStage").Int()),
		MaxSamplersPerShaderStage:                 uint32(j.Get("maxSamplersPerShaderStage").Int()),
		MaxStorageBuffersPerShaderStage:           uint32(j.Get("maxStorageBuffersPerShaderStage").Int()),
		MaxStorageTexturesPerShaderStage:          uint32(j.Get("maxStorageTexturesPerShaderStage").Int()),
		MaxUniformBuffersPerShaderStage:           uint32(j.Get("maxUniformBuffersPerShaderStage").Int()),
		MaxUniformBufferBindingSize:               uint64(j.Get("maxUniformBufferBindingSize").Int()),
		MaxStorageBufferBindingSize:               uint64(j.Get("maxStorageBufferBindingSize").Int()),
		MinUniformBufferOffsetAlignment:           uint32(j.Get("minUniformBufferOffsetAlignment").Int()),
		MinStorageBufferOffsetAlignment:           uint32(j.Get("minStorageBufferOffsetAlignment").Int()),
		MaxVertexBuffers:                          uint32(j.Get("maxVertexBuffers").Int()),
		MaxBufferSize:                             uint64(j.Get("maxBufferSize").Int()),
		MaxVertexAttributes:                       uint32(j.Get("maxVertexAttributes").Int()),
		MaxVertexBufferArrayStride:                uint32(j.Get("maxVertexBufferArrayStride").Int()),
		// MaxInterStageShaderComponents:             uint32(j.Get("maxInterStageShaderComponents").Int()), // no present on firefox
		MaxInterStageShaderVariables:      uint32(j.Get("maxInterStageShaderVariables").Int()),
		MaxColorAttachments:               uint32(j.Get("maxColorAttachments").Int()),
		MaxColorAttachmentBytesPerSample:  uint32(j.Get("maxColorAttachmentBytesPerSample").Int()),
		MaxComputeWorkgroupStorageSize:    uint32(j.Get("maxComputeWorkgroupStorageSize").Int()),
		MaxComputeInvocationsPerWorkgroup: uint32(j.Get("maxComputeInvocationsPerWorkgroup").Int()),
		MaxComputeWorkgroupSizeX:          uint32(j.Get("maxComputeWorkgroupSizeX").Int()),
		MaxComputeWorkgroupSizeY:          uint32(j.Get("maxComputeWorkgroupSizeY").Int()),
		MaxComputeWorkgroupSizeZ:          uint32(j.Get("maxComputeWorkgroupSizeZ").Int()),
		MaxComputeWorkgroupsPerDimension:  uint32(j.Get("maxComputeWorkgroupsPerDimension").Int()),
		MaxImmediateSize:                  128,
	}
}
