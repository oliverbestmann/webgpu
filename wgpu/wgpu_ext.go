package wgpu

const (
	// Buffer-Texture copies must have `TextureDataLayout.BytesPerRow` aligned to this number.
	//
	// This doesn't apply to `(*Queue).TryWriteTexture()`.
	CopyBytesPerRowAlignment = 256
	// An offset into the query resolve buffer has to be aligned to this.
	QueryResolveBufferAlignment = 256
	// Buffer to buffer copy as well as buffer clear offsets and sizes must be aligned to this number.
	CopyBufferAlignment = 4
	// Size to align mappings.
	MapAlignment = 8
	// Vertex buffer strides have to be aligned to this number.
	VertexStrideAlignment = 4
	// Alignment all push constants need
	PushConstantAlignment = 4
	// Maximum queries in a query set
	QuerySetMaxQueries = 8192
	// Size of a single piece of query data.
	QuerySize = 8
)

var (
	ColorTransparent = Color{0, 0, 0, 0}
	ColorBlack       = Color{0, 0, 0, 1}
	ColorWhite       = Color{1, 1, 1, 1}
	ColorRed         = Color{1, 0, 0, 1}
	ColorGreen       = Color{0, 1, 0, 1}
	ColorBlue        = Color{0, 0, 1, 1}
)

var (
	BlendComponentReplace = BlendComponent{
		SrcFactor: BlendFactorOne,
		DstFactor: BlendFactorZero,
		Operation: BlendOperationAdd,
	}

	BlendComponentOver = BlendComponent{
		SrcFactor: BlendFactorOne,
		DstFactor: BlendFactorOneMinusSrcAlpha,
		Operation: BlendOperationAdd,
	}

	BlendComponentAdd = BlendComponent{
		SrcFactor: BlendFactorSrcAlpha,
		DstFactor: BlendFactorOne,
		Operation: BlendOperationAdd,
	}

	BlendComponentMultiply = BlendComponent{
		SrcFactor: BlendFactorDst,
		DstFactor: BlendFactorZero,
		Operation: BlendOperationAdd,
	}
)

var (
	BlendStateReplace = BlendState{
		Color: BlendComponentReplace,
		Alpha: BlendComponentReplace,
	}

	BlendStateAdd = BlendState{
		Color: BlendComponentAdd,
		Alpha: BlendComponentAdd,
	}

	BlendStateMultiply = BlendState{
		Color: BlendComponentMultiply,
		Alpha: BlendComponentReplace,
	}

	BlendStateAlphaBlending = BlendState{
		Color: BlendComponent{
			SrcFactor: BlendFactorSrcAlpha,
			DstFactor: BlendFactorOneMinusSrcAlpha,
			Operation: BlendOperationAdd,
		},
		Alpha: BlendComponentOver,
	}

	BlendStatePremultipliedAlphaBlending = BlendState{
		Color: BlendComponentOver,
		Alpha: BlendComponentOver,
	}
)

func (v VertexFormat) Size() uint64 {
	return uint64(v.ByteSize())
}
