//go:build !js

package wgpu

// #include "gen_wgpu_wrappers.h"
import "C"
import (
	"errors"
)

type ComputePassDescriptor struct {
	Label           string
	TimestampWrites *PassTimestampWrites
}

func cPassTimestampWrites(w *PassTimestampWrites) (*C.WGPUPassTimestampWrites, func()) {
	if w == nil || w.QuerySet == nil {
		return nil, func() {}
	}
	writes := callocOne[C.WGPUPassTimestampWrites]()
	writes.querySet = w.QuerySet.ref
	writes.beginningOfPassWriteIndex = C.uint32_t(w.BeginningOfPassWriteIndex)
	writes.endOfPassWriteIndex = C.uint32_t(w.EndOfPassWriteIndex)
	return writes, func() { free(writes) }
}

func (p *CommandEncoder) BeginComputePass(descriptor *ComputePassDescriptor) *ComputePassEncoder {
	var desc C.WGPUComputePassDescriptor

	if descriptor != nil {
		label := stringViewOf(descriptor.Label)
		defer label.Release()
		desc.label = label.ToC()

		timestampWrites, release := cPassTimestampWrites(descriptor.TimestampWrites)
		defer release()
		desc.timestampWrites = timestampWrites
	}

	ref := C.wgpuCommandEncoderBeginComputePass(p.ref, &desc)
	if ref == nil {
		err := errors.New("failed to acquire ComputePassEncoder")
		panic(wrap(err, ""))
	}

	C.wgpuDeviceAddRef(p.device)
	return releaseOnGC(&ComputePassEncoder{device: p.device, ref: ref})
}

func (p *CommandEncoder) TryBeginRenderPass(descriptor *RenderPassDescriptor) (*RenderPassEncoder, error) {
	var desc C.WGPURenderPassDescriptor

	if descriptor != nil {
		label := stringViewOf(descriptor.Label)
		defer label.Release()
		desc.label = label.ToC()

		colorAttachmentCount := len(descriptor.ColorAttachments)
		if colorAttachmentCount > 0 {
			colorAttachments, colorAttachmentsSlice := callocSlice[C.WGPURenderPassColorAttachment](colorAttachmentCount)
			defer free(colorAttachments)

			for i, v := range descriptor.ColorAttachments {
				colorAttachment := C.WGPURenderPassColorAttachment{
					loadOp:     C.WGPULoadOp(v.LoadOp),
					storeOp:    C.WGPUStoreOp(v.StoreOp),
					depthSlice: C.WGPU_DEPTH_SLICE_UNDEFINED,
					clearValue: C.WGPUColor{
						r: C.double(v.ClearValue.R),
						g: C.double(v.ClearValue.G),
						b: C.double(v.ClearValue.B),
						a: C.double(v.ClearValue.A),
					},
				}
				if v.View != nil {
					colorAttachment.view = v.View.ref
				}
				if v.ResolveTarget != nil {
					colorAttachment.resolveTarget = v.ResolveTarget.ref
				}

				colorAttachmentsSlice[i] = colorAttachment
			}

			desc.colorAttachmentCount = C.size_t(colorAttachmentCount)
			desc.colorAttachments = colorAttachments
		}

		if descriptor.DepthStencilAttachment != nil {
			depthStencilAttachment := callocOne[C.WGPURenderPassDepthStencilAttachment]()
			defer free(depthStencilAttachment)

			if descriptor.DepthStencilAttachment.View != nil {
				depthStencilAttachment.view = descriptor.DepthStencilAttachment.View.ref
			}
			depthStencilAttachment.depthLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.DepthLoadOp)
			depthStencilAttachment.depthStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.DepthStoreOp)
			depthStencilAttachment.depthClearValue = C.float(descriptor.DepthStencilAttachment.DepthClearValue)
			depthStencilAttachment.depthReadOnly = cBool(descriptor.DepthStencilAttachment.DepthReadOnly)
			depthStencilAttachment.stencilLoadOp = C.WGPULoadOp(descriptor.DepthStencilAttachment.StencilLoadOp)
			depthStencilAttachment.stencilStoreOp = C.WGPUStoreOp(descriptor.DepthStencilAttachment.StencilStoreOp)
			depthStencilAttachment.stencilClearValue = C.uint32_t(descriptor.DepthStencilAttachment.StencilClearValue)
			depthStencilAttachment.stencilReadOnly = cBool(descriptor.DepthStencilAttachment.StencilReadOnly)

			desc.depthStencilAttachment = depthStencilAttachment
		}

		timestampWrites, release := cPassTimestampWrites(descriptor.TimestampWrites)
		defer release()
		desc.timestampWrites = timestampWrites
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	ref := C.go_wgpuCommandEncoderBeginRenderPass(p.device, errh.ToPointer(), p.ref, &desc)
	if err := errh.ToError(); err != nil {
		return nil, err
	}

	C.wgpuDeviceAddRef(p.device)
	return releaseOnGC(&RenderPassEncoder{device: p.device, ref: ref}), nil
}

