//go:build js

package wgpu

import (
	"syscall/js"
	"testing"
)

func jsonOf(v js.Value) string {
	return js.Global().Get("JSON").Call("stringify", v).String()
}

func TestRenderPassDescriptorReuse(t *testing.T) {
	view := &TextureView{jsValue: js.ValueOf("view")}

	full := &RenderPassDescriptor{
		ColorAttachments: []RenderPassColorAttachment{
			{View: view, LoadOp: LoadOpClear, StoreOp: StoreOpStore, ClearValue: Color{R: 1}},
			{View: view, LoadOp: LoadOpLoad, StoreOp: StoreOpStore},
		},
		DepthStencilAttachment: &RenderPassDepthStencilAttachment{
			View: view, DepthLoadOp: LoadOpClear, DepthStoreOp: StoreOpStore, DepthClearValue: 1,
			StencilReadOnly: true,
		},
	}

	got := jsonOf(full.toJS())
	want := `{"colorAttachments":[` +
		`{"view":"view","loadOp":"clear","storeOp":"store","clearValue":[1,0,0,0]},` +
		`{"view":"view","loadOp":"load","storeOp":"store","clearValue":[0,0,0,0]}],` +
		`"depthStencilAttachment":{"view":"view","depthClearValue":1,"stencilClearValue":0,` +
		`"depthReadOnly":false,"stencilReadOnly":true,"depthLoadOp":"clear","depthStoreOp":"store"}}`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}

	// a smaller descriptor must not keep fields or attachments of the previous one
	small := &RenderPassDescriptor{
		ColorAttachments: []RenderPassColorAttachment{{View: view, LoadOp: LoadOpLoad, StoreOp: StoreOpDiscard}},
	}

	got = jsonOf(small.toJS())
	want = `{"colorAttachments":[{"view":"view","loadOp":"load","storeOp":"discard","clearValue":[0,0,0,0]}]}`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}

func TestBindGroupDescriptorReuse(t *testing.T) {
	layout := &BindGroupLayout{jsValue: js.ValueOf("layout")}
	buffer := &Buffer{jsValue: js.ValueOf("buffer")}
	sampler := &Sampler{jsValue: js.ValueOf("sampler")}

	first := BindGroupDescriptor{
		Layout: layout,
		Entries: []BindGroupEntry{
			{Binding: 0, Buffer: buffer, Offset: 16, Size: LimitU64Undefined},
			{Binding: 1, Sampler: sampler},
		},
	}

	got := jsonOf(first.toJS())
	want := `{"layout":"layout","entries":[{"binding":0,"resource":{"buffer":"buffer","offset":16}},{"binding":1,"resource":"sampler"}]}`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}

	// entry 0 switches from a buffer to a sampler, the entry count shrinks
	second := BindGroupDescriptor{
		Layout:  layout,
		Entries: []BindGroupEntry{{Binding: 2, Sampler: sampler}},
	}

	got = jsonOf(second.toJS())
	want = `{"layout":"layout","entries":[{"binding":2,"resource":"sampler"}]}`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
