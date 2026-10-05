//go:build !js

package wgpu

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// wgpu-native returns an empty AdapterInfo string as a dangling pointer (0x1).
// If GetInfo holds one on the Go stack and the stack grows, the runtime aborts
// with "invalid pointer found on stack". Calling GetInfo at every stack depth
// on fresh goroutines makes the stack grow inside GetInfo at some depth.
//
// Fresh goroutines must start at the fixed stack size for the depths to line
// up, so the test runs itself again with adaptive starting stacks off.
func TestGetInfoSurvivesStackGrowth(t *testing.T) {
	if os.Getenv("WGPU_TEST_GETINFO_STACK_CHILD") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestGetInfoSurvivesStackGrowth$", "-test.count=1")
		cmd.Env = append(os.Environ(),
			"WGPU_TEST_GETINFO_STACK_CHILD=1",
			"GODEBUG="+os.Getenv("GODEBUG")+",adaptivestackstart=0",
		)
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "%s", out)
		return
	}

	instance := CreateInstance(nil)
	defer instance.Release()

	adapters := instance.EnumerateAdapters(nil)
	require.NotEmpty(t, adapters)
	defer func() {
		for _, adapter := range adapters {
			adapter.Release()
		}
	}()

	// Coarse steps walk past the first stack growths; fine steps shift each
	// one by a few bytes, so some call lands with the stack limit inside GetInfo.
	for coarse := range 192 {
		for fine := range 16 {
			done := make(chan struct{})
			go func() {
				defer close(done)
				atStackDepth(coarse, func() {
					atStackDepth8(fine, func() {
						for _, adapter := range adapters {
							adapter.GetInfo()
						}
					})
				})
			}()
			<-done
		}
	}
}

// atStackDepth calls fn depth frames of about 100 bytes below the caller.
//
//go:noinline
func atStackDepth(depth int, fn func()) {
	var pad [64]byte
	if depth == 0 {
		fn()
		return
	}
	atStackDepth(depth-1, fn)
	pad[0]++
}

// atStackDepth8 calls fn depth of the smallest frames below the caller.
//
//go:noinline
func atStackDepth8(depth int, fn func()) {
	if depth == 0 {
		fn()
		return
	}
	atStackDepth8(depth-1, fn)
}
