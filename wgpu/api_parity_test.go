//go:build !js

package wgpu

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// apiOf returns the exported API of this package as seen by the given build
// context: a map from a declaration ("func (*Buffer) Unmap", "type Color") to
// its normalized signature or field list.
func apiOf(t *testing.T, ctx build.Context) map[string]string {
	t.Helper()

	pkg, err := ctx.ImportDir(".", 0)
	if err != nil {
		t.Fatal(err)
	}

	api := map[string]string{}
	fset := token.NewFileSet()

	for _, name := range slices.Concat(pkg.GoFiles, pkg.CgoFiles) {
		file, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}

		for _, decl := range file.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if !decl.Name.IsExported() {
					continue
				}

				key := "func " + decl.Name.Name
				if decl.Recv != nil {
					recv := types.ExprString(decl.Recv.List[0].Type)
					if !ast.IsExported(strings.TrimPrefix(recv, "*")) {
						continue
					}
					key = "func (" + recv + ") " + decl.Name.Name
				}

				api[key] = "(" + fieldTypes(decl.Type.Params) + ") (" + fieldTypes(decl.Type.Results) + ")"

			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.IsExported() {
							api["type "+spec.Name.Name] = typeShape(spec)
						}

					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if name.IsExported() {
								api[decl.Tok.String()+" "+name.Name] = ""
							}
						}
					}
				}
			}
		}
	}

	return api
}

// fieldTypes lists the types of a parameter or result list, without names.
func fieldTypes(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}

	var out []string
	for _, field := range fields.List {
		for range max(1, len(field.Names)) {
			out = append(out, types.ExprString(field.Type))
		}
	}

	return strings.Join(out, ", ")
}

// typeShape describes a type by its underlying kind, and for structs by the
// exported fields, so platform specific unexported fields are ignored.
func typeShape(spec *ast.TypeSpec) string {
	st, ok := spec.Type.(*ast.StructType)
	if !ok {
		return types.ExprString(spec.Type)
	}

	var fields []string
	for _, field := range st.Fields.List {
		typ := types.ExprString(field.Type)
		if len(field.Names) == 0 {
			fields = append(fields, typ)
		}
		for _, name := range field.Names {
			if name.IsExported() {
				fields = append(fields, name.Name+" "+typ)
			}
		}
	}

	return "struct{" + strings.Join(fields, "; ") + "}"
}

