//go:build cgo

// Package needscgo is empty and only builds with cgo. The wgpu package imports it
// so that a native build without cgo fails with an error naming this package.
package needscgo