func (p *CommandEncoder) TryClearBuffer(buffer *Buffer, offset uint64, size uint64) error {
	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderClearBuffer(
		p.device,
		errh.errStr,
		p.ref,
		buffer.ref,
		C.uint64_t(offset),
		C.uint64_t(size),
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryCopyBufferToBuffer(source *Buffer, sourceOffset uint64, destination *Buffer, destinationOffset uint64, size uint64) error {
	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderCopyBufferToBuffer(
		p.device,
		errh.ToPointer(),
		p.ref,
		source.ref,
		C.uint64_t(sourceOffset),
		destination.ref,
		C.uint64_t(destinationOffset),
		C.uint64_t(size),
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryCopyBufferToTexture(source *TexelCopyBufferInfo, destination *TexelCopyTextureInfo, copySize *Extent3D) error {
	var src C.WGPUTexelCopyBufferInfo
	if source != nil {
		if source.Buffer != nil {
			src.buffer = source.Buffer.ref
		}
		src.layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(source.Layout.Offset),
			bytesPerRow:  C.uint32_t(source.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(source.Layout.RowsPerImage),
		}
	}

	var dst C.WGPUTexelCopyTextureInfo
	if destination != nil {
		dst = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderCopyBufferToTexture(
		p.device,
		errh.ToPointer(),
		p.ref,
		&src,
		&dst,
		&cpySize,
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryCopyTextureToBuffer(source *TexelCopyTextureInfo, destination *TexelCopyBufferInfo, copySize *Extent3D) error {
	var src C.WGPUTexelCopyTextureInfo
	if source != nil {
		src = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUTexelCopyBufferInfo
	if destination != nil {
		if destination.Buffer != nil {
			dst.buffer = destination.Buffer.ref
		}
		dst.layout = C.WGPUTexelCopyBufferLayout{
			offset:       C.uint64_t(destination.Layout.Offset),
			bytesPerRow:  C.uint32_t(destination.Layout.BytesPerRow),
			rowsPerImage: C.uint32_t(destination.Layout.RowsPerImage),
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderCopyTextureToBuffer(
		p.device,
		errh.ToPointer(),
		p.ref,
		&src,
		&dst,
		&cpySize,
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryCopyTextureToTexture(source *TexelCopyTextureInfo, destination *TexelCopyTextureInfo, copySize *Extent3D) error {
	var src C.WGPUTexelCopyTextureInfo
	if source != nil {
		src = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(source.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(source.Origin.X),
				y: C.uint32_t(source.Origin.Y),
				z: C.uint32_t(source.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(source.Aspect),
		}
		if source.Texture != nil {
			src.texture = source.Texture.ref
		}
	}

	var dst C.WGPUTexelCopyTextureInfo
	if destination != nil {
		dst = C.WGPUTexelCopyTextureInfo{
			mipLevel: C.uint32_t(destination.MipLevel),
			origin: C.WGPUOrigin3D{
				x: C.uint32_t(destination.Origin.X),
				y: C.uint32_t(destination.Origin.Y),
				z: C.uint32_t(destination.Origin.Z),
			},
			aspect: C.WGPUTextureAspect(destination.Aspect),
		}
		if destination.Texture != nil {
			dst.texture = destination.Texture.ref
		}
	}

	var cpySize C.WGPUExtent3D
	if copySize != nil {
		cpySize = C.WGPUExtent3D{
			width:              C.uint32_t(copySize.Width),
			height:             C.uint32_t(copySize.Height),
			depthOrArrayLayers: C.uint32_t(copySize.DepthOrArrayLayers),
		}
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderCopyTextureToTexture(
		p.device,
		errh.ToPointer(),
		p.ref,
		&src,
		&dst,
		&cpySize,
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryFinish(descriptor *CommandBufferDescriptor) (*CommandBuffer, error) {
	var desc *C.WGPUCommandBufferDescriptor

	if descriptor != nil && descriptor.Label != "" {
		label := C.CString(descriptor.Label)
		defer free(label)

		desc = &C.WGPUCommandBufferDescriptor{
			label: C.WGPUStringView{data: label, length: C.WGPU_STRLEN},
		}
	}

	errh := acquireErrorCallback()
	defer errh.Done()

	ref := C.go_wgpuCommandEncoderFinish(
		p.device,
		errh.errStr,
		p.ref,
		desc,
	)
	if err := errh.ToError(); err != nil {
		C.wgpuCommandBufferRelease(ref)
		return nil, err
	}

	return releaseOnGC(&CommandBuffer{ref: ref}), nil
}

func (p *CommandEncoder) TryInsertDebugMarker(markerLabel string) error {
	label := stringViewOf(markerLabel)
	defer label.Release()

	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderInsertDebugMarker(
		p.device,
		errh.ToPointer(),
		p.ref,
		label.ToC(),
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryPopDebugGroup() error {
	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderPopDebugGroup(
		p.device,
		errh.ToPointer(),
		p.ref,
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryPushDebugGroup(groupLabel string) error {
	errh := acquireErrorCallback()
	defer errh.Done()

	label := stringViewOf(groupLabel)
	defer label.Release()
	labelC := label.ToC()

	C.go_wgpuCommandEncoderPushDebugGroup(
		p.device,
		errh.ToPointer(),
		p.ref,
		labelC,
	)

	return errh.ToError()
}

func (p *CommandEncoder) TryResolveQuerySet(querySet *QuerySet, firstQuery uint32, queryCount uint32, destination *Buffer, destinationOffset uint64) error {
	errh := acquireErrorCallback()
	defer errh.Done()

	C.go_wgpuCommandEncoderResolveQuerySet(
		p.device,
		errh.ToPointer(),
		p.ref,
		querySet.ref,
		C.uint32_t(firstQuery),
		C.uint32_t(queryCount),
		destination.ref,
		C.uint64_t(destinationOffset),
	)

	return errh.ToError()
}
