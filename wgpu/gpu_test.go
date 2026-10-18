//go:build !js

package wgpu

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// newTestDevice returns a device on the fallback adapter. Everything is
// released when the test ends.
func newTestDevice(t testing.TB) (*Device, *Queue) {
	t.Helper()

	instance := CreateInstance(nil)
	t.Cleanup(instance.Release)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(t, err)
	t.Cleanup(adapter.Release)

	device, err := adapter.RequestDevice(nil)
	require.NoError(t, err)
	t.Cleanup(device.Release)

	queue := device.GetQueue()
	t.Cleanup(queue.Release)

	return device, queue
}

// readBuffer copies src to a mappable buffer and returns its content.
func readBuffer(t testing.TB, device *Device, queue *Queue, src *Buffer) []byte {
	t.Helper()

	size := src.GetSize()
	staging := device.CreateBuffer(&BufferDescriptor{
		Size:  size,
		Usage: BufferUsageMapRead | BufferUsageCopyDst,
	})
	defer staging.Release()

	encoder := device.CreateCommandEncoder(nil)
	defer encoder.Release()
	encoder.CopyBufferToBuffer(src, 0, staging, 0, size)
	commands := encoder.Finish(nil)
	defer commands.Release()
	queue.Submit(commands)

	return mapRead(t, device, staging)
}

// mapRead maps a MapRead buffer and returns a copy of its content.
func mapRead(t testing.TB, device *Device, buffer *Buffer) []byte {
	t.Helper()

	var status MapAsyncStatus
	err := buffer.TryMapAsync(MapModeRead, 0, buffer.GetSize(), func(s MapAsyncStatus) { status = s })
	require.NoError(t, err)
	device.Poll(true, nil)
	require.Equal(t, MapAsyncStatusSuccess, status)

	data := append([]byte(nil), buffer.GetMappedRange(0, WholeMapSize)...)
	require.NoError(t, buffer.TryUnmap())
	return data
}

func TestBufferRoundTrip(t *testing.T) {
	device, queue := newTestDevice(t)

	src := device.CreateBufferInit(&BufferInitDescriptor{
		Contents: ToBytes([]uint32{1, 2, 3, 4, 5, 6, 7, 8}),
		Usage:    BufferUsageCopySrc | BufferUsageCopyDst,
	})
	defer src.Release()

	require.EqualValues(t, 32, src.GetSize())
	require.Equal(t, BufferUsageCopySrc|BufferUsageCopyDst, src.GetUsage())

	queue.WriteBuffer(src, 8, ToBytes([]uint32{30, 40}))

	require.Equal(t, []uint32{1, 2, 30, 40, 5, 6, 7, 8}, FromBytes[uint32](readBuffer(t, device, queue, src)))
}

func TestBufferMappedRangeWithOffset(t *testing.T) {
	device, _ := newTestDevice(t)

	buffer := device.CreateBuffer(&BufferDescriptor{
		Size:             16,
		Usage:            BufferUsageMapWrite | BufferUsageCopySrc,
		MappedAtCreation: true,
	})
	defer buffer.Release()

	copy(buffer.GetMappedRange(0, 16), ToBytes([]uint32{1, 2, 3, 4}))
	tail := buffer.GetMappedRange(8, WholeMapSize)
	require.Equal(t, []uint32{3, 4}, FromBytes[uint32](tail))
	require.NoError(t, buffer.TryUnmap())
}

const doubleShader = `
@group(0) @binding(0) var<storage, read_write> values: array<u32>;

@compute @workgroup_size(1)
fn main(@builtin(global_invocation_id) id: vec3<u32>) {
	values[id.x] = values[id.x] * 2u;
}
`

