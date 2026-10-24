//go:build js

package wgpu

import "github.com/oliverbestmann/webgpu/jsx"

// Submit as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-submit
func (g *Queue) Submit(commandBuffers ...*CommandBuffer) SubmissionIndex {
	g.jsValue.Call("submit", sliceToJS(commandBuffers, pointerToJS[CommandBuffer]))
	return SubmissionIndex(0)
}

// OnSubmittedWorkDone as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-onsubmittedworkdone
func (g *Queue) OnSubmittedWorkDone(callback QueueWorkDoneCallback) {
	// TODO should probably just schedule the callback using .then
	jsx.Await(g.jsValue.Call("onSubmittedWorkDone"))
	callback(QueueWorkDoneStatusSuccess)
}

// GetTimestampPeriod always returns 1: browsers resolve timestamp queries in
// nanoseconds, so no conversion is needed.
func (g *Queue) GetTimestampPeriod() float32 {
	return 1
}
