package wgpu

import "testing"

func TestEnumFromJS(t *testing.T) {
	if v, ok := enumFromJS[TextureDimension]("3d"); !ok || v != TextureDimension3D {
		t.Errorf("3d: got %v, %v", v, ok)
	}
	if v, ok := enumFromJS[FeatureName]("texture-component-swizzle"); !ok || v != FeatureNameTextureComponentSwizzle {
		t.Errorf("last standard feature: got %v, %v", v, ok)
	}
	if _, ok := enumFromJS[FeatureName]("no-such-feature"); ok {
		t.Error("unknown name must not match")
	}
}
