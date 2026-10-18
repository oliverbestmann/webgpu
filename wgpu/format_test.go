package wgpu

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
	"testing"
)

// TestTextureFormatByteSize derives the size of plain formats from their
// WebGPU name, so it does not repeat the table in ByteSize.
func TestTextureFormatByteSize(t *testing.T) {
	plain := regexp.MustCompile(`^(r|rg|rgba|bgra)(8|16|32)(unorm|snorm|uint|sint|float)(-srgb)?$`)
	special := map[string]uint32{
		"rgb10a2uint": 4, "rgb10a2unorm": 4, "rg11b10ufloat": 4, "rgb9e5ufloat": 4,
		"stencil8": 1, "depth16unorm": 2, "depth24plus": 4, "depth32float": 4,
	}

	for v := TextureFormat(1); v.String() != ""; v++ {
		name := v.String()

		var want uint32
		if m := plain.FindStringSubmatch(name); m != nil {
			bits, _ := strconv.Atoi(m[2])
			want = uint32(len(m[1]) * bits / 8)
		} else {
			// compressed and combined depth-stencil formats have no fixed size
			want = special[name]
		}

		if got := v.ByteSize(); got != want {
			t.Errorf("%s: ByteSize is %d, want %d", name, got, want)
		}
	}
}

func TestVertexFormatByteSize(t *testing.T) {
	plain := regexp.MustCompile(`^(?:uint|sint|unorm|snorm|float)(8|16|32)(?:x([234]))?$`)
	special := map[string]uint32{"unorm10-10-10-2": 4, "unorm8x4-bgra": 4}

	for v := VertexFormat(1); v.String() != ""; v++ {
		name := v.String()

		want, ok := special[name]
		if m := plain.FindStringSubmatch(name); m != nil {
			bits, _ := strconv.Atoi(m[1])
			count := 1
			if m[2] != "" {
				count, _ = strconv.Atoi(m[2])
			}
			want, ok = uint32(bits/8*count), true
		}

		if !ok {
			t.Errorf("%s: unknown vertex format, extend the test", name)
			continue
		}

		if got := v.ByteSize(); got != want {
			t.Errorf("%s: ByteSize is %d, want %d", name, got, want)
		}
	}
}

// TestEnumStrings checks the generated code: every constant has a non-empty
// String name, and no two constants of a type share one, or enumFromJS and the
// JS bindings would mix them up.
func TestEnumStrings(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "gen_enums.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}

	consts := map[string][]string{}         // type -> constant names
	names := map[string]map[string]string{} // type -> constant -> String result

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok && spec.Type != nil {
					typ := spec.Type.(*ast.Ident).Name
					consts[typ] = append(consts[typ], spec.Names[0].Name)
				}
			}

		case *ast.FuncDecl:
			if decl.Name.Name != "String" || decl.Recv == nil {
				continue
			}

			typ := decl.Recv.List[0].Type.(*ast.Ident).Name
			names[typ] = map[string]string{}

			ast.Inspect(decl.Body, func(n ast.Node) bool {
				clause, ok := n.(*ast.CaseClause)
				if !ok || len(clause.List) != 1 {
					return true
				}

				ret := clause.Body[0].(*ast.ReturnStmt).Results[0].(*ast.BasicLit)
				name, _ := strconv.Unquote(ret.Value)
				names[typ][clause.List[0].(*ast.Ident).Name] = name
				return true
			})
		}
	}

	if len(consts) == 0 {
		t.Fatal("found no enums in gen_enums.go")
	}

	for typ, constNames := range consts {
		seen := map[string]string{}

		for _, constName := range constNames {
			name, ok := names[typ][constName]
			if !ok || name == "" {
				t.Errorf("%s has no String name", constName)
				continue
			}

			if other, dup := seen[name]; dup {
				t.Errorf("%s and %s are both %q", other, constName, name)
			}
			seen[name] = constName
		}
	}
}
