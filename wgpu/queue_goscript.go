//go:build js && goscript

package wgpu

// GoScript compiles Go to JavaScript, so there is no wasm memory to view;
// the data is copied into a JavaScript Uint8Array instead.

// TryWriteBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writebuffer
func (g *Queue) TryWriteBuffer(buffer *Buffer, offset uint64, data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}

	defer handleJsException(&err)
	g.jsValue.Call("writeBuffer", pointerToJS(buffer), offset, bytesToJS(data), 0, len(data))
	return
}

// TryWriteTexture as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writetexture
func (g *Queue) TryWriteTexture(destination *TexelCopyTextureInfo, data []byte, dataLayout *TexelCopyBufferLayout, writeSize *Extent3D) (err error) {
	if len(data) == 0 {
		return nil
	}

	defer handleJsException(&err)
	g.jsValue.Call("writeTexture", pointerToJS(destination), bytesToJS(data), pointerToJS(dataLayout), pointerToJS(writeSize))
	return
}
