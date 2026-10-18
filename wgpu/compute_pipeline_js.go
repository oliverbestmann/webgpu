//go:build js

package wgpu

// ComputePipeline Descriptor as described:
// https://gpuweb.github.io/gpuweb/#dictdef-gpucomputepipelinedescriptor
type ComputePipelineDescriptor struct {
	Label   string
	Layout  *PipelineLayout
	Compute ProgrammableStageDescriptor
}

func (g ComputePipelineDescriptor) toJS() any {
	result := map[string]any{"label": g.Label}
	if g.Layout != nil {
		result["layout"] = pointerToJS(g.Layout)
	} else {
		result["layout"] = "auto"
	}
	result["compute"] = g.Compute.toJS()
	return result
}

func (g *ComputePipeline) GetBindGroupLayout(groupIndex uint32) *BindGroupLayout {
	jsValue := g.jsValue.Call("getBindGroupLayout", groupIndex)
	return &BindGroupLayout{jsValue}
}
