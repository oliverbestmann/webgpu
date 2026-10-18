package wgpu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesRoundTrip(t *testing.T) {
	values := []uint32{1, 0xdeadbeef}
	bytes := ToBytes(values)
	require.Len(t, bytes, 8)
	require.Equal(t, values, FromBytes[uint32](bytes))

	require.Nil(t, ToBytes([]uint32{}))
	require.Nil(t, FromBytes[uint32](nil))
	require.Panics(t, func() { FromBytes[uint32](make([]byte, 6)) })
}

func TestVersionString(t *testing.T) {
	require.Equal(t, "0x1d000101", Version(0x1d000101).String())
}

func TestHandle(t *testing.T) {
	h := newHandle("value")
	other := newHandle("other")
	require.NotEqual(t, h.ToPointer(), other.ToPointer())

	h = lookupHandle(h.ToPointer())
	require.Equal(t, "value", h.Value())

	h.Delete()
	require.Panics(t, func() { h.Value() })
	require.Panics(t, func() { h.Delete() })

	require.Equal(t, "other", other.Value())
	other.Delete()
}
