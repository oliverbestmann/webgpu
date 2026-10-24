//go:build js

package wgpu

import "syscall/js"

func (g *RenderBundleDescriptor) toJS() js.Value {
	return js.ValueOf(map[string]any{"label": labelToJS(g.Label)})
}

// SetPipeline as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setpipeline
func (g *RenderBundleEncoder) SetPipeline(pipeline *RenderPipeline) {
	g.jsValue.Call("setPipeline", pointerToJS(pipeline))
}

// SetBindGroup as described:
// https://gpuweb.github.io/gpuweb/#gpubindingcommandsmixin-setbindgroup
func (g *RenderBundleEncoder) SetBindGroup(groupIndex uint32, group *BindGroup, dynamicOffsets []uint32) {
	g.jsValue.Call("setBindGroup", groupIndex, pointerToJS(group), mapSlice(dynamicOffsets, func(offset uint32) any {
		return offset
	}))
}

// SetVertexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setvertexbuffer
func (g *RenderBundleEncoder) SetVertexBuffer(slot uint32, buffer *Buffer, offset uint64, size uint64) {
	g.jsValue.Call("setVertexBuffer", slot, pointerToJS(buffer), offset, uint64ToJS(size))
}

// SetIndexBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-setindexbuffer
func (g *RenderBundleEncoder) SetIndexBuffer(buffer *Buffer, format IndexFormat, offset uint64, size uint64) {
	g.jsValue.Call("setIndexBuffer", pointerToJS(buffer), enumToJS(format), offset, uint64ToJS(size))
}

// Draw as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-draw
func (g *RenderBundleEncoder) Draw(vertexCount, instanceCount, firstVertex, firstInstance uint32) {
	g.jsValue.Call("draw", vertexCount, instanceCount, firstVertex, firstInstance)
}

// DrawIndexed as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexed
func (g *RenderBundleEncoder) DrawIndexed(indexCount uint32, instanceCount uint32, firstIndex uint32, baseVertex int32, firstInstance uint32) {
	g.jsValue.Call("drawIndexed", indexCount, instanceCount, firstIndex, baseVertex, firstInstance)
}

// DrawIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindirect
func (g *RenderBundleEncoder) DrawIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	g.jsValue.Call("drawIndirect", pointerToJS(indirectBuffer), indirectOffset)
}

// DrawIndexedIndirect as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurendercommandsmixin-drawindexedindirect
func (g *RenderBundleEncoder) DrawIndexedIndirect(indirectBuffer *Buffer, indirectOffset uint64) {
	g.jsValue.Call("drawIndexedIndirect", pointerToJS(indirectBuffer), indirectOffset)
}

// InsertDebugMarker as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-insertdebugmarker
func (g *RenderBundleEncoder) InsertDebugMarker(markerLabel string) {
	g.jsValue.Call("insertDebugMarker", labelToJS(markerLabel))
}

// PopDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-popdebuggroup
func (g *RenderBundleEncoder) PopDebugGroup() {
	g.jsValue.Call("popDebugGroup")
}

// PushDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-pushdebuggroup
func (g *RenderBundleEncoder) PushDebugGroup(groupLabel string) {
	g.jsValue.Call("pushDebugGroup", labelToJS(groupLabel))
}

// Finish as described:
// https://gpuweb.github.io/gpuweb/#dom-gpurenderbundleencoder-finish
func (g *RenderBundleEncoder) Finish(descriptor *RenderBundleDescriptor) *RenderBundle {
	return &RenderBundle{jsValue: g.jsValue.Call("finish", pointerToJS(descriptor))}
}
