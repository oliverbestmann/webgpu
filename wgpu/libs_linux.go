//go:build linux && !android && cgo

package wgpu

// Imported per OS, so that a build only downloads the libraries it links.
import _ "github.com/oliverbestmann/webgpu/libs-linux"