// knownAPIDifferences lists declarations that may be missing on one side or
// differ between native and JS, with the reason. An entry that no longer
// differs fails the test, so remove it when a gap gets closed.
var knownAPIDifferences = map[string]string{
	// Platform specific: surfaces come from window handles natively, from a canvas on web.
	"type SurfaceDescriptor":                 "platform",
	"type SurfaceSourceAndroidNativeWindow":  "platform",
	"type SurfaceSourceMetalLayer":           "platform",
	"type SurfaceSourceWaylandSurface":       "platform",
	"type SurfaceSourceWindowsHWND":          "platform",
	"type SurfaceSourceXcbWindow":            "platform",
	"type SurfaceSourceXlibWindow":           "platform",
	"type CanvasContext":                     "platform",
	"func NewCanvasContext":                  "platform",
	"func (CanvasContext) GetCurrentTexture": "platform",

	// wgpu-native extensions, not part of WebGPU.
	"func (*Instance) GenerateReport":                         "wgpu-native",
	"type GlobalReport":                                       "wgpu-native",
	"type HubReport":                                          "wgpu-native",
	"type RegistryReport":                                     "wgpu-native",
	"type ShaderModuleDescriptor":                             "wgpu-native: SPIR-V and GLSL sources",
	"type ShaderSourceGLSL":                                   "wgpu-native",
	"type PipelineLayoutDescriptor":                           "wgpu-native: immediates",
	"func (*ComputePassEncoder) BeginPipelineStatisticsQuery": "wgpu-native",
	"func (*ComputePassEncoder) EndPipelineStatisticsQuery":   "wgpu-native",
	"func (*RenderPassEncoder) BeginPipelineStatisticsQuery":  "wgpu-native",
	"func (*RenderPassEncoder) EndPipelineStatisticsQuery":    "wgpu-native",
	"func (*RenderPassEncoder) MultiDrawIndexedIndirect":      "wgpu-native",
	"func (*RenderPassEncoder) MultiDrawIndexedIndirectCount": "wgpu-native",
	"func (*RenderPassEncoder) MultiDrawIndirect":             "wgpu-native",
	"func (*RenderPassEncoder) MultiDrawIndirectCount":        "wgpu-native",

	// Only in WebGPU.
	"type BindGroupLayoutEntry":         "js: external textures",
	"type ExternalTextureBindingLayout": "js: external textures",

	// Gaps in the JS port.
	"func (*Adapter) GetFeatures":                           "missing in js",
	"func (*Device) GetFeatures":                            "missing in js",
	"func (*Device) HasFeature":                             "missing in js",
	"func (*ComputePassEncoder) DispatchWorkgroupsIndirect": "missing in js",
	"type RenderBundleDescriptor":                           "missing in js",
	"func (*RenderBundleEncoder) Draw":                      "missing in js",
	"func (*RenderBundleEncoder) DrawIndexed":               "missing in js",
	"func (*RenderBundleEncoder) DrawIndexedIndirect":       "missing in js",
	"func (*RenderBundleEncoder) DrawIndirect":              "missing in js",
	"func (*RenderBundleEncoder) Finish":                    "missing in js",
	"func (*RenderBundleEncoder) InsertDebugMarker":         "missing in js",
	"func (*RenderBundleEncoder) PopDebugGroup":             "missing in js",
	"func (*RenderBundleEncoder) PushDebugGroup":            "missing in js",
	"func (*RenderBundleEncoder) SetBindGroup":              "missing in js",
	"func (*RenderBundleEncoder) SetIndexBuffer":            "missing in js",
	"func (*RenderBundleEncoder) SetPipeline":               "missing in js",
	"func (*RenderBundleEncoder) SetVertexBuffer":           "missing in js",
	"func (*Texture) Destroy":                               "missing in js",
	"func (*Texture) GetDimension":                          "missing in js",
	"func (*Texture) GetUsage":                              "missing in js",
}

// TestNativeAndJSAPIMatch compares the exported declarations of the native and
// the JS builds (with and without goscript) by parsing the sources, so it runs
// without a wasm toolchain.
func TestNativeAndJSAPIMatch(t *testing.T) {
	native := build.Default
	native.CgoEnabled = true
	nativeAPI := apiOf(t, native)

	for _, tags := range [][]string{nil, {"goscript"}} {
		js := build.Default
		js.GOOS, js.GOARCH, js.CgoEnabled = "js", "wasm", false
		js.BuildTags = tags
		jsAPI := apiOf(t, js)

		for _, key := range slices.Sorted(maps.Keys(knownAPIDifferences)) {
			nativeSig, inNative := nativeAPI[key]
			jsSig, inJS := jsAPI[key]
			if inNative && inJS && nativeSig == jsSig {
				t.Errorf("tags %v: %s matches now, remove it from knownAPIDifferences", tags, key)
			}
		}

		for _, key := range slices.Sorted(maps.Keys(nativeAPI)) {
			if _, known := knownAPIDifferences[key]; known {
				continue
			}

			jsSig, ok := jsAPI[key]
			if !ok {
				t.Errorf("tags %v: only native: %s", tags, key)
			} else if jsSig != nativeAPI[key] {
				t.Errorf("tags %v: differs: %s\n  native: %s\n  js:     %s", tags, key, nativeAPI[key], jsSig)
			}
		}

		for _, key := range slices.Sorted(maps.Keys(jsAPI)) {
			if _, known := knownAPIDifferences[key]; known {
				continue
			}

			if _, ok := nativeAPI[key]; !ok {
				t.Errorf("tags %v: only js: %s", tags, key)
			}
		}
	}
}
