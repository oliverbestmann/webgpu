package wgpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkCreateView(b *testing.B) {
	instance := CreateInstance(nil)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(b, err)
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	require.NoError(b, err)
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	texture := device.CreateTexture(&TextureDescriptor{
		Label:         "test",
		Usage:         TextureUsageCopyDst,
		Dimension:     TextureDimension2D,
		Size:          Extent3D{Width: 16, Height: 16, DepthOrArrayLayers: 1},
		Format:        TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})

	defer texture.Release()

	b.ReportAllocs()

	for range b.N {
		view := texture.CreateView(&TextureViewDescriptor{
			Label:           "test label",
			MipLevelCount:   1,
			ArrayLayerCount: 1,
		})

		view.Release()
	}
}

func BenchmarkCreateBindGroup(b *testing.B) {
	instance := CreateInstance(nil)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(b, err)
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	require.NoError(b, err)
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	texture := device.CreateTexture(&TextureDescriptor{
		Label:         "test",
		Usage:         TextureUsageTextureBinding,
		Dimension:     TextureDimension2D,
		Size:          Extent3D{Width: 16, Height: 16, DepthOrArrayLayers: 1},
		Format:        TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})
	defer texture.Release()

	view := texture.CreateView(&TextureViewDescriptor{
		Label:           "test label",
		MipLevelCount:   1,
		ArrayLayerCount: 1,
	})
	defer view.Release()

	bindGroupLayout := device.CreateBindGroupLayout(&BindGroupLayoutDescriptor{
		Label: "BindGroupLayout",
		Entries: []BindGroupLayoutEntry{
			{
				Binding:    0,
				Visibility: ShaderStageFragment,
				Texture: TextureBindingLayout{
					SampleType:    TextureSampleTypeFloat,
					ViewDimension: TextureViewDimension2D,
					Multisampled:  false,
				},
			},
		},
	})
	defer bindGroupLayout.Release()

	b.ReportAllocs()

	label := fmt.Sprintf("BindGroup %T", device)

	for range b.N {
		bindGroup := device.CreateBindGroup(&BindGroupDescriptor{
			Label:  label,
			Layout: bindGroupLayout,
			Entries: []BindGroupEntry{
				{
					TextureView: view,
				},
			},
		})

		bindGroup.Release()
	}
}

func BenchmarkBeginRenderPass(b *testing.B) {
	instance := CreateInstance(nil)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(b, err)
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	require.NoError(b, err)
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	b.ReportAllocs()

	texture := device.CreateTexture(&TextureDescriptor{
		Label:         "test",
		Usage:         TextureUsageRenderAttachment,
		Dimension:     TextureDimension2D,
		Size:          Extent3D{Width: 16, Height: 16, DepthOrArrayLayers: 1},
		Format:        TextureFormatRGBA8Unorm,
		MipLevelCount: 1,
		SampleCount:   1,
	})

	defer texture.Release()

	view := texture.CreateView(&TextureViewDescriptor{
		Label:           "test label",
		MipLevelCount:   1,
		ArrayLayerCount: 1,
	})

	defer view.Release()

	for range b.N {
		enc := device.CreateCommandEncoder(nil)

		pass := enc.BeginRenderPass(&RenderPassDescriptor{
			Label: "RenderPass",
			ColorAttachments: []RenderPassColorAttachment{
				{
					View:    view,
					LoadOp:  LoadOpClear,
					StoreOp: StoreOpStore,
				},
			},
		})

		pass.End()
		pass.Release()

		buf := enc.Finish(nil)
		buf.Release()

		enc.Release()
	}
}

func TestRequestDeviceWithAdapterLimits(t *testing.T) {
	instance := CreateInstance(nil)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(t, err)
	defer adapter.Release()

	limits := adapter.GetLimits()
	require.NotZero(t, limits.MaxNonSamplerBindings)

	device, err := adapter.RequestDevice(&DeviceDescriptor{
		RequiredLimits:     &limits,
		DeviceLostCallback: func(DeviceLostReason, string) {},
	})
	require.NoError(t, err)
	defer device.Release()

	require.NotZero(t, device.GetLimits().MaxNonSamplerBindings)
}

func TestOnSubmittedWorkDone(t *testing.T) {
	instance := CreateInstance(nil)

	adapter, err := instance.RequestAdapter(&RequestAdapterOptions{
		ForceFallbackAdapter: true,
	})
	require.NoError(t, err)
	defer adapter.Release()

	device, err := adapter.RequestDevice(nil)
	require.NoError(t, err)
	defer device.Release()

	queue := device.GetQueue()
	defer queue.Release()

	done := make(chan QueueWorkDoneStatus, 1)
	queue.Submit()
	queue.OnSubmittedWorkDone(func(status QueueWorkDoneStatus) { done <- status })
	device.Poll(true, nil)

	require.Equal(t, QueueWorkDoneStatusSuccess, <-done)
}
