//go:build js

package wgpu

import "syscall/js"

// ShaderModuleDescriptor as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpushadermoduledescriptor
type ShaderModuleDescriptor struct {
	Label      string
	WGSLSource *ShaderSourceWGSL
}

func (g ShaderModuleDescriptor) toJS() js.Value {
	return js.ValueOf(map[string]any{
		"code": g.WGSLSource.Code,
	})
}
