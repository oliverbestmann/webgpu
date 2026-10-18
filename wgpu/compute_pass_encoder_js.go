//go:build js

package wgpu

import (
	"syscall/js"
)

// ComputePassDescriptor as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpucomputepassdescriptor
type ComputePassDescriptor struct {
	Label           string
	TimestampWrites *PassTimestampWrites
}

func (g *ComputePassDescriptor) toJS() any {
	result := map[string]any{"label": g.Label}
	if writes := g.TimestampWrites.toJS(); writes != nil {
		result["timestampWrites"] = writes
	}
	return result
}

// SetPipeline as described:
// https://gpuweb.github.io/gpuweb/#dom-gpucomputepassencoder-setpipeline
func (g *ComputePassEncoder) SetPipeline(pipeline *ComputePipeline) {
	g.jsValue.Call("setPipeline", pointerToJS(pipeline))
}

// SetBindGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpubindingcommandsmixin-setbindgroup
func (g *ComputePassEncoder) SetBindGroup(index uint32, bindGroup *BindGroup, dynamicOffsets []uint32) {
	params := make([]any, 3)
	params[0] = index
	params[1] = pointerToJS(bindGroup)
	params[2] = mapSlice(dynamicOffsets, func(offset uint32) any {
		return offset
	})
	g.jsValue.Call("setBindGroup", params...)
}

// DispatchWorkgroups as described:
// https://gpuweb.github.io/gpuweb/#dom-gpucomputepassencoder-dispatchworkgroups
func (g *ComputePassEncoder) DispatchWorkgroups(workgroupCountX, workgroupCountY, workgroupCountZ uint32) {
	params := make([]any, 3)
	params[0] = workgroupCountX
	if workgroupCountY > 0 {
		params[1] = workgroupCountY
	} else {
		params[1] = js.Undefined()
	}
	if workgroupCountZ > 0 {
		params[2] = workgroupCountZ
	} else {
		params[2] = js.Undefined()
	}
	g.jsValue.Call("dispatchWorkgroups", params...)
}

// InsertDebugMarker as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-insertdebugmarker
func (g *ComputePassEncoder) InsertDebugMarker(markerLabel string) {
	g.jsValue.Call("insertDebugMarker", markerLabel)
}

// PopDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-popdebuggroup
func (g *ComputePassEncoder) PopDebugGroup() {
	g.jsValue.Call("popDebugGroup")
}

// PushDebugGroup as described:
// https://gpuweb.github.io/gpuweb/#dom-gpudebugcommandsmixin-pushdebuggroup
func (g *ComputePassEncoder) PushDebugGroup(groupLabel string) {
	g.jsValue.Call("pushDebugGroup", groupLabel)
}

// TryEnd as described:
// https://gpuweb.github.io/gpuweb/#dom-gpucomputepassencoder-end
func (g *ComputePassEncoder) TryEnd() (err error) {
	defer handleJsException(&err)

	g.jsValue.Call("end")
	return nil
}
