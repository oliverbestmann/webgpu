//go:build js && !goscript

package wgpu

import (
	"runtime"
	"syscall/js"
	"unsafe"
)

// TryWriteBuffer as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writebuffer
func (g *Queue) TryWriteBuffer(buffer *Buffer, offset uint64, data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}

	defer handleJsException(&err)
	address := uintptr(unsafe.Pointer(&data[0]))
	queueWriteBuffer.Invoke(g.jsValue, pointerToJS(buffer), offset, address, uint64(0), len(data))
	runtime.KeepAlive(data)
	return
}

// TryWriteTexture as described:
// https://gpuweb.github.io/gpuweb/#dom-gpuqueue-writetexture
func (g *Queue) TryWriteTexture(destination *TexelCopyTextureInfo, data []byte, dataLayout *TexelCopyBufferLayout, writeSize *Extent3D) (err error) {
	if len(data) == 0 {
		return nil
	}

	defer handleJsException(&err)
	address := uintptr(unsafe.Pointer(&data[0]))
	queueWriteTexture.Invoke(g.jsValue, pointerToJS(destination), address, len(data), pointerToJS(dataLayout), pointerToJS(writeSize))
	runtime.KeepAlive(data)
	return
}

var queueWriteBuffer js.Value
var queueWriteTexture js.Value

func init() {
	queueWriteBuffer = js.Global().Call("eval", `
		(queue, buf, offset, addr, x, n) => {
			const mem = wasm.instance.exports.mem.buffer;
			const data = new Uint8ClampedArray(mem, addr, n);
			return queue.writeBuffer(buf, offset, data, x, n);
		} 
	`)

	queueWriteTexture = js.Global().Call("eval", `
		(queue, tex, addr, n, layout, writeSize) => {
			const mem = wasm.instance.exports.mem.buffer;
			const data = new Uint8ClampedArray(mem, addr, n);
			return queue.writeTexture(tex, data, layout, writeSize);
		} 
	`)
}
