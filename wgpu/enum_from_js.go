package wgpu

import "fmt"

// enumFromJS finds the enum value for a WebGPU string, ok is false for an
// unknown string. It scans from 1 until String returns "", so it only sees the
// standard values, which are numbered without gaps.
func enumFromJS[T interface {
	~uint32
	fmt.Stringer
}](s string) (_ T, ok bool) {
	for v := T(1); v.String() != ""; v++ {
		if v.String() == s {
			return v, true
		}
	}

	return 0, false
}
