//go:build js

package wgpu

import "syscall/js"

// sliceToJS converts each element with fn into a slice that js.ValueOf turns into a JS array.
func sliceToJS[S any](slice []S, fn func(S) js.Value) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func SetLogLevel(level LogLevel) {}
func GetVersion() Version        { return 0 }