func TestComputeDoublesValues(t *testing.T) {
	device, queue := newTestDevice(t)

	shader := device.CreateShaderModule(&ShaderModuleDescriptor{
		WGSLSource: &ShaderSourceWGSL{Code: doubleShader},
	})
	defer shader.Release()

	bindGroupLayout := device.CreateBindGroupLayout(&BindGroupLayoutDescriptor{
		Entries: []BindGroupLayoutEntry{{
			Binding:    0,
			Visibility: ShaderStageCompute,
			Buffer:     BufferBindingLayout{Type: BufferBindingTypeStorage},
		}},
	})
	defer bindGroupLayout.Release()

	pipelineLayout := device.CreatePipelineLayout(&PipelineLayoutDescriptor{
		BindGroupLayouts: []*BindGroupLayout{bindGroupLayout},
	})
	defer pipelineLayout.Release()

	pipeline := device.CreateComputePipeline(&ComputePipelineDescriptor{
		Layout:  pipelineLayout,
		Compute: ProgrammableStageDescriptor{Module: shader, EntryPoint: "main"},
	})
	defer pipeline.Release()

	values := device.CreateBufferInit(&BufferInitDescriptor{
		Contents: ToBytes([]uint32{1, 2, 3, 4}),
		Usage:    BufferUsageStorage | BufferUsageCopySrc,
	})
	defer values.Release()

	bindGroup := device.CreateBindGroup(&BindGroupDescriptor{
		Layout:  bindGroupLayout,
		Entries: []BindGroupEntry{{Binding: 0, Buffer: values, Size: WholeSize}},
	})
	defer bindGroup.Release()

	encoder := device.CreateCommandEncoder(nil)
	defer encoder.Release()

	pass := encoder.BeginComputePass(nil)
	pass.SetPipeline(pipeline)
	pass.SetBindGroup(0, bindGroup, nil)
	pass.DispatchWorkgroups(4, 1, 1)
	pass.End()
	pass.Release()

	commands := encoder.Finish(nil)
	defer commands.Release()
	queue.Submit(commands)

	require.Equal(t, []uint32{2, 4, 6, 8}, FromBytes[uint32](readBuffer(t, device, queue, values)))
}

const quadShader = `
@group(0) @binding(0) var<uniform> color: vec4f;

@vertex
fn vs(@location(0) pos: vec2f) -> @builtin(position) vec4f {
	return vec4f(pos, 0.0, 1.0);
}

@fragment
fn fs() -> @location(0) vec4f {
	return color;
}
`

// TestRenderQuad clears a texture to red and draws a green quad over all of
// it, directly and through a render bundle. The quad vertices start at index
// 2 and are reached through baseVertex, so a wrong baseVertex leaves red pixels.
func TestRenderQuad(t *testing.T) {
	for _, useBundle := range []bool{false, true} {
		name := map[bool]string{false: "pass", true: "bundle"}[useBundle]
		t.Run(name, func(t *testing.T) { testRenderQuad(t, useBundle) })
	}
}

