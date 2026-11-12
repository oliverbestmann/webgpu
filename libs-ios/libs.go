package lib

/*
#cgo ios,amd64 LDFLAGS: -L${SRCDIR}/amd64 -lwgpu_native
#cgo ios,arm64,!iossimulator LDFLAGS: -L${SRCDIR}/arm64 -lwgpu_native
#cgo ios,arm64,iossimulator LDFLAGS: -L${SRCDIR}/arm64-simulator -lwgpu_native
*/
import "C"
