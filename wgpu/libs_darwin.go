//go:build darwin && !ios && cgo

package wgpu

// Imported per OS, so that a build only downloads the libraries it links.
import _ "github.com/oliverbestmann/webgpu/libs-darwin"