func testRenderQuad(t *testing.T, useBundle bool) {
	const size = 4
	const format = TextureFormatRGBA8Unorm

	device, queue := newTestDevice(t)

	shader := device.CreateShaderModule(&ShaderModuleDescriptor{
		WGSLSource: &ShaderSourceWGSL{Code: quadShader},
	})
	defer shader.Release()

	pipeline := device.CreateRenderPipeline(&RenderPipelineDescriptor{
		Vertex: VertexState{
			Module:     shader,
			EntryPoint: "vs",
			Buffers: []VertexBufferLayout{{
				ArrayStride: 8,
				StepMode:    VertexStepModeVertex,
				Attributes:  []VertexAttribute{{Format: VertexFormatFloat32x2, ShaderLocation: 0}},
			}},
		},
		Primitive:   PrimitiveState{Topology: PrimitiveTopologyTriangleList},
		Multisample: MultisampleState{Count: 1, Mask: 0xffffffff},
		Fragment: &FragmentState{
			Module:     shader,
			EntryPoint: "fs",
			Targets:    []ColorTargetState{{Format: format, WriteMask: ColorWriteMaskAll}},
		},
	})
	defer pipeline.Release()

	vertices := device.CreateBufferInit(&BufferInitDescriptor{
		// two unused vertices, then the quad
		Contents: ToBytes([]float32{0, 0, 0, 0, -1, -1, 1, -1, 1, 1, -1, 1}),
		Usage:    BufferUsageVertex,
	})
	defer vertices.Release()

	indices := device.CreateBufferInit(&BufferInitDescriptor{
		Contents: ToBytes([]uint16{0, 1, 2, 0, 2, 3}),
		Usage:    BufferUsageIndex,
	})
	defer indices.Release()

	uniform := device.CreateBufferInit(&BufferInitDescriptor{
		Contents: ToBytes([]float32{0, 1, 0, 1}),
		Usage:    BufferUsageUniform,
	})
	defer uniform.Release()

	bindGroupLayout := pipeline.GetBindGroupLayout(0)
	defer bindGroupLayout.Release()

	bindGroup := device.CreateBindGroup(&BindGroupDescriptor{
		Layout:  bindGroupLayout,
		Entries: []BindGroupEntry{{Binding: 0, Buffer: uniform, Size: WholeSize}},
	})
	defer bindGroup.Release()

	texture := device.CreateTexture(&TextureDescriptor{
		Usage:         TextureUsageRenderAttachment | TextureUsageCopySrc,
		Dimension:     TextureDimension2D,
		Size:          Extent3D{Width: size, Height: size, DepthOrArrayLayers: 1},
		Format:        format,
		MipLevelCount: 1,
		SampleCount:   1,
	})
	defer texture.Release()

	view := texture.CreateView(nil)
	defer view.Release()

	encoder := device.CreateCommandEncoder(nil)
	defer encoder.Release()

	pass := encoder.BeginRenderPass(&RenderPassDescriptor{
		ColorAttachments: []RenderPassColorAttachment{{
			View:       view,
			LoadOp:     LoadOpClear,
			StoreOp:    StoreOpStore,
			ClearValue: ColorRed,
		}},
	})

	if useBundle {
		bundleEncoder := device.CreateRenderBundleEncoder(&RenderBundleEncoderDescriptor{
			ColorFormats: []TextureFormat{format},
			SampleCount:  1,
		})
		bundleEncoder.SetPipeline(pipeline)
		bundleEncoder.SetBindGroup(0, bindGroup, nil)
		bundleEncoder.SetVertexBuffer(0, vertices, 0, WholeSize)
		bundleEncoder.SetIndexBuffer(indices, IndexFormatUint16, 0, WholeSize)
		bundleEncoder.DrawIndexed(6, 1, 0, 2, 0)
		bundle := bundleEncoder.Finish(nil)
		bundleEncoder.Release()
		defer bundle.Release()

		pass.ExecuteBundles(bundle)
	} else {
		pass.SetPipeline(pipeline)
		pass.SetBindGroup(0, bindGroup, nil)
		pass.SetVertexBuffer(0, vertices, 0, WholeSize)
		pass.SetIndexBuffer(indices, IndexFormatUint16, 0, WholeSize)
		pass.DrawIndexed(6, 1, 0, 2, 0)
	}

	pass.End()
	pass.Release()

	const bytesPerRow = CopyBytesPerRowAlignment
	output := device.CreateBuffer(&BufferDescriptor{
		Size:  bytesPerRow * size,
		Usage: BufferUsageMapRead | BufferUsageCopyDst,
	})
	defer output.Release()

	encoder.CopyTextureToBuffer(
		texture.AsImageCopy(),
		&TexelCopyBufferInfo{
			Buffer: output,
			Layout: TexelCopyBufferLayout{BytesPerRow: bytesPerRow, RowsPerImage: size},
		},
		&Extent3D{Width: size, Height: size, DepthOrArrayLayers: 1},
	)

	commands := encoder.Finish(nil)
	defer commands.Release()
	queue.Submit(commands)

	data := mapRead(t, device, output)
	for y := range size {
		for x := range size {
			offset := y*bytesPerRow + x*4
			require.Equal(t, []byte{0, 255, 0, 255}, data[offset:offset+4], "pixel %d,%d", x, y)
		}
	}
}

func TestInvalidDescriptorReturnsError(t *testing.T) {
	device, _ := newTestDevice(t)

	// MapRead may only be combined with CopyDst
	invalid := &BufferDescriptor{Size: 16, Usage: BufferUsageMapRead | BufferUsageStorage}

	_, err := device.TryCreateBuffer(invalid)
	var wgpuErr Error
	require.ErrorAs(t, err, &wgpuErr)
	require.Contains(t, wgpuErr.Context, "TryCreateBuffer")
	require.NotEmpty(t, wgpuErr.Wrapped.Error())

	_, err = device.TryCreateShaderModule(&ShaderModuleDescriptor{
		WGSLSource: &ShaderSourceWGSL{Code: "this is not wgsl"},
	})
	require.Error(t, err)

	// the error must not leak into the next call
	buffer, err := device.TryCreateBuffer(&BufferDescriptor{Size: 16, Usage: BufferUsageCopyDst})
	require.NoError(t, err)
	buffer.Release()
}

func TestCreateWithoutTryPanicsWithContext(t *testing.T) {
	device, _ := newTestDevice(t)

	defer func() {
		err, ok := recover().(error)
		require.True(t, ok, "expected a panic with an error")
		require.Contains(t, err.Error(), "CreateBuffer")

		var wgpuErr Error
		require.True(t, errors.As(err, &wgpuErr))
	}()

	device.CreateBuffer(&BufferDescriptor{Size: 16, Usage: BufferUsageMapRead | BufferUsageStorage})
}
