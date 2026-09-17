//go:build !js

package wgpu

// #include "gen_wgpu_wrappers.h"
import "C"

type stringView struct {
	view C.WGPUStringView
}

func stringViewOf(value string) stringView {
	if value == "" {
		return stringView{}
	}

	return stringView{
		view: C.WGPUStringView{
			data:   C.CString(value),
			length: C.size_t(len(value)),
		},
	}
}

func (v stringView) ToC() C.WGPUStringView {
	return v.view
}

func (v stringView) Release() {
	if v.view.data != nil {
		free(v.view.data)
	}
}
