//go:build js

package wgpu

import "syscall/js"

// SetPipeline as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setpipeline
func (g *RenderPassEncoder) SetPipeline(pipeline *RenderPipeline) {
	g.jsValue.Call("setPipeline", pointerToJS(pipeline))
}

// SetScissorRect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setscissorrect
func (g *RenderPassEncoder) SetScissorRect(x, y, width, height uint32) {
	params := make([]any, 4)
	params[0] = x
	params[1] = y
	params[2] = width
	params[3] = height
	g.jsValue.Call("setScissorRect", params...)
}

// SetStencilReference as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setstencilreference
func (g *RenderPassEncoder) SetStencilReference(ref uint32) {
	g.jsValue.Call("setStencilReference", ref)
}

// SetVertexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setvertexbuffer
func (g *RenderPassEncoder) SetVertexBuffer(slot uint32, vertexBuffer *Buffer, offset, size uint64) {
	params := make([]any, 4)
	params[0] = slot
	params[1] = pointerToJS(vertexBuffer)
	params[2] = offset
	params[3] = uint64ToJS(size)
	g.jsValue.Call("setVertexBuffer", params...)
}

// SetIndexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setindexbuffer
func (g *RenderPassEncoder) SetIndexBuffer(indexBuffer *Buffer, format IndexFormat, offset, size uint64) {
	params := make([]any, 4)
	params[0] = pointerToJS(indexBuffer)
	params[1] = enumToJS(format)
	params[2] = offset
	params[3] = uint64ToJS(size)
	g.jsValue.Call("setIndexBuffer", params...)
}

// SetBindGroup as described:
// https://gpuweb.github.io/gpuweb/#gpubindingcommandsmixin-setbindgroup
func (g *RenderPassEncoder) SetBindGroup(index uint32, bindGroup *BindGroup, dynamicOffsets []uint32) {
	params := make([]any, 3)
	params[0] = index
	params[1] = pointerToJS(bindGroup)
	params[2] = mapSlice(dynamicOffsets, func(offset uint32) any {
		return offset
	})
	g.jsValue.Call("setBindGroup", params...)
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
	params := make([]any, 4)
	params[0] = vertexCount
	params[1] = instanceCount
	params[2] = firstVertex
	params[3] = firstInstance
	g.jsValue.Call("draw", params...)
}

// DrawIndexed as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexed
func (g *RenderPassEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	params := make([]any, 5)
	params[0] = indexCount
	params[1] = instanceCount
	params[2] = firstIndex
	params[3] = baseVertex
	params[4] = firstInstance
	g.jsValue.Call("drawIndexed", params...)
}

// DrawIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindirect
func (g *RenderPassEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	g.jsValue.Call("drawIndirect", pointerToJS(indirectBuffer), indirectOffset)
}

// DrawIndexedIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexedindirect
func (g *RenderPassEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	g.jsValue.Call("drawIndexedIndirect", pointerToJS(indirectBuffer), indirectOffset)
}

// SetViewport as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderpassencoder-setviewport
func (g *RenderPassEncoder) SetViewport(x, y, width, height, minDepth, maxDepth float32) {
	g.jsValue.Call("setViewport", x, y, width, height, minDepth, maxDepth)
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
