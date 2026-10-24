//go:build js

package wgpu

import "syscall/js"

// dynamic offsets are copied into this array instead of being converted to []any.
// One array is enough: wasm is single threaded and setBindGroup copies the offsets.
var dynamicOffsetsArray = js.Global().Get("Uint32Array").New(128)

// Invoking a cached function avoids encoding and looking up the method name
// on every call, which matters for methods called many times per frame.
var (
	jsRenderPassEncoderSetPipeline         = renderPassMethod("setPipeline")
	jsRenderPassEncoderSetBindGroup        = renderPassMethod("setBindGroup")
	jsRenderPassEncoderSetVertexBuffer     = renderPassMethod("setVertexBuffer")
	jsRenderPassEncoderSetIndexBuffer      = renderPassMethod("setIndexBuffer")
	jsRenderPassEncoderSetScissorRect      = renderPassMethod("setScissorRect")
	jsRenderPassEncoderSetViewport         = renderPassMethod("setViewport")
	jsRenderPassEncoderDraw                = renderPassMethod("draw")
	jsRenderPassEncoderDrawIndexed         = renderPassMethod("drawIndexed")
	jsRenderPassEncoderDrawIndirect        = renderPassMethod("drawIndirect")
	jsRenderPassEncoderDrawIndexedIndirect = renderPassMethod("drawIndexedIndirect")
)

// renderPassMethod returns a function that is called as f(encoder, ...args).
func renderPassMethod(name string) js.Value {
	class := js.Global().Get("GPURenderPassEncoder")
	if class.IsUndefined() {
		// no WebGPU in this browser, do not panic during package init
		return js.Undefined()
	}

	method := class.Get("prototype").Get(name)
	return js.Global().Get("Function").Get("prototype").Get("call").Call("bind", method)
}

// SetPipeline as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setpipeline
func (g *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	jsRenderPassEncoderSetPipeline.Invoke(g.jsValue, pipeline.jsValue)
}

// SetScissorRect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setscissorrect
func (g *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	jsRenderPassEncoderSetScissorRect.Invoke(g.jsValue, x, y, width, height)
}

// SetStencilReference as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setstencilreference
func (g *RenderPassEncoder) SetStencilReference(ref uint32) {
	g.jsValue.Call("setStencilReference", ref)
}

// SetVertexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setvertexbuffer
func (g *RenderPassEncoder) SetVertexBuffer(slot uint32, vertexBuffer *Buffer, offset, size uint64) {
	jsRenderPassEncoderSetVertexBuffer.Invoke(g.jsValue, slot, pointerToJS(vertexBuffer), offset, uint64ToJS(size))
}

// SetIndexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setindexbuffer
func (g *RenderPassEncoder) SetIndexBuffer(indexBuffer *Buffer, format IndexFormat, offset, size uint64) {
	jsRenderPassEncoderSetIndexBuffer.Invoke(g.jsValue, pointerToJS(indexBuffer), enumToJS(format), offset, uint64ToJS(size))
}

// SetBindGroup as described:
// https://gpuweb.github.io/gpuweb/#gpubindingcommandsmixin-setbindgroup
func (g *RenderPassEncoder) SetBindGroup(index uint32, bindGroup *BindGroup, dynamicOffsets []uint32) {
	for idx, value := range dynamicOffsets {
		dynamicOffsetsArray.SetIndex(idx, value)
	}

	jsRenderPassEncoderSetBindGroup.Invoke(g.jsValue, index, pointerToJS(bindGroup), dynamicOffsetsArray, 0, len(dynamicOffsets))
}

// SetBlendConstant as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setblendconstant
func (g *RenderPassEncoder) SetBlendConstant(color *Color) {
	jsColor := js.ValueOf([]any{color.R, color.G, color.B, color.A})
	g.jsValue.Call("setBlendConstant", jsColor)
}

func (g *RenderPassEncoder) SetImmediates(offset uint32, data []byte) {
	panic("not supported in javascript")
}

// Draw as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-draw
func (g *RenderPassEncoder) Draw(vertexCount uint32, instanceCount, firstVertex, firstInstance uint32) {
	jsRenderPassEncoderDraw.Invoke(g.jsValue, vertexCount, instanceCount, firstVertex, firstInstance)
}

// DrawIndexed as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexed
func (g *RenderPassEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	jsRenderPassEncoderDrawIndexed.Invoke(
		g.jsValue,
		indexCount, instanceCount, firstIndex, baseVertex, firstInstance,
	)
}

// DrawIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindirect
func (g *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	jsRenderPassEncoderDrawIndirect.Invoke(g.jsValue, pointerToJS(indirectBuffer), indirectOffset)
}

// DrawIndexedIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexedindirect
func (g *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	jsRenderPassEncoderDrawIndexedIndirect.Invoke(g.jsValue, pointerToJS(indirectBuffer), indirectOffset)
}

// SetViewport as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setviewport
func (g *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {
	jsRenderPassEncoderSetViewport.Invoke(g.jsValue, x, y, width, height, minDepth, maxDepth)
}

// BeginOcclusionQuery as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-beginocclusionquery
func (g *RenderPassEncoder) BeginOcclusionQuery(queryIndex uint32) {
	g.jsValue.Call("beginOcclusionQuery", queryIndex)
}

// EndOcclusionQuery as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-endocclusionquery
func (g *RenderPassEncoder) EndOcclusionQuery() {
	g.jsValue.Call("endOcclusionQuery")
}

// ExecuteBundles as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-executebundles
func (g *RenderPassEncoder) ExecuteBundles(bundles ...*RenderBundle) {
	g.jsValue.Call("executeBundles", mapSlice(bundles, func(bundle *RenderBundle) any {
		return pointerToJS(bundle)
	}))
}

// InsertDebugMarker as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-insertdebugmarker
func (g *RenderPassEncoder) InsertDebugMarker(markerLabel string) {
	g.jsValue.Call("insertDebugMarker", markerLabel)
}

// PopDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-popdebuggroup
func (g *RenderPassEncoder) PopDebugGroup() {
	g.jsValue.Call("popDebugGroup")
}

// PushDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-pushdebuggroup
func (g *RenderPassEncoder) PushDebugGroup(groupLabel string) {
	g.jsValue.Call("pushDebugGroup", groupLabel)
}

// TryEnd as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-end
func (g *RenderPassEncoder) TryEnd() (err error) {
	defer handleJsException(&err)
	g.jsValue.Call("end")
	return nil
}
