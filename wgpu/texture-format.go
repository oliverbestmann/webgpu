package wgpu

// ByteSize returns the number of bytes per pixel, or zero if unknown or variable
func (v TextureFormat) ByteSize() uint32 {
	switch v {
	case TextureFormatR8Unorm,
		TextureFormatR8Snorm,
		TextureFormatR8Uint,
		TextureFormatR8Sint:
		return 1

	case TextureFormatR16Unorm,
		TextureFormatR16Snorm,
		TextureFormatR16Uint,
		TextureFormatR16Sint,
		TextureFormatR16Float,
		TextureFormatRG8Unorm,
		TextureFormatRG8Snorm,
		TextureFormatRG8Uint,
		TextureFormatRG8Sint:
		return 2

	case TextureFormatR32Float,
		TextureFormatR32Uint,
		TextureFormatR32Sint,
		TextureFormatRG16Unorm,
		TextureFormatRG16Snorm,
		TextureFormatRG16Uint,
		TextureFormatRG16Sint,
		TextureFormatRG16Float,
		TextureFormatRGBA8Unorm,
		TextureFormatRGBA8UnormSrgb,
		TextureFormatRGBA8Snorm,
		TextureFormatRGBA8Uint,
		TextureFormatRGBA8Sint,
		TextureFormatBGRA8Unorm,
		TextureFormatBGRA8UnormSrgb,
		TextureFormatRGB10A2Uint,
		TextureFormatRGB10A2Unorm,
		TextureFormatRG11B10Ufloat,
		TextureFormatRGB9E5Ufloat:
		return 4

	case TextureFormatRG32Float,
		TextureFormatRG32Uint,
		TextureFormatRG32Sint,
		TextureFormatRGBA16Unorm,
		TextureFormatRGBA16Snorm,
		TextureFormatRGBA16Uint,
		TextureFormatRGBA16Sint,
		TextureFormatRGBA16Float,
		TextureFormatRGBA32Float:
		return 8

	case TextureFormatRGBA32Uint, TextureFormatRGBA32Sint:
		return 16

	case TextureFormatStencil8:
		return 1

	case TextureFormatDepth16Unorm:
		return 2

	case TextureFormatDepth24Plus:
		return 4

	case TextureFormatDepth32Float:
		return 4

	default:
		// i.e. unknown, variable length
		return 0
	}
}
